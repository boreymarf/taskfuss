from pydantic import BaseModel, Field
from typing import Literal

from src.logging.enums import LogFormatStyle


class ConsoleSettings(BaseModel):
    enabled: bool = True
    level: str = "INFO"
    style: LogFormatStyle = LogFormatStyle.RICH


class FileSettings(BaseModel):
    enabled: bool = True
    path: str = "logs/app.jsonl"
    level: str = "INFO"
    style: LogFormatStyle = LogFormatStyle.JSON
    max_bytes: int = 10 * 1024 * 1024
    backup_count: int = 7
    rotation: Literal["size", "daily"] = "size"


class ErrorFileSettings(BaseModel):
    enabled: bool = True
    path: str = "logs/errors.jsonl"
    level: str = "ERROR"
    style: LogFormatStyle = LogFormatStyle.JSON
    max_bytes: int = 10 * 1024 * 1024
    backup_count: int = 5


class LoggingConfig(BaseModel):
    root_level: str = "DEBUG"
    shushed_modules: list[str] = Field(
        default_factory=list
    )  
    console: ConsoleSettings = ConsoleSettings()
    file: FileSettings = FileSettings()
    error_file: ErrorFileSettings = ErrorFileSettings()
