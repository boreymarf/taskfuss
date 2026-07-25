from datetime import timedelta
import logging

from sqlalchemy import select
from sqlalchemy.orm import Session

from src.db import UserDB
from src.domain.auth import Token
from src.domain.user import User
from src.dto.user import UserLogin
from src.exceptions import InvalidCredentialsError
from src.security.auth import create_access_token
from src.security.hashing import verify_password

logger = logging.getLogger(__name__)


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
        user: User, *, token_expire_minutes: int | None = None
    ) -> Token:
        """Generate JWT token for authenticated user."""
        logger.debug(f"Generating token for user id={user.id}")

        token_data = {"sub": str(user.id)}

        if token_expire_minutes is not None:
            if token_expire_minutes == 0:
                access_token = create_access_token(data=token_data, expires_delta=None)
            else:
                access_token = create_access_token(
                    data=token_data,
                    expires_delta=timedelta(minutes=token_expire_minutes),
                )
        else:
            access_token = create_access_token(data=token_data)

        logger.info(f"Token generated for user id={user.id}")

        return Token(access_token=access_token, token_type="bearer")
