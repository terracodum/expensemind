from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse

from app.routes.predict import router
from app.schemas.response import ErrorDetail, ErrorResponse

app = FastAPI(title="ExpenseMind ML Service")

app.include_router(router)


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    first_msg = exc.errors()[0]["msg"] if exc.errors() else "invalid input"

    # distinguish insufficient data from generic validation
    if "at least 3 points" in first_msg:
        code = "INSUFFICIENT_DATA"
    else:
        code = "VALIDATION_ERROR"

    return JSONResponse(
        status_code=400,
        content=ErrorResponse(
            error=ErrorDetail(code=code, message=first_msg)
        ).model_dump(),
    )


@app.get("/health")
async def health():
    return {"status": "ok"}
