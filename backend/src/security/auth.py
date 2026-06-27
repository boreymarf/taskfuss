from datetime import datetime, timedelta, timezone
import logging
from typing import Any
from jose import JWTError, jwt
from src.config import get_config
from src.exceptions import InvalidCredentialsError

ALGORITHM = "HS256"

logger = logging.getLogger(__name__)


def create_access_token(
    data: dict[str, Any], expires_delta: timedelta | None = None
) -> str:
    to_encode = data.copy()
    config = get_config()
    secret_key = config.security.auth_secret_key

    if secret_key == "change-me" and config.app.environment != "dev":
        raise RuntimeError(
            "JWT authentication is not properly configured: "
            + "default secret key used in non-dev environment. "
            + "Please set a secure key under [backend.security]."
        )

    expire = datetime.now(timezone.utc) + (
        expires_delta or timedelta(minutes=config.security.auth_token_lifespan_minutes)
    )
    to_encode["exp"] = expire
    return jwt.encode(to_encode, secret_key, algorithm=ALGORITHM)


def verify_token(token: str) -> dict[str, Any]:
    secret_key = get_config().security.auth_secret_key
    try:
        payload: dict[str, Any] = jwt.decode(token, secret_key, algorithms=[ALGORITHM])
        return payload
    except JWTError:
        raise InvalidCredentialsError()
