from src.domain.field_validation_errors import FieldError
from src.exceptions.base import AppException


class SetupFormDataValidationError(AppException):
    def __init__(self, plan_id: str, details: dict[str, list[FieldError]]):
        json_details = {
            field: [err.model_dump() for err in error_list]
            for field, error_list in details.items()
        }
        super().__init__(
            message=f"Setup form data for plan '{plan_id}' is invalid. See details for more information.",
            status_code=400,
            details=json_details,
        )
