from pydantic import BaseModel
from typing import List


class ForecastPoint(BaseModel):
    t: int
    balance: float


class PredictResponse(BaseModel):
    forecast: List[ForecastPoint]
    predicted_balance: float
    confidence: float


class ErrorDetail(BaseModel):
    code: str
    message: str


class ErrorResponse(BaseModel):
    error: ErrorDetail
