from datetime import timedelta
import logging
from sqlalchemy import select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

from src.db import UserDB
from src.domain.auth import Token
from src.domain.user import User
from src.dto.user import UserCreate, UserLogin
from src.exceptions import (
    InvalidCredentialsError,
    NotFoundError,
    UserAlreadyExistsError,
)
from src.security.auth import create_access_token
from src.security.hashing import hash_password, verify_password


logger = logging.getLogger(__name__)


class UserService:
    @staticmethod
    def create_user(
        db: Session,
        data: UserCreate,
        *,
        overwrite_id: int | None = None,
        auto_commit: bool = False,
    ) -> User:
        logger.debug("Creating user with login=%s", data.login)

        # Build DB model
        user_data = data.model_dump(exclude={"password"})
        user_data["password_hash"] = hash_password(data.password)

        # Flags
        if overwrite_id is not None:
            user_data["id"] = overwrite_id

        try:
            # Create
            user_db = UserDB(**user_data)
            db.add(user_db)
            db.flush()
        except IntegrityError:
            db.rollback()
            logger.warning("User with login=%s already exists", data.login)
            raise UserAlreadyExistsError(data.login)

        if auto_commit:
            db.commit()

        logger.info(
            "User created successfully: id=%s, login=%s", user_db.id, data.login
        )

        # Translate to domain model
        return User.model_validate(user_db)

    @staticmethod
    def get_user(db: Session, user_id: int) -> User:
        logger.debug("Fetching user with id=%s", user_id)

        # Fetch
        user_db = db.get(UserDB, user_id)

        # Checks
        if user_db is None:
            logger.warning("User not found: id=%s", user_id)
            raise NotFoundError("User", user_id)

        logger.info("User fetched successfully: id=%s", user_id)

        # Translate to domain model
        return User.model_validate(user_db)

    @staticmethod
    def get_all_users(db: Session) -> list[User]:
        logger.debug("Fetching all users")

        # Fetch
        users_db = db.scalars(select(UserDB)).all()

        logger.info("Fetched %d users successfully", len(users_db))

        # Translate to domain models
        return [User.model_validate(user_db) for user_db in users_db]

    @staticmethod
    def authenticate_user(
        db: Session,
        data: UserLogin,
        *,
        token_expire_minutes: int | None = None,
    ) -> Token:
        logger.debug("Authenticating user with login=%s", data.login)

        # Checks
        user_db = db.scalar(select(UserDB).where(UserDB.login == data.login))
        if not user_db:
            logger.warning("Authentication failed: invalid login=%s", data.login)
            raise InvalidCredentialsError

        verify_password(data.password, user_db.password_hash)

        # Create token
        token_data = {"sub": str(user_db.id)}

        # Overwrite
        if token_expire_minutes is not None:
            if token_expire_minutes == 0:
                # Infinite expiration
                access_token = create_access_token(data=token_data, expires_delta=None)
            else:
                access_token = create_access_token(
                    data=token_data,
                    expires_delta=timedelta(minutes=token_expire_minutes),
                )
        # Default expiration
        else:
            access_token = create_access_token(data=token_data)

        logger.info(
            "User authenticated successfully: id=%s, login=%s", user_db.id, data.login
        )

        # Translate to domain model and return
        return Token(
            access_token=access_token,
            token_type="bearer",
            user=User.model_validate(user_db),
        )
