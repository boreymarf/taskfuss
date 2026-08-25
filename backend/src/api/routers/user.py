from typing import Annotated

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from dependency_injector.wiring import inject, Provide

from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import BAD_REQUEST, CONFLICT, UNAUTHORIZED
from src.container import Container
from src.domain.auth import Token
from src.domain.user import User, UserCreate, UserLogin
from src.service.user import UserService

router = APIRouter(prefix="/api/users", tags=["users"])

@router.post("/sign-up", response_model=Token, responses={**CONFLICT})
@inject
def sign_up(
    data: UserCreate,
    user_service: Annotated[UserService, Depends(Provide[Container.user_service])],
    db: Session = Depends(get_session),
):
    user = user_service.create_user(db, data, auto_commit=True)
    token = user_service.generate_token_for_user(user)
    return token


@router.post(
    "/sign-in", response_model=Token, responses={**BAD_REQUEST, **UNAUTHORIZED}
)
@inject
def sign_in(
    data: UserLogin,
    user_service: Annotated[UserService, Depends(Provide[Container.user_service])],
    db: Session = Depends(get_session),
):
    user = user_service.authenticate_user(db, data)
    token = user_service.generate_token_for_user(user)
    return token


@router.get("/me", response_model=User, responses={**UNAUTHORIZED})
@inject
def get_current_user_profile(
    user_service: Annotated[UserService, Depends(Provide[Container.user_service])],
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return user_service.get_user(db, current_user)
