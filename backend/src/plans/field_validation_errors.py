from __future__ import annotations
from typing import Literal

from pydantic import BaseModel, ConfigDict


class RequiredError(BaseModel):
    discriminator: Literal["required"] = "required"
    message: str = "It's required to return a value"

    model_config = ConfigDict(from_attributes=True)


class IncorrectTypeError(BaseModel):
    discriminator: Literal["incorrect_type_error"] = "incorrect_type_error"
    message: str = "The value returned is incorrect for the field."
    current_type: str
    correct_type: str

    model_config = ConfigDict(from_attributes=True)


class MinSizeError(BaseModel):
    discriminator: Literal["min_size_error"] = "min_size_error"
    message: str = "The value is too low."
    current_value: int
    min_value: int

    model_config = ConfigDict(from_attributes=True)


class MaxSizeError(BaseModel):
    discriminator: Literal["max_size_error"] = "max_size_error"
    message: str = "The value is too high."
    current_value: int
    max_value: int

    model_config = ConfigDict(from_attributes=True)


class ListErrors(BaseModel):
    discrimanator: Literal["list_errors"] = "list_errors"
    message: str = "There's one or more errors in the list."
    errors: dict[int, list[FieldError]]

    model_config = ConfigDict(from_attributes=True)


FieldError = RequiredError | IncorrectTypeError | MinSizeError | MaxSizeError | ListErrors

# This fixes recursion error in openapi docs generation
ListErrors.model_rebuild()
