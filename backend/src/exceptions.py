from __future__ import annotations
from uuid import UUID
from typing import Any, TypeAlias


class AppException(Exception):
    """Base exception."""

    def __init__(
        self,
        message: str,
        status_code: int = 400,
        details: dict[str, Any] | None = None,
    ):
        super().__init__(message)
        self.message = message
        self.status_code = status_code
        self.details = details or {}


# Generic


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


# User


class InvalidToken(AppException):
    def __init__(self, details: dict[str, Any] | None = None):
        super().__init__(message="Invalid token.", status_code=401, details=details)


class UserAlreadyExistsError(AppException):
    def __init__(self, login: str, details: dict[str, Any] | None = None):
        super().__init__(
            message=f"User '{login}' already exists",
            status_code=409,
            details={"login": login, **(details or {})},
        )


class UserNotFound(AppException):
    def __init__(self, user_id: int, details: dict[str, Any] | None = None):
        super().__init__(
            message=f"User with id {user_id} does not exist.",
            status_code=422,
            details={"user_id": user_id, **(details or {})},
        )


class InvalidCredentialsError(AppException):
    def __init__(self, details: dict[str, Any] | None = None):
        super().__init__(
            message="Invalid credentials.", status_code=401, details=details
        )


# Quest

PlanSetupFormDetail: TypeAlias = dict[str, list[str] | 'PlanSetupFormDetail']
class PlanSetupFormValidationError(AppException):
    def __init__(
        self,
        message: str = "Failed to validate setup form data",
        details: PlanSetupFormDetail = {},
        status_code: int = 400,
    ):
        super().__init__(
            message=message, status_code=status_code, details=details or {}
        )


class PlanSetupValidationError(AppException):
    def __init__(
        self,
        message: str = "Validation error",
        details: dict[str, Any] | None = None,
        status_code: int = 400,
    ):
        super().__init__(
            message=message, status_code=status_code, details=details or {}
        )
