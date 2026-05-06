package service

import (
	"io"
	"log/slog"
	"sort"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/ml"
	"github.com/terracodum/expensemind/backend/internal/parser"
	"github.com/terracodum/expensemind/backend/internal/repository"
)

type MLClient interface {
	Predict(req ml.PredictRequest) (domain.Forecast, error)
}

type service struct {
	parserFactory parser.Factory
	txRepo        repository.TransactionRepository
	recurRepo     repository.RecurringRuleRepository
	forecastRepo  repository.ForecastJobRepository
	ml            MLClient
}

func New(
	pf parser.Factory,
	txRepo repository.TransactionRepository,
	recurRepo repository.RecurringRuleRepository,
	forecastRepo repository.ForecastJobRepository,
	ml MLClient,
) Service {
	return &service{parserFactory: pf, txRepo: txRepo, recurRepo: recurRepo, forecastRepo: forecastRepo, ml: ml}
}

func (s *service) forecastWorker(id int) {
	fail := func(err error) {
		slog.Error("forecast worker failed", "job_id", id, "err", err)
		s.forecastRepo.Update(domain.ForecastJob{ID: id, Status: "failed"})
	}

	s.forecastRepo.Update(domain.ForecastJob{ID: id, Status: "running"})

	today := time.Now()

	trans, err := s.txRepo.FindForForecast(time.Time{}, today)
	if err != nil {
		fail(err)
		return
	}

	rules, err := s.recurRepo.FindActive(today)
	if err != nil {
		fail(err)
		return
	}

	byDay := make(map[time.Time][]domain.Transaction)
	for _, tx := range trans {
		if tx.Category == "transfer" {
			continue
		}
		day := tx.Date.Truncate(24 * time.Hour)
		byDay[day] = append(byDay[day], tx)
	}

	days := make([]time.Time, 0, len(byDay))
	for d := range byDay {
		days = append(days, d)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Before(days[j]) })

	var balance, totalDailyExpense float64
	timepoints := make([]ml.TimePoint, 0, len(days))

	for t, day := range days {
		txs := byDay[day]
		var foodTotal, transportTotal, entertainmentTotal, dayAmount float64
		for _, tx := range txs {
			dayAmount += tx.Amount
			switch tx.Category {
			case "food":
				foodTotal += tx.Amount
			case "transport":
				transportTotal += tx.Amount
			case "entertainment":
				entertainmentTotal += tx.Amount
			}
		}
		balance += dayAmount
		totalDailyExpense += dayAmount

		var avgSize float64
		if len(txs) > 0 {
			avgSize = dayAmount / float64(len(txs))
		}

		dow := int(day.Weekday())
		if dow == 0 {
			dow = 7
		}

		timepoints = append(timepoints, ml.TimePoint{
			T:                  t + 1,
			Balance:            balance,
			DayOfWeek:          dow,
			IsWeekend:          day.Weekday() == time.Saturday || day.Weekday() == time.Sunday,
			FoodTotal:          foodTotal,
			TransportTotal:     transportTotal,
			EntertainmentTotal: entertainmentTotal,
			AvgTransactionSize: avgSize,
			TransactionCount:   len(txs),
		})
	}

	lastDay := today
	if len(days) > 0 {
		lastDay = days[len(days)-1]
	}

	incomeEvents := []ml.IncomeEvent{}
	for _, rule := range rules {
		next := nextRuleOccurrence(lastDay, rule.Day)
		t := int(next.Sub(lastDay).Hours() / 24)
		if t > 0 && t <= 30 {
			incomeEvents = append(incomeEvents, ml.IncomeEvent{
				T:      t,
				Amount: rule.Amount,
				Label:  rule.Label,
			})
		}
	}

	var avgDailyExpense float64
	if len(days) > 0 {
		avgDailyExpense = totalDailyExpense / float64(len(days))
	}

	req := ml.PredictRequest{
		Timeseries: timepoints,
		Horizon:    30,
		Features: ml.Features{
			AvgDailyExpense: avgDailyExpense,
			IncomeEvents:    incomeEvents,
		},
	}

	forecast, err := s.ml.Predict(req)
	if err != nil {
		fail(err)
		return
	}

	s.forecastRepo.Update(domain.ForecastJob{ID: id, Status: "done", Result: &forecast})
}

func nextRuleOccurrence(from time.Time, dayOfMonth int) time.Time {
	y, m, _ := from.Date()
	loc := from.Location()
	candidate := time.Date(y, m, dayOfMonth, 0, 0, 0, 0, loc)
	if candidate.After(from) {
		return candidate
	}
	if m == time.December {
		return time.Date(y+1, time.January, dayOfMonth, 0, 0, 0, 0, loc)
	}
	return time.Date(y, m+1, dayOfMonth, 0, 0, 0, 0, loc)
}

func (s *service) UploadTransactions(contentType string, file io.Reader) error {
	pars, err := s.parserFactory.Create(contentType)
	if err != nil {
		return err
	}
	trans, err := pars.Parse(file)
	if err != nil {
		return err
	}
	err = s.txRepo.SaveAll(trans)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTransactions(filters domain.Filters) ([]domain.Transaction, error) {
	trans, err := s.txRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	return trans, nil
}

func (s *service) UpdateTransaction(tx domain.Transaction) error {
	err := s.txRepo.Update(tx)
	return err
}

func (s *service) DeleteTransaction(id int) error {
	err := s.txRepo.Delete(id)
	return err
}

func (s *service) CreateForecastJob() (int, error) {
	id, err := s.forecastRepo.Create()
	if err != nil {
		return 0, err
	}
	go s.forecastWorker(id)
	return id, nil
}

func (s *service) GetForecastJob(id int) (domain.ForecastJob, error) {
	job, err := s.forecastRepo.FindByID(id)
	if err != nil {
		return domain.ForecastJob{}, err
	}

	return job, nil
}

func (s *service) GetForecastJobs() ([]domain.ForecastJob, error) {
	jobs, err := s.forecastRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *service) SaveRecurringRule(rule domain.RecurringRule) error {
	err := s.recurRepo.Save(rule)
	return err
}

func (s *service) GetRecurringRules() ([]domain.RecurringRule, error) {
	rules, err := s.recurRepo.FindAll()
	if err != nil {
		return nil, err
	}

	rulesById := make(map[string][]domain.RecurringRule)
	past := []domain.RecurringRule{}
	future := []domain.RecurringRule{}
	today := time.Now()
	result := []domain.RecurringRule{}

	for _, rule := range rules {
		if rule.StartDate.Before(today) {
			past = append(past, rule)
		} else {
			future = append(future, rule)
		}
	}

	result = append(result, future...)

	for _, rule := range past {
		rulesById[rule.SourceID] = append(rulesById[rule.SourceID], rule)
	}

	for _, rules := range rulesById {
		newest := rules[0]
		for _, r := range rules[1:] {
			if r.StartDate.After(newest.StartDate) {
				newest = r
			}
		}
		result = append(result, newest)
	}

	return result, nil
}

func (s *service) DeleteRecurringRule(sourceID string) error {
	err := s.recurRepo.Delete(sourceID)
	return err
}
