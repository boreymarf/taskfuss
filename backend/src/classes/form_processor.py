from typing import Any, cast

from src.domain.field_validation_errors import FieldError, ListErrors
from src.domain.fields import FormFields


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
        raw_errors: list[FieldError] = []
        for field_name, field_def in self.fields.items():
            value = self.data.get(field_name)
            raw_errors.extend(field_def.validate_value(value, loc=field_name))

        # Разворачиваем все ListErrors – в итоге каждый FieldError имеет свой loc
        return self._flatten_errors(raw_errors)

    @staticmethod
    def _flatten_errors(errors: list[FieldError]) -> list[FieldError]:
        flat: list[FieldError] = []
        for err in errors:
            if isinstance(err, ListErrors):
                # Все ошибки внутри уже имеют полный путь (например "items.0.name")
                for idx, child_err_list in err.errors.items():
                    flat.extend(FormProcessor._flatten_errors(child_err_list))
            else:
                flat.append(err)
        return flat

    def get_value(self, *args: str | int) -> Any | None:
        if len(args) == 1 and isinstance(args[0], str):
            path = args[0]
        else:
            path = ".".join(str(part) for part in args)
        return get_nested_value(self.data, path)
