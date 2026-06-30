from typing import Literal

from pydantic import BaseModel, ConfigDict
from enum import Enum


class ValueType(str, Enum):
    STR = "str"
    LIST_OBJECT = "list[object]"


class ListField(BaseModel):
    discriminator: Literal["list"] = "list"
    value_type: Literal[ValueType.LIST_OBJECT] = ValueType.LIST_OBJECT
    label: str | None = None
    description: str | None = None
    item_field: Field
    min_items: int | None = None
    max_items: int | None = None

    model_config = ConfigDict(from_attributes=True)


class StrField(BaseModel):
    discriminator: Literal["str"] = "str"
    value_type: Literal[ValueType.STR] = ValueType.STR
    label: str | None = None
    description: str | None = None
    required: bool = False
    default: str | None = None
    max_length: int | None = None
    min_length: int | None = None

    model_config = ConfigDict(from_attributes=True)


Field = StrField | ListField
