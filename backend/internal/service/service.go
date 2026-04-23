package service

import (
	"io"
	"log/slog"
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
	today := time.Now()

	trans, err := s.txRepo.FindForForecast(time.Time{}, today)
	if err != nil {
		fail(err)
	}

	rules, err := s.recurRepo.FindActive(today)
	if err != nil {
		fail(err)
	}

	_, _ = trans, rules
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
