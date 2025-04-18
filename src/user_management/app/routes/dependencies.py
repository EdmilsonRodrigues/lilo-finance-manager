from functools import wraps

from app.models.base import ErrorResponse


def catch_errors(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except Exception as exc:
            return ErrorResponse.from_exception(exc)

    return wrapper
