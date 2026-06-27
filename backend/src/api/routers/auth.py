
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from src.api.openapi_responses import BAD_REQUEST, CONFLICT, UNAUTHORIZED
from src.database import get_session
from src.domain.auth import Token
from src.dto.user import UserCreate, UserLogin
from src.service.auth import AuthService
from src.service.user import UserService


router = APIRouter(prefix="/api/auth", tags=["authentication"])

@router.post(
    "/sign-up",
    response_model=Token,
    responses={
        **CONFLICT,
    },
)
def sign_up(
    data: UserCreate,
    db: Session = Depends(get_session),
):
    domain_data = UserCreate(**data.model_dump())

    user = UserService.create_user(db, domain_data, auto_commit=True)
    token = AuthService.generate_token_for_user(user)

    return token


@router.post(
    "/sign-in",
    response_model=Token,
    responses={
        **BAD_REQUEST,
        **UNAUTHORIZED,
    },
)
def sign_in(
    data: UserLogin,
    db: Session = Depends(get_session),
):
    domain_data = UserLogin(**data.model_dump())

    user = AuthService.authenticate_user(db, domain_data)
    token = AuthService.generate_token_for_user(user)

    return token
