from pydantic import BaseModel, field_validator
from typing import List
import math


class TimePoint(BaseModel):
    t: int
    balance: float
    day_of_week: int
    is_weekend: bool
    food_total: float
    transport_total: float
    entertainment_total: float
    avg_transaction_size: float
    transaction_count: int

    @field_validator("balance", "food_total", "transport_total", "entertainment_total",
                     "avg_transaction_size", mode="after")
    @classmethod
    def no_nan_inf(cls, v: float) -> float:
        if math.isnan(v) or math.isinf(v):
            raise ValueError("value must be finite")
        return v


class IncomeEvent(BaseModel):
    t: int
    amount: float
    label: str = ""


class Features(BaseModel):
    avg_daily_expense: float
    income_events: List[IncomeEvent] = []


class PredictRequest(BaseModel):
    timeseries: List[TimePoint]
    horizon: int
    features: Features

    @field_validator("timeseries")
    @classmethod
    def timeseries_min_length(cls, v):
        if len(v) < 3:
            raise ValueError("timeseries must have at least 3 points")
        return v

    @field_validator("horizon")
    @classmethod
    def horizon_valid(cls, v):
        if v < 1 or v > 365:
            raise ValueError("horizon must be between 1 and 365")
        return v
