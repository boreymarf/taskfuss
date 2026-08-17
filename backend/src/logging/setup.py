import logging
import os
from logging.handlers import RotatingFileHandler
from pythonjsonlogger.json import JsonFormatter
from rich.console import Console

from src.logging.enums import LogFormatStyle
from src.logging.settings import LoggingConfig
from src.logging.formatters import RichFormatter, UvicornAccessFormatter
from src.logging.handlers import RichConsoleHandler, DefaultConsoleHandler, DailySizeRotatingFileHandler
from src.logging.utils import to_log_level, clean_handlers, set_root_logger, set_modules_log_level


def setup_logging(config: LoggingConfig) -> None:
    """
    Fully configure logging according to the provided settings.
    Clears all existing handlers and sets up new ones.
    """
    # 1. Reset root logger
    clean_handlers()
    set_root_logger(config.root_level)

    # 2. Console logging
    if config.console.enabled:
        if config.console.style == LogFormatStyle.RICH:
            # Rich console output
            console = Console(highlight=False)
            handler = RichConsoleHandler(console)
            formatter = RichFormatter()
            handler.setFormatter(formatter)
            handler.setLevel(to_log_level(config.console.level))
            logging.getLogger().addHandler(handler)

            # Special handling for uvicorn.access
            uvicorn_logger = logging.getLogger("uvicorn.access")
            uvicorn_logger.handlers.clear()
            uvicorn_logger.propagate = False
            uvicorn_handler = RichConsoleHandler(console)
            uvicorn_handler.setFormatter(UvicornAccessFormatter())
            uvicorn_handler.setLevel(to_log_level(config.console.level))
            uvicorn_logger.addHandler(uvicorn_handler)
        else:
            # Plain console output (non‑Rich)
            handler = DefaultConsoleHandler()

            # Default format (used for unknown styles, e.g. if someone adds a new enum value)
            fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"

            if config.console.style == LogFormatStyle.SIMPLE:
                fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
            elif config.console.style == LogFormatStyle.VERBOSE:
                fmt = "%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s"
            elif config.console.style == LogFormatStyle.DEBUG:
                fmt = "%(asctime)s - %(name)s - %(levelname)s - %(module)s:%(funcName)s:%(lineno)d - %(message)s"
            elif config.console.style == LogFormatStyle.JSON:
                formatter = JsonFormatter(
                    "%(asctime)s %(name)s %(levelname)s %(message)s",
                    rename_fields={"asctime": "timestamp"},
                    datefmt="%Y-%m-%dT%H:%M:%SZ",
                )
                handler.setFormatter(formatter)
            # else: keep default fmt (e.g. if style is not one of the above, fallback to SIMPLE)

            # For non‑JSON styles, use a standard logging.Formatter
            if config.console.style != LogFormatStyle.JSON:
                formatter = logging.Formatter(fmt, datefmt="%Y-%m-%d %H:%M:%S")
                handler.setFormatter(formatter)

            handler.setLevel(to_log_level(config.console.level))
            logging.getLogger().addHandler(handler)

    # 3. Main file logging
    if config.file.enabled:
        file_cfg = config.file
        os.makedirs(os.path.dirname(file_cfg.path) or ".", exist_ok=True)

        if file_cfg.rotation == "daily":
            handler = DailySizeRotatingFileHandler(
                filename=file_cfg.path,
                max_bytes=file_cfg.max_bytes,
                backup_count=file_cfg.backup_count,
                when="midnight",
                interval=1,
                encoding="utf-8",
            )
        else:  # size rotation
            handler = RotatingFileHandler(
                filename=file_cfg.path,
                maxBytes=file_cfg.max_bytes,
                backupCount=file_cfg.backup_count,
                encoding="utf-8",
            )

        if file_cfg.style == LogFormatStyle.JSON:
            formatter = JsonFormatter(
                "%(asctime)s %(name)s:%(lineno)d %(levelname)s %(message)s",
                rename_fields={"asctime": "timestamp"},
                datefmt="%Y-%m-%dT%H:%M:%SZ",
            )
        else:
            fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
            formatter = logging.Formatter(fmt, datefmt="%Y-%m-%d %H:%M:%S")

        handler.setFormatter(formatter)
        handler.setLevel(to_log_level(file_cfg.level))
        logging.getLogger().addHandler(handler)

    # 4. Error file logging
    if config.error_file.enabled:
        err_cfg = config.error_file
        os.makedirs(os.path.dirname(err_cfg.path) or ".", exist_ok=True)

        handler = RotatingFileHandler(
            filename=err_cfg.path,
            maxBytes=err_cfg.max_bytes,
            backupCount=err_cfg.backup_count,
            encoding="utf-8",
        )

        if err_cfg.style == LogFormatStyle.JSON:
            formatter = JsonFormatter(
                "%(asctime)s %(name)s:%(lineno)d %(levelname)s %(message)s",
                rename_fields={"asctime": "timestamp"},
                datefmt="%Y-%m-%dT%H:%M:%SZ",
            )
        else:
            fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
            formatter = logging.Formatter(fmt, datefmt="%Y-%m-%d %H:%M:%S")

        handler.setFormatter(formatter)
        handler.setLevel(to_log_level(err_cfg.level))
        logging.getLogger().addHandler(handler)

    # 5. Shush specific modules (set to WARNING or higher)
    if config.shushed_modules:
        set_modules_log_level(config.shushed_modules, logging.WARNING)
