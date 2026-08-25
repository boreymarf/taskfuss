from sqlalchemy import select
from sqlalchemy.orm import Session
from sqlalchemy.exc import IntegrityError
import logging

from src.db.user import UserDB
from src.domain.user import User, UserCreate
from src.exceptions.generic import NotFoundError
from src.exceptions.user import UserAlreadyExistsError
from src.security.hashing import hash_password

logger = logging.getLogger(__name__)


class UserRepository:

    def create_user(
        self,
        db: Session,
        data: UserCreate,
        *,
        overwrite_id: int | None = None,
    ) -> User:
        user_data = data.model_dump(exclude={"password"})
        user_data["password_hash"] = hash_password(data.password)
        if overwrite_id is not None:
            user_data["id"] = overwrite_id

        try:
            user_db = UserDB(**user_data)
            db.add(user_db)
            db.flush()
        except IntegrityError:
            db.rollback()
            raise UserAlreadyExistsError(data.login)

        return User.model_validate(user_db)

    def get_user(self, db: Session, user_id: int) -> User:
        user_db = db.get(UserDB, user_id)
        if user_db is None:
            raise NotFoundError("User", user_id)
        return User.model_validate(user_db)

    def get_all_users(self, db: Session) -> list[User]:
        users_db = db.scalars(select(UserDB)).all()
        return [User.model_validate(user_db) for user_db in users_db]

    def get_user_by_login(self, db: Session, login: str) -> User | None:
        user_db = db.scalar(select(UserDB).where(UserDB.login == login))
        if user_db:
            return User.model_validate(user_db)
        return None
