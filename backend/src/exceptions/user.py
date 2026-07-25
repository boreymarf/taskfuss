from typing import Any

from src.exceptions.base import AppException


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


class UserNotFoundError(AppException):
    def __init__(self, user_id: int, details: dict[str, Any] | None = None):
        super().__init__(
            message=f"User '{user_id}' does not exist.",
            status_code=422,
            details={"user_id": user_id, **(details or {})},
        )


class InvalidCredentialsError(AppException):
    def __init__(self, details: dict[str, Any] | None = None):
        super().__init__(
            message="Invalid credentials.", status_code=401, details=details
        )
