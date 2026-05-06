from fastapi import APIRouter
from fastapi.responses import JSONResponse

from app.schemas.request import PredictRequest
from app.schemas.response import PredictResponse, ForecastPoint, ErrorDetail, ErrorResponse
from app.models.forecaster import Forecaster

router = APIRouter()
forecaster = Forecaster()


@router.post("/internal/v1/predict", response_model=PredictResponse)
async def predict(req: PredictRequest):
    try:
        forecast_raw, predicted_balance, confidence = forecaster.predict(
            req.timeseries, req.horizon, req.features
        )
    except Exception as e:
        return JSONResponse(
            status_code=500,
            content=ErrorResponse(
                error=ErrorDetail(code="PREDICTION_ERROR", message=str(e))
            ).model_dump(),
        )

    return PredictResponse(
        forecast=[ForecastPoint(t=p["t"], balance=p["balance"]) for p in forecast_raw],
        predicted_balance=predicted_balance,
        confidence=confidence,
    )
