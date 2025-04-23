import logging
from dataclasses import dataclass, field
from typing import Any, Self

from email_validator import validate_email
from flask_bcrypt import check_password_hash, generate_password_hash

from app.config import Unset, UnsetType
from app.models.base import BaseClass, BaseModel
from app.models.errors import (
    UnauthorizedException,
    UnprocessableContentException,
)
from app.services.auth_service import AuthService
from app.sessions import db

logger = logging.getLogger()


class User(BaseModel):
    """SQLAlchemy model for the users table."""

    __tablename__ = 'users'

    email = db.Column(db.String(120), unique=True, nullable=False)
    password = db.Column(db.String(128), nullable=False)
    full_name = db.Column(db.String(80), nullable=False)
    role = db.Column(db.String(10), nullable=False)

    @staticmethod
    def hash_password(password: str) -> bytes:
        """
        Hashes the password using bcrypt.

        :param password: The password to hash.
        :type password: str
        :return: The hashed password.
        :rtype: bytes
        """
        return generate_password_hash(password)

    def check_password(self, password: str) -> bool:
        """
        Checks if the provided password matches the stored hash.

        :param password: The password to check.
        :type password: str
        :return: True if the password matches, False otherwise.
        :rtype: bool
        """
        return check_password_hash(self.password, password)

    @classmethod
    def update(cls, id: int, fields: dict) -> Self:
        if 'old_password' in fields:
            logger.debug(f'Updating password of user with id: {id}')
            old_password = fields.pop('old_password')
            user = cls.get_one(id)
            if not user.check_password(old_password):
                logger.warning('Invalid password')
                raise UnauthorizedException('Invalid password')
            return cls.update(id, {'password': fields['new_password']})
        return super().update(id, fields)

    @classmethod
    def login(cls, email: str, password: str) -> tuple[str, str]:
        """
        Logs in a user.

        :param email: The email of the user.
        :type email: str
        :param password: The password of the user.
        :type password: str
        :return: A tuple containing the JWT token and the expiration time.
        :rtype: tuple[str, str]
        """
        try:
            user = cls.query.filter_by(email=email).first()
            if user is None:
                raise UnauthorizedException('User not found')
            if user.check_password(password):
                return AuthService.generate_token(user.id)
            raise UnauthorizedException('Password is incorrect')
        except Exception as exc:
            logger.exception(exc)
            raise UnauthorizedException('Invalid credentials') from exc

    @classmethod
    def authenticate(cls, token: str) -> int:
        """
        Authenticates a user.

        :param token: The JWT token.
        :type token: str
        :return: The user id if the authentication is successful,
        raises an exception otherwise.
        :rtype: int
        """
        try:
            match token.split(' '):
                case ['Bearer', token]:
                    return AuthService.verify_token(token)
                case _:
                    raise UnauthorizedException('Invalid token')
        except Exception as exc:
            logger.exception(exc)
            raise UnauthorizedException(
                'Invalid Token', headers={'WWW-Authenticate': 'bearer'}
            ) from exc


@dataclass
class UserResponse(BaseClass):
    """Dataclass for the user response serialization."""

    email: str
    full_name: str
    role: str

    def __post_init__(self):
        try:
            self.email = validate_email(self.email).normalized
        except Exception as exc:
            logger.exception(exc)
            raise UnprocessableContentException(exc) from exc


@dataclass
class CreateUser:
    """Dataclass for the user creation request serialization."""

    email: str
    password: str
    full_name: str
    role: str = field(init=False, default='user')

    def __post_init__(self):
        self.__dict__['role'] = self.role
        try:
            self.email = validate_email(self.email).normalized
            self.password = User.hash_password(self.password).decode('utf-8')
        except Exception as exc:
            logger.exception(exc)
            raise UnprocessableContentException(exc) from exc


@dataclass
class PatchUser:
    """Dataclass for the user update request serialization."""

    full_name: str | UnsetType = Unset


@dataclass
class PatchUserPassword:
    """Dataclass for the user password update request serialization."""

    old_password: str
    new_password: str

    def __post_init__(self):
        if self.new_password == self.old_password:
            logger.warning('New password is the same as the old password')
            raise UnprocessableContentException(
                'New password must be different from the old password'
            )
        try:
            self.new_password = User.hash_password(self.new_password).decode(
                'utf-8'
            )
        except Exception as exc:
            logger.exception(exc)
            raise UnprocessableContentException(
                'Passwords are invalid'
            ) from exc


@dataclass
class PatchUserEmail:
    """Dataclass for the user email update request serialization."""

    email: str

    def __post_init__(self):
        try:
            self.email = validate_email(self.email).normalized
        except Exception as exc:
            logger.exception(exc)
            raise UnprocessableContentException(exc) from exc


def get_patch_fields(
    **data,
) -> dict[str, Any]:
    """
    Get the fields to patch from the request data.

    :param data: The request data.
    :type data: dict
    :return: The fields to patch.
    :rtype: dict[str, str]
    """
    match data:
        case {'email': email}:
            return vars(PatchUserEmail(email=email))

        case {'old_password': old_password, 'new_password': new_password}:
            return vars(
                PatchUserPassword(
                    old_password=old_password, new_password=new_password
                )
            )

        case _:
            try:
                return {
                    key: value
                    for key, value in vars(PatchUser(**data)).items()
                    if value is not Unset
                }
            except TypeError as exc:
                raise UnprocessableContentException(
                    'Unprocessable Content'
                ) from exc
