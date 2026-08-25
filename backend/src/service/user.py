from datetime import datetime, timedelta, timezone
import logging
from typing import Literal

import humanize
from sqlalchemy.orm import Session

from src.config import get_config
from src.domain.auth import Token
from src.domain.user import User, UserCreate, UserLogin
from src.exceptions import (
    NotFoundError,
    UserAlreadyExistsError,
)
from src.exceptions.user import InvalidCredentialsError
from src.repositories.user import UserRepository
from src.security.auth import create_access_token
from src.security.hashing import verify_password

logger = logging.getLogger(__name__)


# TODO: Get back to this later because this is weird looking service
class UserService:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def create_user(
        self,
        db: Session,
        data: UserCreate,
        *,
        overwrite_id: int | None = None,
        auto_commit: bool = False,
    ) -> User:
        try:
            user = self.user_repository.create_user(db, data, overwrite_id=overwrite_id)
        except UserAlreadyExistsError:
            logger.warning(f"Tried to create new user, but user with login={data.login} already exists")
            raise

        if auto_commit:
            db.commit()

        logger.info(f"User created successfully: id={user.id}, login={data.login}")
        return user

    def get_user(self, db: Session, user_id: int) -> User:
        try:
            user = self.user_repository.get_user(db, user_id)
        except NotFoundError:
            logger.warning(f"User not found: id={user_id}")
            raise

        logger.info(f"User fetched successfully: id={user_id}")
        return user

    def get_all_users(self, db: Session) -> list[User]:
        users = self.user_repository.get_all_users(db)

        logger.info(f"Fetched {len(users)} users successfully")
        return users

    def authenticate_user(self, db: Session, data: UserLogin) -> User:
        user = self.user_repository.get_user_by_login(db, data.login)
        if not user:
            logger.warning(f"Authentication failed: invalid login={data.login}")
            raise InvalidCredentialsError

        verify_password(data.password, user.password_hash)

        logger.info(
            f"User authenticated successfully: id={user.id}, login={data.login}"
        )
        return user

    def generate_token_for_user(
        self,
        user: User,
        *,
        token_expire_seconds: int | Literal["infinite"] | None = None,
    ) -> Token:
        config = get_config()
        token_lifespan = token_expire_seconds or config.security.auth_token_lifespan_seconds

        logger.debug(
            f"Generating token for user id={user.id}, auth_token_lifespan_seconds={token_lifespan}"
        )

        token_data = {"sub": str(user.id)}
        issued_at = datetime.now(timezone.utc)

        if token_lifespan == "infinite":
            expires_at = None
        else:  # int
            expires_at = issued_at + timedelta(seconds=token_lifespan)

        # Create the JWT
        access_token = create_access_token(token_data, expires_at=expires_at)

        # Compute duration fields
        if expires_at is None:
            duration_seconds = 0
            duration_formatted = "infinite"
        else:
            delta = expires_at - issued_at
            duration_seconds = int(delta.total_seconds())
            duration_formatted = humanize.precisedelta(delta)

        logger.info(f"Token generated for user id={user.id}, expires_at={expires_at}")

        return Token(
            access_token=access_token,
            token_type="bearer",
            issued_at=issued_at,
            expires_at=expires_at,
            duration_seconds=duration_seconds,
            duration_formatted=duration_formatted,
        )
