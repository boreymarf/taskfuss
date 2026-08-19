import logging
from typing import Any, Literal, cast

from pydantic import TypeAdapter

from src.domain.field_validation_errors import FieldError
from src.domain.fields import Field, TupleField
from src.exceptions.generic import NotFoundError
logger = logging.getLogger(__name__)

class FormProcessorNoDataError(Exception):
    def __init__(self):
        self.message = (
            "Tried to access data in form processor without loading it first!"
        )
        super().__init__(self.message)


class FormProcessorNoFieldsError(Exception):
    def __init__(self):
        self.message = (
            "Tried to access fields in form processor without loading them first!"
        )
        super().__init__(self.message)


class FormProcessorNotFound(Exception):
    def __init__(
        self,
        type: Literal["field", "value"],
        path: str,
        message: str | None = None,
        *,
        current_dir: Any | None = None,
    ):
        self.type = type
        self.path = path
        self.message = message or f"A {type} on a path '{path}' was not found!"
        self.current_dir = current_dir
        super().__init__(self.message)


class FormProcessorIncorrectType(Exception):
    def __init__(
        self,
        path: str,
        current_type: str,
        expected_type: str,
        message: str | None = None,
    ):
        self.path = path
        self.current_type = current_type
        self.expected_type = expected_type
        self.message = message or f"Wrong type on a path '{path}'"
        super().__init__(self.message)


class FormProcessor:
    def __init__(
        self, fields: dict[str, Field] | None = None, data: dict[str, Any] | None = None
    ):
        self.fields = fields or {}
        self.data = data or {}

    def new_data(self) -> None:
        self.data = {}

    def load_fields(self, fields: dict[str, Field]) -> None:
        # Without adapter, this methdo will allow plain JSON to pass
        adapter = TypeAdapter(dict[str, Field])
        self.fields = adapter.validate_python(fields)

    def load_data(self, data: dict[str, Any]) -> None:
        self.data = data

    def validate_data(self) -> list[FieldError]:

        if not self.data:
            raise FormProcessorNoDataError

        if not self.fields:
            raise FormProcessorNoFieldsError

        error_list: list[FieldError] = []
        for i, field in self.fields.items():
            field.validate_value(self.data[i], i, error_list=error_list)

        return error_list

    def validate_value(
        self, *args: str | int, value: Any, error_list: list[FieldError] | None = None
    ) -> list[FieldError]:

        if not self.fields:
            raise FormProcessorNoFieldsError

        field = self.get_field(*args)
        return field.validate_value(value, ".".join(str(*args)), error_list=error_list)

    def get_field(self, *args: str | int) -> Field:

        if not self.fields:
            raise FormProcessorNoFieldsError

        path_str = ".".join(map(str, args))
        path = path_str.split(".")

        current_dir: dict[str, Field] | Field = self.fields
        current_path: list[str] = []

        for p in path:

            current_path.append(p)

            # If it's a tuple
            if isinstance(current_dir, TupleField) and p.isdigit():
                try:
                    current_dir = current_dir.get_field(int(p))
                except NotFoundError:
                    raise NotFoundError("field", ".".join(current_path))
                continue

            # If it's a dict (only a first layer)
            # Should be last or things will break
            elif isinstance(current_dir, dict):
                current_dir = current_dir[p]
                continue

            # Default
            else:
                raise FormProcessorNotFound(
                    "field", ".".join(current_path), current_dir=current_dir
                )

        assert isinstance(current_dir, Field)
        return current_dir

    def get_default(self, *args: str | int) -> Any:

        if not self.fields:
            raise FormProcessorNoFieldsError

        return self.get_field(*args).get_default()

    def get_value(self, *args: str | int) -> Any | None:

        if not self.data:
            raise FormProcessorNoDataError

        path_str = ".".join(map(str, args))
        path = path_str.split(".")

        current_dir: dict[str, Any] | list[Any] = self.data
        current_path: list[str] = []

        for p in path:

            current_path.append(p)
            current_dir = cast(dict[str, Any] | list[Any], current_dir)

            if isinstance(current_dir, dict):
                current_dir = current_dir[p]
            elif isinstance(current_dir, list) and p.isdigit():
                current_dir = current_dir[int(p)]
            else:
                raise FormProcessorNotFound(
                    "field", ".".join(current_path), current_dir=current_dir
                )

        return current_dir

    def get_value_or_default(self, *args: str | int) -> Any:

        if not self.data:
            raise FormProcessorNoDataError

        if not self.fields:
            raise FormProcessorNoFieldsError

        try:
            return self.get_value(*args)
        except FormProcessorNotFound:
            return self.get_default(*args)

    def insert_value(self, value: Any, *args: str | int) -> None:

        path_str = ".".join(map(str, args))
        path = path_str.split(".")

        current_dir: dict[str, Any] | list[Any] = self.data
        current_path: list[str] = []

        for i, p in enumerate(path):
            current_path.append(p)
            is_last = i == len(path) - 1
            next_is_list = not is_last and path[i + 1].isdigit()

            # === Последний сегмент: записываем значение ===
            if is_last:
                if isinstance(current_dir, dict):
                    current_dir[p] = value
                elif isinstance(current_dir, list) and p.isdigit():
                    idx = int(p)
                    current_list = cast(list[Any], current_dir)
                    while len(current_list) <= idx:
                        current_list.append(None)
                    current_list[idx] = value
                else:
                    raise FormProcessorIncorrectType(
                        path=".".join(current_path),
                        current_type=type(current_dir).__name__,
                        expected_type="dict" if isinstance(current_dir, list) else "dict/list",
                    )
                return

            # === Промежуточный сегмент: получаем или создаём контейнер ===
            expected_type = list if next_is_list else dict

            if isinstance(current_dir, dict):
                if p not in current_dir or current_dir[p] is None:
                    current_dir[p] = list[Any]() if next_is_list else dict[str, Any]()
                target = current_dir[p]
            elif isinstance(current_dir, list) and p.isdigit():
                idx = int(p)
                current_list = cast(list[Any], current_dir)
                while len(current_list) <= idx:
                    current_list.append(None)
                if current_list[idx] is None:
                    current_list[idx] = list[Any]() if next_is_list else dict[str, Any]()
                target = current_list[idx]
            else:
                raise FormProcessorIncorrectType(
                    path=".".join(current_path),
                    current_type=type(current_dir).__name__,
                    expected_type="dict" if isinstance(current_dir, list) else "dict/list",
                )

            # Проверяем, что target имеет ожидаемый тип
            if not isinstance(target, expected_type):
                raise FormProcessorIncorrectType(
                    path=".".join(current_path),
                    current_type=type(target).__name__,
                    expected_type=expected_type.__name__,
                )

            current_dir = cast(dict[str, Any] | list[Any], target)
