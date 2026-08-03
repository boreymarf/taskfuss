from __future__ import annotations
from typing import Any, Literal, cast, override
from pydantic import BaseModel, ConfigDict

from src.domain.field_validation_errors import (
    ExactLengthError,
    FieldError,
    IncorrectTypeError,
    ListErrors,
    MaxSizeError,
    MinSizeError,
    RequiredError,
)


class BaseField(BaseModel):
    discriminator: str
    automatic: bool = False

    def validate_value(self, _value: Any, _loc: str) -> list[FieldError]:
        raise NotImplementedError


class CheckboxField(BaseField):
    discriminator: Literal["bool"] = "bool"
    label: str | None = None
    description: str | None = None
    default: bool = False
    model_config = ConfigDict(from_attributes=True)

    @override
    def validate_value(self, value: Any, loc: str) -> list[FieldError]:
        if not isinstance(value, bool):
            return [
                IncorrectTypeError(
                    loc=loc,
                    current_type=value.__class__.__name__,
                    correct_type="bool",
                )
            ]
        return []


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
    def validate_value(self, value: Any, loc: str) -> list[FieldError]:
        errors: list[FieldError] = []

        if value is None and self.treat_none_as_default:
            value = self.default

        if not value and self.required:
            errors.append(RequiredError(loc=loc))
            return errors

        if not isinstance(value, str):
            errors.append(
                IncorrectTypeError(
                    loc=loc,
                    current_type=value.__class__.__name__,
                    correct_type="str",
                )
            )
            return errors

        list_length = len(value)
        if self.min_size is not None and list_length < self.min_size:
            errors.append(
                MinSizeError(
                    loc=loc,
                    current_value=list_length,
                    min_value=self.min_size,
                )
            )
        if self.max_size is not None and list_length > self.max_size:
            errors.append(
                MaxSizeError(
                    loc=loc,
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
    item_field: Field
    min_size: int | None = None
    max_size: int | None = None
    model_config = ConfigDict(from_attributes=True)

    @override
    def validate_value(self, value: Any, loc: str) -> list[FieldError]:
        errors: list[FieldError] = []

        if value is None and self.treat_none_as_default:
            value = self.default

        if not isinstance(value, list):
            errors.append(
                IncorrectTypeError(
                    loc=loc,
                    current_type=value.__class__.__name__,
                    correct_type="list",
                )
            )
            return errors

        if not value and self.required:
            errors.append(RequiredError(loc=loc))
            return errors

        value = cast(list[Any], value)
        list_length = len(value)

        if self.min_size is not None and list_length < self.min_size:
            errors.append(
                MinSizeError(
                    loc=loc,
                    current_value=list_length,
                    min_value=self.min_size,
                )
            )
        if self.max_size is not None and list_length > self.max_size:
            errors.append(
                MaxSizeError(
                    loc=loc,
                    current_value=list_length,
                    max_value=self.max_size,
                )
            )

        child_errors: dict[int, list[FieldError]] = {}
        for i, item in enumerate(value):
            item_errs = self.item_field.validate_value(item, loc=f"{loc}.{i}")
            if item_errs:
                child_errors[i] = item_errs

        if child_errors:
            errors.append(ListErrors(loc=loc, errors=child_errors))

        return errors


class TupleField(BaseField):
    """A fixed-size tuple of fields, each with its own schema."""

    discriminator: Literal["tuple"] = "tuple"
    treat_none_as_default: bool = True
    label: str | None = None
    description: str | None = None
    required: bool = False
    fields: tuple[Field, ...]
    model_config = ConfigDict(from_attributes=True)

    @property
    def default(self) -> tuple[Any, ...]:
        return tuple(f.default for f in self.fields)

    @override
    def validate_value(self, value: Any, loc: str) -> list[FieldError]:
        errors: list[FieldError] = []

        if value is None and self.treat_none_as_default:
            value = self.default

        if not isinstance(value, (list, tuple)):
            errors.append(
                IncorrectTypeError(
                    loc=loc,
                    current_type=value.__class__.__name__,
                    correct_type="tuple",
                )
            )
            return errors

        # Yes, this is a tuple
        value = cast(tuple[Any], value)

        if len(value) != len(self.fields):
            errors.append(
                ExactLengthError(
                    loc=loc,
                    current_length=len(value),
                    required_length=len(self.fields),
                )
            )
            return errors

        if not value and self.required:
            errors.append(RequiredError(loc=loc))
            return errors

        child_errors: dict[int, list[FieldError]] = {}
        for i, (field, item) in enumerate(zip(self.fields, value)):
            item_errs = field.validate_value(item, f"{loc}.{i}")
            if item_errs:
                child_errors[i] = item_errs

        if child_errors:
            errors.append(ListErrors(loc=loc, errors=child_errors))

        return errors


Field = CheckboxField | StrField | ListField | TupleField

# This fixes openapi's build failure
# And also weird "name 'Callable' is not defined" error lol
ListField.model_rebuild() 
TupleField.model_rebuild()

FormFields = dict[str, Field]
