from src.domain.field_validation_errors import FieldError
from src.exceptions.base import AppException


class SetupFormDataValidationError(AppException):
    def __init__(self, plan_id: str, errors: list[FieldError]):
        json_details = [err.model_dump() for err in errors]
        super().__init__(
            message=f"Setup form data for plan '{plan_id}' is invalid. See details for more information.",
            status_code=400,
            details={"errors": json_details},
        )


class UnknownQuestActionError(Exception):
    def __init__(self, discriminator: str):
        self.discriminator = discriminator
        super().__init__(f"Unknown quest action type: '{discriminator}'")
