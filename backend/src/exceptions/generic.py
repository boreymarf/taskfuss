from uuid import UUID
from typing import Any

from src.exceptions.base import AppException


class NotFoundError(AppException):
    def __init__(self, resource: str, id: int | UUID | str):
        super().__init__(
            message=f"{resource} with id {id} not found",
            status_code=404,
            details={"resource": resource, "id": str(id)},
        )


class ForbiddenError(AppException):
    def __init__(
        self,
        message: str = "You don't have permission to access this resource",
        details: dict[str, Any] | None = None,
    ):
        super().__init__(message=message, status_code=403, details=details)


class BadRequestError(AppException):
    def __init__(
        self, message: str = "Bad request", details: dict[str, Any] | None = None
    ):
        super().__init__(message=message, status_code=400, details=details)


class TodoError(AppException):
    def __init__(
        self,
        message: str = "This feature is not implemented ¯\\_(ツ)_/¯",
        details: dict[str, Any] | None = None,
    ):
        super().__init__(message=message, status_code=501, details=details)
