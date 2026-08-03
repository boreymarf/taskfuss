from pydantic import BaseModel, ConfigDict
from typing import Literal


class RequiredError(BaseModel):
    loc: str
    discriminator: Literal["required"] = "required"
    message: str = "It's required to return a value"
    model_config = ConfigDict(from_attributes=True)


class IncorrectTypeError(BaseModel):
    loc: str
    discriminator: Literal["incorrect_type_error"] = "incorrect_type_error"
    message: str = "The value returned is incorrect for the field."
    current_type: str
    correct_type: str
    model_config = ConfigDict(from_attributes=True)


class MinSizeError(BaseModel):
    loc: str
    discriminator: Literal["min_size_error"] = "min_size_error"
    message: str = "The value is too low."
    current_value: int
    min_value: int
    model_config = ConfigDict(from_attributes=True)


class MaxSizeError(BaseModel):
    loc: str
    discriminator: Literal["max_size_error"] = "max_size_error"
    message: str = "The value is too high."
    current_value: int
    max_value: int
    model_config = ConfigDict(from_attributes=True)


class ListErrors(BaseModel):
    loc: str
    discriminator: Literal["list_errors"] = "list_errors"
    message: str = "There's one or more errors in the list."
    errors: dict[int, list["FieldError"]]
    model_config = ConfigDict(from_attributes=True)


class CustomError(BaseModel):
    loc: str
    discriminator: Literal["custom"] = "custom"
    message: str = "Custom error message"
    model_config = ConfigDict(from_attributes=True)


class ExactLengthError(BaseModel):
    loc: str
    discriminator: Literal["exact_length_error"] = "exact_length_error"
    message: str = "The tuple has an incorrect number of elements."
    current_length: int
    required_length: int
    model_config = ConfigDict(from_attributes=True)


FieldError = (
    RequiredError
    | IncorrectTypeError
    | MinSizeError
    | MaxSizeError
    | ListErrors
    | CustomError
    | ExactLengthError
)

# This fixes recursion error in openapi docs generation
ListErrors.model_rebuild()
