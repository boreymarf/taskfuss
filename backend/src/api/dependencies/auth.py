import logging
from fastapi import Depends, HTTPException, Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import JWTError, jwt

from src.config import get_config

logger = logging.getLogger(__name__)

ALGORITHM = "HS256"
security = HTTPBearer(auto_error=False)

async def get_current_user(
    request: Request,
    creds: HTTPAuthorizationCredentials | None = Depends(security)
) -> int:
    """
    FastAPI dependency that validates the Bearer token and returns the user ID.
    Attaches the user ID to `request.state.user_id`.
    """
    if creds is None:
        raise HTTPException(status_code=401, detail="Not authenticated")

    token = creds.credentials
    secret_key = get_config().security.auth_secret_key

    try:
        payload = jwt.decode(token, secret_key, algorithms=[ALGORITHM])
        user_id = payload.get("sub")
        if not user_id:
            raise HTTPException(status_code=401, detail="Invalid token")
    except JWTError:
        logger.warning("JWT decode failed")
        raise HTTPException(status_code=401, detail="Invalid token")

    request.state.user_id = int(user_id)
    return int(user_id)
