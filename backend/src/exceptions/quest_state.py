from typing import Any
from uuid import UUID

from src.exceptions.base import AppException


class NoFieldsQuestStateError(AppException):
    def __init__(self, id: UUID, details: dict[str, Any] | None = None):
        super().__init__(
            message=f"Quest state with '{id}' doesn't have any fields!",
            status_code=401,
            details=details,
        )
