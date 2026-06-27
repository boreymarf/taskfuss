from fastapi import APIRouter
from src.config import FrontendConfig, get_config, get_frontend_config


router = APIRouter(prefix="/api/config", tags=["config"])


@router.get(
    "/",
    response_model=FrontendConfig,
)
def get_frontend_configuration():
    config = get_config()
    return get_frontend_config(config)
