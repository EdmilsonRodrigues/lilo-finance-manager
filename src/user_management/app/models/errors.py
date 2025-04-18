from http import HTTPStatus


class ApplicationError(Exception):
    """Base class for all application errors."""

    status_code = HTTPStatus.BAD_REQUEST


class UnprocessableContentException(ApplicationError):
    """Exception raised for unprocessable content."""

    status_code = HTTPStatus.UNPROCESSABLE_ENTITY


class NotFoundException(ApplicationError):
    """Exception raised for not found."""

    status_code = HTTPStatus.NOT_FOUND


class UnauthorizedException(ApplicationError):
    """Exception raised for unauthorized."""

    status_code = HTTPStatus.UNAUTHORIZED


class FobiddenException(ApplicationError):
    """Exception raised for forbidden."""

    status_code = HTTPStatus.FORBIDDEN


class TooManyRequestsException(ApplicationError):
    """Exception raised for too many requests."""

    status_code = HTTPStatus.TOO_MANY_REQUESTS
