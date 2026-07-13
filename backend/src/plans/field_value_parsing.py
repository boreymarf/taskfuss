
import logging
from typing import Any

from src.exceptions import PlanSetupValidationError
from src.plans.fields import Field, StrField
from src.plans.value_types import ValueType



logger = logging.getLogger(__name__)

def parse_str(value: str, field: StrField) -> list[str]:
    try:
        str(value)
    except

def parse_value(field_id: str, value: Any, field: Field):
    if type(field) is StrField:
        parse_str(field_id, value, field)
    else:
        logger.error(f"There's no parse function for field type {field}!")
        raise PlanSetupValidationError(f"Invalid field type for {field_id}!")

