from uuid import UUID


class AppException(Exception):
    """Base exception."""

    def __init__(self, detail: str, status_code: int = 400):
        super().__init__(detail)
        self.detail = detail
        self.status_code = status_code


# Generic


class NotFoundError(AppException):
    def __init__(self, resource: str, id: int | UUID | str):
        super().__init__(detail=f"{resource} with id {id} not found", status_code=404)


class ForbiddenError(AppException):
    def __init__(
        self, detail: str = "You don't have permission to access this resource"
    ):
        super().__init__(detail=detail, status_code=403)


class BadRequestError(AppException):
    def __init__(self, detail: str = "Bad request"):
        super().__init__(detail=detail, status_code=400)


class TodoError(AppException):
    def __init__(self, detail: str = "This feature is not implemented ¯\\_(ツ)_/¯"):
        super().__init__(detail=detail, status_code=501)


# User


class InvalidToken(AppException):
    def __init__(self):
        super().__init__(detail="Invalid token.", status_code=401)


class UserAlreadyExistsError(AppException):
    def __init__(self, login: str):
        super().__init__(detail=f"User '{login}' already exists", status_code=409)


class UserNotFound(AppException):
    def __init__(self, user_id: int):
        super().__init__(
            detail=f"User with id {user_id} does not exist.", status_code=422
        )


class InvalidCredentialsError(AppException):
    def __init__(self):
        super().__init__(detail="Invalid credentials.", status_code=401)
