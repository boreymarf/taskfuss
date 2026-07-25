from fastapi import status
from pydantic import BaseModel



class ErrorDetail(BaseModel):
    detail: str


BAD_REQUEST = {
    status.HTTP_400_BAD_REQUEST: {
        "description": "Invalid request",
        "model": ErrorDetail,
    }
}

UNAUTHORIZED = {
    status.HTTP_401_UNAUTHORIZED: {
        "description": "Authentication required",
        "model": ErrorDetail,
    }
}

FORBIDDEN = {
    status.HTTP_403_FORBIDDEN: {
        "description": "Permission denied",
        "model": ErrorDetail,
    }
}

NOT_FOUND = {
    status.HTTP_404_NOT_FOUND: {
        "description": "Resource not found",
        "model": ErrorDetail,
    }
}

CONFLICT = {
    status.HTTP_409_CONFLICT: {
        "description": "Resource already exists",
        "model": ErrorDetail,
    }
}

NOT_IMPLEMENTED = {
    status.HTTP_501_NOT_IMPLEMENTED: {
        "description": "This feature is not implemented ¯\\_(ツ)_/¯",
        "model": ErrorDetail,
    }
}

INTERNAL_ERROR = {
    status.HTTP_500_INTERNAL_SERVER_ERROR: {
        "description": "Internal error.",
        "model": ErrorDetail,
    }
}

# Combinations
PROTECTED = {**UNAUTHORIZED, **FORBIDDEN}
