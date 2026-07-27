from __future__ import annotations
from typing import Any, Literal, cast, override

from pydantic import BaseModel, ConfigDict

from src.domain import FieldError, IncorrectTypeError, ListErrors, MaxSizeError, MinSizeError, RequiredError


def validate_form(
    fields: FormFields, data: dict[str, Any]
) -> dict[str, list[FieldError]]:
    errors: dict[str, list[FieldError]] = {}
    for field_name, field_def in fields.items():
        value = data.get(field_name)
        field_errors = field_def.validate_value(value)
        if field_errors:
            errors[field_name] = field_errors
    return errors


class BaseField(BaseModel):
    discriminator: str
    automatic: bool = False

    def validate_value(self, _value: Any) -> list["FieldError"]:
        raise NotImplementedError


class CheckboxField(BaseField):
    discriminator: Literal["bool"] = "bool"

    label: str | None = None
    description: str | None = None
    default: bool = False

    model_config = ConfigDict(from_attributes=True)

    @override
    def validate_value(self, value: Any) -> list[FieldError]:
        errors: list["FieldError"] = []

        # Check type
        if not isinstance(value, bool):
            errors.append(
                IncorrectTypeError(
                    current_type=value.__class__.__name__, correct_type="bool"
                )
            )
            return errors

        return errors


class StrField(BaseField):
    discriminator: Literal["str"] = "str"

    treat_none_as_default: bool = True
    default: str = ""
    label: str | None = None
    description: str | None = None
    required: bool = False
    max_size: int | None = None
    min_size: int | None = None

    model_config = ConfigDict(from_attributes=True)

    @override
    def validate_value(self, value: Any) -> list[FieldError]:
        errors: list["FieldError"] = []

        # Convert none to default
        if value is None and self.treat_none_as_default:
            value = self.default

        # Check for no value
        if not value and self.required:
            errors.append(RequiredError())
            return errors

        # Check type
        if not isinstance(value, str):
            errors.append(
                IncorrectTypeError(
                    current_type=value.__class__.__name__, correct_type="str"
                )
            )
            return errors

        # Check length
        list_length = len(value)

        if self.min_size is not None and list_length < self.min_size:
            errors.append(
                MinSizeError(
                    current_value=list_length,
                    min_value=self.min_size,
                )
            )

        if self.max_size is not None and list_length > self.max_size:
            errors.append(
                MaxSizeError(
                    current_value=list_length,
                    max_value=self.max_size,
                )
            )

        return errors


class ListField(BaseField):
    discriminator: Literal["list"] = "list"

    treat_none_as_default: bool = True
    default: list[Any] = []
    label: str | None = None
    description: str | None = None
    required: bool = False
    item_field: "Field"
    min_size: int | None = None
    max_size: int | None = None

    model_config = ConfigDict(from_attributes=True)

    @override
    def validate_value(self, value: Any) -> list["FieldError"]:
        errors: list["FieldError"] = []
        
        # Convert none to default
        if value is None and self.treat_none_as_default:
            value = self.default

        # Check type
        if not isinstance(value, list):
            errors.append(
                IncorrectTypeError(
                    current_type=value.__class__.__name__, correct_type="list"
                )
            )
            return errors

        # Check for no value
        if not value and self.required:
            errors.append(RequiredError())
            return errors

        # Yep, it's a list
        value = cast(list[Any], value)

        # Check length
        list_length = len(value)

        if self.min_size is not None and list_length < self.min_size:
            errors.append(
                MinSizeError(
                    message="The list is too small",
                    current_value=list_length,
                    min_value=self.min_size,
                )
            )

        if self.max_size is not None and list_length > self.max_size:
            errors.append(
                MaxSizeError(
                    message="The list is too big",
                    current_value=list_length,
                    max_value=self.max_size,
                )
            )

        # Check children via my item field
        list_errors: dict[int, list[FieldError]] = {}
        for i, v in enumerate(value):
            item_errors = self.item_field.validate_value(v)

            if item_errors:
                list_errors[i] = item_errors

        if list_errors:
            errors.append(ListErrors(errors=list_errors))

        return errors


Field = CheckboxField | StrField | ListField
ListField.model_rebuild()  # This fixes recursion error in openapi docs generation
FormFields = dict[str, Field]
