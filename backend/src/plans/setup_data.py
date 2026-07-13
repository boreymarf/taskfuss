from typing import Any

from pydantic import BaseModel, ConfigDict





class SetupData(BaseModel):
    """All info gathered during setup"""

    form_data: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)
