from enum import Enum

class LogFormatStyle(str, Enum):
    SIMPLE = "simple"
    VERBOSE = "verbose"
    RICH = "rich"
    JSON = "json"
    DEBUG = "debug"
