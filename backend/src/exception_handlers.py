import logging
import traceback

from fastapi import Request
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
from pydantic import ValidationError

from src.exceptions.base import AppException

logger = logging.getLogger(__name__)


async def app_exception_handler(_request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"message": exc.message, "details": exc.details},
    )


async def validation_exception_handler(_request: Request, exc: ValidationError):
    return JSONResponse(
        status_code=422, content={"detail": jsonable_encoder(exc.errors())}
    )


async def generic_exception_handler(request: Request, exc: Exception):
    tb_str = traceback.format_exc()
    request_info = {
        "method": request.method,
        "url": str(request.url),
        "client": request.client.host if request.client else "unknown",
        "headers": dict(request.headers),
    }
    logger.error(
        "Unhandled exception: %s\nRequest: %s\nTraceback:\n%s",
        exc,
        request_info,
        tb_str,
        exc_info=True,
    )
    return JSONResponse(status_code=500, content={"detail": "Internal server error"})
