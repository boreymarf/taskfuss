from pydantic import BaseModel, ConfigDict


class User(BaseModel):
    """Domain-level model to create a user."""

    id: int
    login: str

    model_config = ConfigDict(from_attributes=True)
