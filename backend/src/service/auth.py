from datetime import datetime, timedelta, timezone
import logging
import humanize
from typing import Literal

from sqlalchemy import select
from sqlalchemy.orm import Session

from src.config import get_config
from src.db import UserDB
from src.domain.auth import Token
from src.domain.user import User, UserLogin
from src.exceptions import InvalidCredentialsError
from src.security.auth import create_access_token
from src.security.hashing import verify_password

logger = logging.getLogger(__name__)

# TODO: Needs cleaning up later
class AuthService:
    @staticmethod
    def authenticate_user(db: Session, data: UserLogin) -> User:
        """Authenticate user by login and password."""
        logger.debug(f"Authenticating user with login={data.login}")

        user_db = db.scalar(select(UserDB).where(UserDB.login == data.login))
        if not user_db:
            logger.warning(f"Authentication failed: invalid login={data.login}")
            raise InvalidCredentialsError

        verify_password(data.password, user_db.password_hash)

        logger.info(
            f"User authenticated successfully: id={user_db.id}, login={data.login}"
        )

        return User.model_validate(user_db)

    @staticmethod
    def generate_token_for_user(
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
