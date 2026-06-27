import logging
from pathlib import Path
import tomllib
from typing import Callable, Generic, TypeVar, final
from pydantic import BaseModel, ValidationError
import tomli_w


class ConfigError(Exception):
    """Base exception for all configuration-related errors."""


class ConfigNotFoundError(ConfigError):
    """Raised when the configuration file does not exist."""


class ConfigInvalidError(ConfigError):
    """Raised when the configuration file contains invalid TOML or fails
    Pydantic validation.
    """


@final
class ConfigLogger:
    def __init__(self, backend: logging.Logger | Callable[[str], None] | None = None):
        self._backend = backend

    def info(self, msg: str) -> None:
        if self._backend is None:
            return
        elif isinstance(self._backend, logging.Logger):
            self._backend.info(msg)
        elif callable(self._backend):
            self._backend(msg)

    def debug(self, msg: str) -> None:
        if self._backend is None:
            return
        elif isinstance(self._backend, logging.Logger):
            self._backend.debug(msg)
        elif callable(self._backend):
            self._backend(msg)

    def warning(self, msg: str) -> None:
        if self._backend is None:
            return
        elif isinstance(self._backend, logging.Logger):
            self._backend.warning(msg)
        elif callable(self._backend):
            self._backend(msg)

    def error(self, msg: str) -> None:
        if self._backend is None:
            return
        elif isinstance(self._backend, logging.Logger):
            self._backend.error(msg)
        elif callable(self._backend):
            self._backend(msg)


T = TypeVar("T", bound=BaseModel)


@final
class ConfigManager(Generic[T]):
    def __init__(
        self,
        model_class: type[T],
        default_path: Path | None = None,
        *,
        auto_create: bool = True,
        auto_update: bool = True,
        logger: Callable[[str], None] | logging.Logger | None = None,
    ):
        self.model_class = model_class
        self.default_path = Path(default_path or "config.toml")
        self.auto_create = auto_create
        self.auto_update = auto_update
        self._logger = ConfigLogger(logger)
        self._config: T | None = None
        self._config_path: Path | None = None
        self._config_mtime: float | None = None

    def load_raw(self, path: Path) -> dict[str, object]:
        try:
            with open(path, "rb") as f:
                return tomllib.load(f)
        except FileNotFoundError as e:
            raise ConfigNotFoundError(f"File not found: {path}") from e
        except tomllib.TOMLDecodeError as e:
            raise ConfigInvalidError(f"Invalid TOML: {e}") from e

    def validate(self, data: dict[str, object]) -> T:
        try:
            return self.model_class.model_validate(data)
        except ValidationError as e:
            raise ConfigInvalidError(f"Validation failed: {e}") from e

    def save(self, config: T, path: Path) -> None:
        try:
            with open(path, "wb") as f:
                tomli_w.dump(config.model_dump(), f)
        except OSError as e:
            raise ConfigError(f"Cannot write config: {e}") from e

    def get(self, path: Path | str | None = None) -> T:
        """Get config from cache if valid, otherwise load from disk."""
        path = Path(path or self.default_path)

        if self._is_cached_config_valid(path):
            assert self._config is not None
            return self._config

        # Load fresh config and update cache
        config = self.load_or_create(path)
        self._config = config
        self._config_path = path
        try:
            self._config_mtime = path.stat().st_mtime
        except OSError:
            self._config_mtime = None
        return config

    def reload(self, path: Path | str | None = None) -> T:
        """Force reload config from disk and update cache."""
        path = Path(path or self.default_path)

        self._config = None
        self._config_path = None
        self._config_mtime = None

        return self.get(path)

    def load_or_create(self, path: Path | str | None = None) -> T:
        """Load config from disk, creating if missing and auto_update enabled."""
        path = Path(path or self.default_path)

        if not path.exists():
            if self.auto_create:
                self._logger.info(f"Config not found, creating default at {path}")

                default = self.model_class()

                self.save(default, path)

                self._config = default
                self._config_path = path
                try:
                    self._config_mtime = path.stat().st_mtime
                except OSError:
                    self._config_mtime = None
                return default
            raise ConfigNotFoundError(f"Config file missing: {path}")

        data = self.load_raw(path)
        config = self.validate(data)

        if self.auto_update:
            current_dict = config.model_dump()
            if current_dict != data:
                self._logger.info(f"Updating config file {path} with new defaults")
                self.save(config, path)
                # Update mtime after save
                try:
                    self._config_mtime = path.stat().st_mtime
                except OSError:
                    self._config_mtime = None

        # Update cache
        self._config = config
        self._config_path = path
        if self._config_mtime is None:  # if not set above
            try:
                self._config_mtime = path.stat().st_mtime
            except OSError:
                pass

        return config

    def _is_cached_config_valid(self, path: Path) -> bool:
        """Check if cached config exists and file hasn't been modified."""
        if self._config is None or self._config_path != path:
            return False
        try:
            current_mtime = path.stat().st_mtime
            return current_mtime == self._config_mtime
        except OSError:
            return False
