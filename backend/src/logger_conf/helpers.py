import logging


def set_modules_log_level(modules: list[str] | str, level: int | str) -> None:
    """Set log level for specific modules."""
    module_names = [modules] if isinstance(modules, str) else modules
    for name in module_names:
        logging.getLogger(name).setLevel(level)
