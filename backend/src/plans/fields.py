from __future__ import annotations
from typing import Literal

from pydantic import BaseModel, ConfigDict

from src.plans.value_types import ValueType




class ListField(BaseModel):
    discriminator: Literal["list"] = "list"
    value_type: Literal[ValueType.LIST_OBJECT] = ValueType.LIST_OBJECT
    automatic: bool = False

    label: str | None = None
    description: str | None = None
    item_field: Field
    min_items: int | None = None
    max_items: int | None = None

    model_config = ConfigDict(from_attributes=True)


class StrField(BaseModel):
    discriminator: Literal["str"] = "str"
    value_type: Literal[ValueType.STR] = ValueType.STR
    automatic: bool = False

    label: str | None = None
    description: str | None = None
    required: bool = False
    default: str | None = None
    max_length: int | None = None
    min_length: int | None = None

    model_config = ConfigDict(from_attributes=True)


class CheckboxField(BaseModel):
    discriminator: Literal["bool"] = "bool"
    value_type: Literal[ValueType.BOOL] = ValueType.BOOL
    automatic: bool = False

    label: str | None = None
    description: str | None = None
    default: bool = False

    model_config = ConfigDict(from_attributes=True)


Field = StrField | ListField

# This fixes recursion error in openapi docs generation
ListField.model_rebuild()

FormFields = dict[str, Field]
