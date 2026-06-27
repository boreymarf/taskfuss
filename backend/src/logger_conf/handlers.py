import logging
from logging.handlers import TimedRotatingFileHandler
import os
import pprint
import sys
from typing import Any, TextIO, override

from rich.console import Console


class DefaultConsoleHandler(logging.StreamHandler[TextIO]):
    """Writes logs to stdout/stderr. Used for basic console logging."""

    def __init__(self, stream: TextIO | None = None) -> None:
        super().__init__(stream if stream is not None else sys.stdout)


class RichConsoleHandler(logging.Handler):
    """Outputs formatted/colored logs using Rich library. Used for pretty console output."""

    def __init__(self, console: Console, *args: Any, **kwargs: Any) -> None:
        super().__init__()
        self.console = console

    @override
    def emit(self, record: logging.LogRecord) -> None:
        msg = self.format(record)
        self.console.print(msg)


class DebugHandler(logging.Handler):
    """Prints full LogRecord objects for deep inspection. Used for debugging logging itself."""

    @override
    def emit(self, record: logging.LogRecord) -> None:
        pprint.pprint(record, stream=sys.stdout)


class DailySizeRotatingFileHandler(TimedRotatingFileHandler):
    def __init__(
        self,
        filename: str,
        max_bytes: int = 10 * 1024 * 1024,
        backup_count: int = 7,
        when: str = "midnight",
        interval: int = 1,
        encoding: str = "utf-8",
    ) -> None:
        super().__init__(
            filename,
            when=when,
            interval=interval,
            backupCount=backup_count,
            encoding=encoding,
        )
        self.max_bytes: int = max_bytes

    @override
    def shouldRollover(self, record: logging.LogRecord) -> bool:
        if super().shouldRollover(record):
            return True
        if self.max_bytes > 0 and os.path.exists(self.baseFilename):
            cur_size: int = os.path.getsize(self.baseFilename)
            if cur_size >= self.max_bytes:
                return True
        return False
