import logging
from typing import Any, cast

from src.domain.field_validation_errors import FieldError
from src.domain.fields import FormFields
from src.exceptions.generic import NotFoundError


logger = logging.getLogger(__name__)

# NOTE: Vibe coded :(
def get_nested_value(data: dict[str, Any], path: str) -> Any | None:
    parts = path.split(".")
    current: Any = data
    for part in parts:
        if current is None:
            return None
        if part.isdigit():
            idx = int(part)
            if not isinstance(current, list):
                return None
            current_list = cast(list[Any], current)
            if idx < 0 or idx >= len(current_list):
                return None
            current = current_list[idx]
        else:
            if not isinstance(current, dict):
                return None
            current_dict = cast(dict[str, Any], current)
            if part not in current_dict:
                return None
            current = current_dict[part]
    return current


class FormProcessor:
    def __init__(self, fields: FormFields | None = None, data: dict[str, Any] | None = None):
        self.fields = fields or {}
        self.data = data or {}

    def load_fields(self, fields: FormFields) -> None:
        self.fields = fields

    def load_data(self, data: dict[str, Any]) -> None:
        self.data = data

    def validate(self) -> list[FieldError]:
        """Валидирует все поля и возвращает плоский список ошибок с точными loc."""
        errors: list[FieldError] = []
        for field_name, field_def in self.fields.items():
            value = self.data.get(field_name)
            field_def.validate_value(value, loc=field_name, error_list=errors)
        return errors

    def validate_value(self, *args: str | int, value: Any, error_list: list[FieldError] | None = None) -> list[FieldError]:
        if len(args) == 1 and isinstance(args[0], str):
            path = args[0]
        else:
            path = ".".join(str(part) for part in args)

        logger.debug(f"Looking for field: '{path}' in fields: {list(self.fields.keys())}")
        logger.debug(f"Full fields dict: {self.fields}")

        field_def = self.fields.get(path)
        if field_def is None:
            raise NotFoundError("Field", path)

        if error_list is None:
            local_errors: list[FieldError] = []
            field_def.validate_value(value, loc=path, error_list=local_errors)
            return local_errors
        else:
            field_def.validate_value(value, loc=path, error_list=error_list)
            return error_list
