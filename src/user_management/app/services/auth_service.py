import logging
from datetime import UTC, datetime, timedelta

import jwt

from app.config import JWT_EXPIRATION_TIME, SECRET_KEY
from app.models.errors import UnauthorizedException

logger = logging.getLogger()


class AuthService:
    """Service for authentication."""

    @staticmethod
    def generate_token(
        user_id: int, expiration_time: int = JWT_EXPIRATION_TIME
    ) -> tuple[str, str]:
        """
        Generates a JWT token for the given user ID.

        :param user_id: The ID of the user.
        :type user_id: int
        :return: The JWT token and the expiration time.
        :rtype: tuple[str, str]
        """
        expires_at = datetime.now(tz=UTC) + timedelta(seconds=expiration_time)
        payload = {'sub': str(user_id), 'exp': expires_at}
        token = jwt.encode(payload, SECRET_KEY, algorithm='HS256')
        return token, expires_at.strftime('%Y-%m-%d %H:%M:%S')

    @staticmethod
    def verify_token(token: str) -> int:
        """
        Verifies the given JWT token.
        Raises an exception if the token is invalid or expired.

        :param token: The JWT token.
        :type token: str
        :return: The user ID.
        :rtype: int
        """
        try:
            payload = jwt.decode(token, SECRET_KEY, algorithms=['HS256'])
            return int(payload['sub'])
        except (jwt.ExpiredSignatureError, jwt.InvalidTokenError) as exc:
            logger.exception(exc)
            raise UnauthorizedException('Invalid token') from exc
