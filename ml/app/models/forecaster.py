import pandas as pd
import numpy as np
from prophet import Prophet


class Forecaster:
    def predict(self, timeseries, horizon, features):
        base_date = pd.Timestamp.today().normalize()
        n = len(timeseries)

        df = pd.DataFrame([
            {
                "ds": base_date - pd.Timedelta(days=n - p.t),
                "y": p.balance,
            }
            for p in timeseries
        ])

        model = Prophet(daily_seasonality=False, weekly_seasonality=True)
        model.fit(df)

        future = model.make_future_dataframe(periods=horizon)
        forecast_df = model.predict(future)

        future_rows = forecast_df.tail(horizon)
        last_t = int(timeseries[-1].t)

        # apply expected income events as post-processing
        income_map = {e.t: e.amount for e in features.income_events}
        predicted = future_rows["yhat"].values.copy()
        for i in range(horizon):
            t = last_t + i + 1
            if t in income_map:
                predicted[i] += income_map[t]

        forecast = [
            {"t": last_t + i + 1, "balance": float(b)}
            for i, b in enumerate(predicted)
        ]

        # confidence: based on coefficient of variation of the interval width
        interval_width = (future_rows["yhat_upper"] - future_rows["yhat_lower"]).mean()
        avg_abs_y = future_rows["yhat"].abs().mean()
        scale = max(avg_abs_y, interval_width, 1.0)
        confidence = float(np.clip(1.0 - interval_width / scale / 2, 0.0, 1.0))

        return forecast, float(predicted[-1]), confidence
