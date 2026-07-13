from enum import Enum


class ValueType(str, Enum):
    STR = "str"
    LIST_OBJECT = "list[object]"
    BOOL = "bool"
