import logging
from typing import Any, cast

from src.domain.field_validation_errors import FieldError
from src.domain.fields import FormFields, Field
from src.exceptions.generic import NotFoundError


logger = logging.getLogger(__name__)


# NOTE: I'm gonna be honest, this whole file is vibecoded and is crap
# I hate every single line in this file, but I can't be bothered to rewrite it yet

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
        path = ".".join(map(str, args))

        logger.debug(f"Looking for field: '{path}' in fields: {list(self.fields.keys())}")
        logger.debug(f"Full fields dict: {self.fields}")

        field_def = self.get_field(*args)

        if error_list is None:
            local_errors: list[FieldError] = []
            field_def.validate_value(value, loc=path, error_list=local_errors)
            return local_errors
        else:
            field_def.validate_value(value, loc=path, error_list=error_list)
            return error_list

    def get_field(self, *args: str | int) -> Field:
        """Возвращает определение поля по пути (например, 'user.name' или 'items.0')."""
        path = ".".join(map(str, args))
        if not path:
            raise NotFoundError("Field", path)

        current: Any = self.fields

        for part in path.split("."):
            if hasattr(current, "fields"):
                nested: Any = current.fields

                if isinstance(nested, list):
                    nested_list = cast(list[Any], nested)
                    if not part.isdigit():
                        raise NotFoundError("Field", path)
                    idx = int(part)
                    if idx < 0 or idx >= len(nested_list):
                        raise NotFoundError("Field", path)
                    current = nested_list[idx]

                elif isinstance(nested, tuple):
                    nested_tuple = cast(tuple[Any, ...], nested)
                    if not part.isdigit():
                        raise NotFoundError("Field", path)
                    idx = int(part)
                    if idx < 0 or idx >= len(nested_tuple):
                        raise NotFoundError("Field", path)
                    current = nested_tuple[idx]

                elif isinstance(nested, dict):
                    nested_dict = cast(dict[str, Any], nested)
                    if part not in nested_dict:
                        raise NotFoundError("Field", path)
                    current = nested_dict[part]

                else:
                    raise NotFoundError("Field", path)

            elif isinstance(current, dict):
                # current — это словарь FormFields (dict[str, Field])
                current_dict = cast(dict[str, Any], current)
                if part not in current_dict:
                    raise NotFoundError("Field", path)
                current = current_dict[part]

            else:
                raise NotFoundError("Field", path)

        return cast(Field, current)


    def get_default(self, *args: str | int) -> Any:
        """Возвращает значение по умолчанию для поля по пути."""
        field = self.get_field(*args)
        return field.get_default()

    def get_value(self, *args: str | int) -> Any | None:
        """Возвращает значение из данных по пути (использует get_nested_value)."""
        path = ".".join(map(str, args))
        return get_nested_value(self.data, path)

    def get_value_or_default(self, *args: str | int) -> Any:
        """
        Возвращает значение из данных, если оно есть, иначе — значение по умолчанию.
        Если путь отсутствует в данных, возвращает дефолт.
        """
        path = ".".join(map(str, args))
        value = get_nested_value(self.data, path)
        if value is None:
            # Пытаемся получить дефолт, если поле существует
            try:
                return self.get_default(*args)
            except NotFoundError:
                return None  # или можно вернуть None, если поле не найдено
        return value
