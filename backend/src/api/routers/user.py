
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.domain.user import User
from src.service.user import UserService


router = APIRouter(prefix="/api/users", tags=["users"])

@router.get("/me", response_model=User)
def get_current_user_profile(
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
) -> User:
    """Get current authenticated user profile."""
    return UserService.get_user(db, current_user)
