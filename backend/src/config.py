import logging
import os
from pathlib import Path
import re
from typing import Literal
from pydantic import BaseModel, Field, field_validator

from src.classes.config import ConfigManager
from src.logging.settings import LoggingConfig

logger = logging.getLogger(__name__)


class AppDevConfig(BaseModel):
    reload_app: bool = True
    reload_dirs: list[str] = ["src"]


class AppConfig(BaseModel):
    environment: Literal["dev", "prod"] = "dev"
    dev: AppDevConfig = Field(default_factory=AppDevConfig)


class ServerConfig(BaseModel):
    port: int = 5000
    allow_origins: list[str] = [
        "http://localhost:6006",
        "http://localhost:5173",
    ]  # Add your localhosts
    allow_credentials: bool = True
    allow_methods: list[str] = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers: list[str] = ["Authorization", "Content-Type"]

    @field_validator("allow_origins", mode="after")
    @classmethod
    def validate_origin_urls(cls, values: list[str]) -> list[str]:
        URL_PATTERN = re.compile(
            r"^(https?://)"  # scheme
            + r"([a-zA-Z0-9.-]+)"  # domain
            + r"(:\d+)?$"  # optional port
        )

        for origin in values:
            if not URL_PATTERN.match(origin):
                logger.warning(
                    f"Invalid origin URL format in the config for allow_origins: '{origin}'. "
                    + "Are you sure you didn't fucked up?"
                )

        return values


class SecurityLoginPolicy(BaseModel):
    min_length: int = 3
    max_length: int = 20
    allow_lowercase: bool = True
    allow_uppercase: bool = True
    allow_digits: bool = True
    allow_underscore: bool = True
    allow_hyphen: bool = False
    allow_dots: bool = False
    allow_cyrillic: bool = False
    allow_special: bool = False


class SecurityPasswordPolicy(BaseModel):
    min_length: int = 6
    max_length: int = 128
    require_uppercase: bool = False
    require_digit: bool = False
    require_special: bool = False


class SecurityConfig(BaseModel):
    auth_secret_key: str = "change-me"
    auth_token_lifespan_seconds: int | Literal["infinite"] = 120 * 60
    login_policy: SecurityLoginPolicy = Field(default_factory=SecurityLoginPolicy)
    password_policy: SecurityPasswordPolicy = Field(
        default_factory=SecurityPasswordPolicy
    )


class DatabaseConfig(BaseModel):
    db_type: Literal["sqlite"] = "sqlite"
    db_path: str = "data/data.db"


class Config(BaseModel):
    """General config"""

    app: AppConfig = Field(default_factory=AppConfig)
    server: ServerConfig = Field(default_factory=ServerConfig)
    security: SecurityConfig = Field(default_factory=SecurityConfig)
    logging: LoggingConfig = Field(default_factory=LoggingConfig)
    database: DatabaseConfig = Field(default_factory=DatabaseConfig)


_config_manager: ConfigManager[Config] | None = None


def init_config(config_path: Path | str):
    global _config_manager
    _config_manager = ConfigManager(
        model_class=Config,
        default_path=Path(config_path),
        auto_create=True,
        auto_update=True,
        logger=logger,
    )


def get_config() -> Config:
    global _config_manager
    if _config_manager is None:
        config_path = os.environ.get("APP_CONFIG_PATH") or "config.toml"
        init_config(Path(config_path))
    assert _config_manager is not None
    return _config_manager.get()


class FrontendConfig(BaseModel):
    """This config is getting send to the frontend"""

    password_policy: SecurityPasswordPolicy
    login_policy: SecurityLoginPolicy


def get_frontend_config(config: Config) -> FrontendConfig:
    return FrontendConfig(
        password_policy=config.security.password_policy,
        login_policy=config.security.login_policy,
    )
