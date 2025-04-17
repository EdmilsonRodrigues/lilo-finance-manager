from dataclasses import dataclass

import bcrypt
from email_validator import validate_email

from app.config import Unset, UnsetType
from app.models.base import BaseClass, BaseModel
from app.models.errors import UnprocessableContentException
from app.sessions import db


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
        :return: The hashed password.
        :rtype: bytes
        """
        salt = bcrypt.gensalt()
        hashed_password = bcrypt.hashpw(password.encode('utf-8'), salt)
        return hashed_password

    def check_password(self, password: str) -> bool:
        """
        Checks if the provided password matches the stored hash.

        :param password: The password to check.
        :return: True if the password matches, False otherwise.
        :rtype: bool
        """
        return bcrypt.checkpw(
            password.encode('utf-8'), self.password.encode('utf-8')
        )


@dataclass
class UserResponse(BaseClass):
    email: str
    full_name: str
    role: str

    def __post_init__(self):
        try:
            self.email = validate_email(self.email).normalized
        except Exception as exc:
            raise UnprocessableContentException(exc) from exc


@dataclass
class CreateUser:
    email: str
    password: str
    full_name: str
    role: str = 'user'

    def __post_init__(self):
        try:
            self.email = validate_email(self.email).normalized
            self.password = User.hash_password(self.password)
        except Exception as exc:
            raise UnprocessableContentException(exc) from exc


@dataclass
class PatchUser:
    full_name: str | UnsetType = Unset


@dataclass
class PatchUserPassword:
    old_password: str
    new_password: str

    def __post_init__(self):
        if self.new_password == self.old_password:
            raise UnprocessableContentException(
                'New password must be different from the old password'
            )
        try:
            self.new_password = User.hash_password(self.new_password)
        except Exception as exc:
            raise UnprocessableContentException(
                'Passwords are invalid'
            ) from exc


@dataclass
class PatchUserEmail:
    email: str

    def __post_init__(self):
        try:
            self.email = validate_email(self.email).normalized
        except Exception as exc:
            raise UnprocessableContentException(exc) from exc


def get_patch_fields(
    data: dict,
) -> dict[str, str]:
    match data:
        case {'email': email}:
            return vars(PatchUserEmail(email=email))
        case {'old_password': old_password, 'new_password': new_password}:
            return vars(
                PatchUserPassword(
                    old_password=old_password, new_password=new_password
                )
            )
        case {}:
            try:
                return {
                    key: value
                    for key, value in vars(PatchUser(**data)).items()
                    if value is not Unset
                }
            except TypeError as exc:
                raise UnprocessableContentException('Invalid data') from exc
    raise UnprocessableContentException('Invalid data')
