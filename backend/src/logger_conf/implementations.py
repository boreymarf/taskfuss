from datetime import datetime
import logging
import os

from rich.console import Console
from logging.handlers import RotatingFileHandler
from pythonjsonlogger.json import JsonFormatter

from src.logger_conf.formatters import RichFormatter, UvicornAccessFormatter
from src.logger_conf.handlers import DailySizeRotatingFileHandler, RichConsoleHandler

shared_console = Console(highlight=False)


def to_log_level(level: int | str) -> int:
    """Convert string or int to logging level constant."""
    if isinstance(level, int):
        return level
    try:
        return getattr(logging, level.upper())
    except AttributeError:
        raise ValueError(f"Invalid log level: {level}")


def clean_handlers() -> None:
    root_logger = logging.getLogger()
    root_logger.handlers.clear()


def set_root_logger(level: int | str = logging.DEBUG) -> None:
    level = to_log_level(level)
    root_logger = logging.getLogger()
    root_logger.setLevel(level)


def setup_rich_logging(level: int | str = logging.DEBUG) -> None:
    """Configures Rich handlers for root logger and uvicorn.access with custom formatters."""

    level = to_log_level(level)

    root_logger = logging.getLogger()

    # root
    rich_handler = RichConsoleHandler(shared_console)
    rich_formatter = RichFormatter()
    rich_handler.setFormatter(rich_formatter)
    rich_handler.setLevel(level)
    root_logger.addHandler(rich_handler)

    # uvicorn.access
    uvicorn_access_logger = logging.getLogger("uvicorn.access")
    uvicorn_access_logger.handlers.clear()
    uvicorn_access_logger.propagate = False

    uvicorn_handler = RichConsoleHandler(shared_console)
    uvicorn_formatter = UvicornAccessFormatter()
    uvicorn_handler.setFormatter(uvicorn_formatter)
    uvicorn_handler.setLevel(level)
    uvicorn_access_logger.addHandler(uvicorn_handler)


def setup_daily_json_logger(
    prefix: str = "app",
    log_dir: str = "logs",
    level: int | str = logging.INFO,
    max_bytes: int = 10 * 1024 * 1024,
    backup_days: int = 7,
):
    level = to_log_level(level)

    os.makedirs(log_dir, exist_ok=True)

    today = datetime.now().strftime("%Y-%m-%d")
    log_file = os.path.join(log_dir, f"{prefix}_{today}.jsonl")

    handler = DailySizeRotatingFileHandler(
        filename=log_file,
        max_bytes=max_bytes,
        backup_count=backup_days,
        when="midnight",
        interval=1,
        encoding="utf-8",
    )
    handler.setLevel(level)

    formatter = JsonFormatter(
        "%(asctime)s %(name)s:%(lineno)d %(levelname)s %(message)s",
        rename_fields={"asctime": "timestamp"},
        datefmt="%Y-%m-%dT%H:%M:%SZ",
    )
    handler.setFormatter(formatter)

    root_logger = logging.getLogger()
    root_logger.addHandler(handler)


def setup_json_logger(
    log_file: str = "logs/app.jsonl",
    level: int | str = logging.INFO,
    max_bytes: int = 10 * 1024 * 1024,
    backup_count: int = 5,
):
    level = to_log_level(level)

    os.makedirs(os.path.dirname(log_file), exist_ok=True)

    handler = RotatingFileHandler(
        filename=log_file,
        maxBytes=max_bytes,
        backupCount=backup_count,
        encoding="utf-8",
    )
    handler.setLevel(level)

    formatter = JsonFormatter(
        "%(asctime)s %(name)s %(levelname)s %(message)s",
        rename_fields={"asctime": "timestamp"},
        datefmt="%Y-%m-%dT%H:%M:%SZ",
    )
    handler.setFormatter(formatter)

    root_logger = logging.getLogger()
    root_logger.addHandler(handler)
