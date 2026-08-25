from pydantic import BaseModel, ConfigDict


class User(BaseModel):
    """Domain-level model to create a user."""

    id: int
    login: str
    password_hash: str

    model_config = ConfigDict(from_attributes=True)


class UserCreate(BaseModel):
    login: str
    password: str

    model_config = ConfigDict(from_attributes=True)


class UserLogin(BaseModel):
    login: str
    password: str

    model_config = ConfigDict(from_attributes=True)
