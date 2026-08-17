import logging

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

def set_modules_log_level(modules: list[str] | str, level: int | str) -> None:
    """Set log level for specific modules."""
    module_names = [modules] if isinstance(modules, str) else modules
    for name in module_names:
        logging.getLogger(name).setLevel(level)
