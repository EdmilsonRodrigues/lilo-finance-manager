import logging
from collections.abc import Callable
from functools import wraps

from app.models.base import ErrorResponse
from app.models.errors import (
    UnauthorizedException,
    UnprocessableContentException,
)
from app.models.user import User

logger = logging.getLogger()


def get_request(func):
    """
    Get the request from the request context.

    :param func: The function to decorate.
    :return: The decorated function.
    :rtype: Callable
    """
    from flask import request

    @wraps(func)
    def wrapper(*args, **kwargs):
        return func(*args, request=request, **kwargs)

    return wrapper


def catch_errors(func):
    """
    Catch all exceptions and return an error response.

    :param func: The function to decorate.
    :return: The decorated function.
    :rtype: Callable
    """

    @wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except Exception as exc:
            logger.exception(exc)
            return ErrorResponse.from_exception(exc)

    return wrapper


def authentication_required(func):
    """
    Check if the user is authenticated and add the user_id to the kwargs.

    :param func: The function to decorate.
    :return: The decorated function.
    :rtype: Callable
    """

    @wraps(func)
    @get_request
    def wrapper(*args, request, **kwargs):
        token = request.headers.get('Authorization')
        if not token:
            logger.exception(f'Missing token on request {request.get_json()}')
            raise UnauthorizedException('Missing token')
        logger.debug(f'Authenticating user with token: {token}')
        user_id = User.authenticate(token)
        logger.debug(f'User {user_id} authenticated')
        return func(*args, **kwargs, user_id=user_id)

    return wrapper


def parse_body(name: str, schema: type | Callable):
    """
    Parse the request body and validate it against the given schema.

    :param name: The name of the request body.
    :type name: str
    :param schema: The schema to validate the request body against.
    :type schema: type | Callable
    :return: The decorated function.
    :rtype: Callable
    """

    def decorator(func):
        @get_request
        @wraps(func)
        def wrapper(*args, request, **kwargs):
            try:
                body = request.get_json()
                logger.info(f'Parrsing request body: {body}')
                kwargs[name] = schema(**body)
            except TypeError as exc:
                logger.exception(f'Failed to parse request body {body}: {exc}')
                raise UnprocessableContentException(
                    'Unprocessable Content'
                ) from exc
            return func(*args, **kwargs)

        return wrapper

    return decorator
