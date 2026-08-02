from typing import Any

from starlette.status import HTTP_400_BAD_REQUEST

from src.domain.field_validation_errors import FieldError
from src.exceptions import AppException


class RecordValidationFailed(AppException):
    def __init__(self, value: Any, path: str, errors: list[FieldError]):
        super().__init__(
            message=f"Quest state with '{id}' doesn't have any fields!",
            status_code=HTTP_400_BAD_REQUEST,
            details={"errors": errors},
        )
