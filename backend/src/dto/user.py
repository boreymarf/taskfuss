from pydantic import BaseModel, ConfigDict


class UserCreate(BaseModel):
    login: str
    password: str

    model_config = ConfigDict(from_attributes=True)


class UserLogin(BaseModel):
    login: str
    password: str

    model_config = ConfigDict(from_attributes=True)
