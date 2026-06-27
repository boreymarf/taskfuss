import logging
from typing import Any, Callable, override
from rich.markup import escape

Preprocessor = Callable[[logging.LogRecord], None]


class CompositeFormatter(logging.Formatter):
    """Wraps a base formatter with preprocessors that modify LogRecord before formatting."""

    def __init__(
        self,
        base_formatter: logging.Formatter,
        preprocessors: list[Preprocessor] | None = None,
    ) -> None:
        super().__init__()
        self.base_formatter = base_formatter
        self.preprocessors: list[Preprocessor] = preprocessors or []

    def add_preprocessor(self, preprocessor: Preprocessor) -> None:
        self.preprocessors.append(preprocessor)

    @override
    def format(self, record: logging.LogRecord) -> str:
        for preprocessor in self.preprocessors:
            preprocessor(record)
        return self.base_formatter.format(record)


class RichFormatter(logging.Formatter):
    """Formats logs with Rich colors: LEVEL module:lineno - message."""

    LEVEL_COLORS = {
        "DEBUG": "bright_black",
        "INFO": "green",
        "WARNING": "yellow",
        "ERROR": "bold red",
        "CRITICAL": "bold white on red",
    }

    def __init__(self, *args: Any, **kwargs: Any):
        super().__init__(*args, **kwargs)

    @override
    def format(self, record: logging.LogRecord) -> str:
        message = escape(record.getMessage())
        level_color = self.LEVEL_COLORS.get(record.levelname, "white")

        level_padded = f"{record.levelname:<8}"
        level_str = f"[{level_color}]{level_padded}[/]"
        location = f"{record.name}:{record.lineno}"

        return f"{level_str} [cyan]{location}[/] {message}"


class UvicornAccessFormatter(logging.Formatter):
    """Formats uvicorn.access logs: LEVEL client status method path proto. Used for HTTP access logs with status-based coloring."""

    # Status code colors with background
    STATUS_STYLES = {
        range(200, 300): "white on green",  # 2xx Success
        range(300, 400): "white on cyan",  # 3xx Redirect
        range(400, 500): "white on yellow",  # 4xx Client error
        range(500, 600): "bold white on red",  # 5xx Server error
    }

    def __init__(self, *args: Any, **kwargs: Any):
        super().__init__(*args, **kwargs)

    def _get_status_style(self, status_code: int) -> str:
        for status_range, style in self.STATUS_STYLES.items():
            if status_code in status_range:
                return style
        return "white on black"

    @override
    def format(self, record: logging.LogRecord) -> str:
        args = record.args
        if args and isinstance(args, (tuple, list)) and len(args) >= 5:
            client = str(args[0])
            method = str(args[1])
            path = str(args[2])
            protocol = str(args[3])
            status = int(str(args[4]))

            level_color = "green"
            level_padded = f"{record.levelname:<8}"
            level_str = f"[{level_color}]{level_padded}[/]"

            location = (
                f"{record.name}:{record.lineno}"
                if hasattr(record, "lineno")
                else record.name
            )

            client_str = f"[dim]{client}[/]"
            status_style = self._get_status_style(status)
            status_str = f"[{status_style}] {status} [/]"
            request_str = f'"{method} {path} HTTP/{protocol}"'

            return f"{level_str} [cyan]{location}[/] {client_str} {status_str} {request_str}"

        return super().format(record)
