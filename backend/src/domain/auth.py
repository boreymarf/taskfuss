from datetime import datetime

from pydantic import BaseModel, ConfigDict


class Token(BaseModel):
    """Domain-level token response."""

    access_token: str
    token_type: str

    issued_at: datetime
    expires_at: datetime | None
    duration_formatted: str
    duration_seconds: int

    model_config = ConfigDict(from_attributes=True)
