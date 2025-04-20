import logging
from dataclasses import dataclass
from datetime import datetime
from functools import singledispatchmethod
from typing import Any, Self

from app.models.errors import ApplicationError, NotFoundException
from app.sessions import db

logger = logging.getLogger()


class BaseModel(db.Model):  # type: ignore
    """Base class for all SQLAlchemy models."""

    __abstract__ = True

    id = db.Column(db.Integer, primary_key=True)
    created_at = db.Column(db.DateTime, default=db.func.now())
    updated_at = db.Column(
        db.DateTime, default=db.func.now(), onupdate=db.func.now()
    )

    def create(self) -> Self:
        """
        Creates the model instance in the database.

        :return: The newly created model instance with an id.
        """
        logger.info(f'Creating {self.__class__.__name__}: {self}')
        db.session.add(self)
        db.session.commit()
        return self

    @classmethod
    def update(cls, id: int, fields: dict[str, Any]) -> Self:
        """
        Updates the model instance in the database.

        :param id: The id of the model instance to update.
        :type id: int
        :param fields: The fields to update.
        :type fields: dict[str, Any]
        :return: The updated model instance.
        """
        logger.info(
            f'Updating {cls.__name__} with id: {id} and fields: {fields}'
        )
        updated = cls.query.filter_by(id=id).update(fields)
        if not updated:
            logger.exception(f'{cls.__name__} with id: {id} not found')
            raise NotFoundException(f'{cls.__name__} not found')
        db.session.commit()
        return cls.get_one(id)

    @classmethod
    def delete(cls, id: int) -> None:
        """
        Deletes the model instance from the database.

        :param id: The id of the model instance to delete.
        :type id: int
        """
        logger.info(f'Deleting {cls.__name__} with id: {id}')
        deleted = cls.query.filter_by(id=id).delete()
        if not deleted:
            logger.exception(f'{cls.__name__} with id: {id} not found')
            raise NotFoundException(f'{cls.__name__} not found')
        db.session.commit()

    @classmethod
    def get_one(cls, id: int) -> Self:
        """
        Gets the model instance from the database.

        :param id: The id of the model instance to get.
        :type id: int
        :return: The model instance.
        """
        logger.info(f'Getting {cls.__name__} with id: {id}')
        obj = cls.query.get(id)
        if obj is None:
            logger.exception(f'{cls.__name__} with id: {id} not found')
            raise NotFoundException(f'{cls.__name__} not found')
        return obj

    @classmethod
    def get_many(cls, filters: dict[str, Any] = {}) -> list[Self]:
        """
        Gets a list of model instances from the database.

        :param filters: The filters to apply.
        :type filters: dict[str, Any]
        :return: The list of model instances.
        :rtype: list[Self]
        """
        logger.info(f'Getting {cls.__name__} with filters: {filters}')
        return cls.query.filter_by(**filters).all()

    def __repr__(self) -> str:
        """
        Returns a string representation of the model instance.

        :return: The string representation of the model instance.
        :rtype: str
        """
        fields = {
            field: value
            for field, value in vars(self).items()
            if (not field.startswith('_')) and not callable(value)
        }
        fields_str = ', '.join(
            f'{field}={value!r}' for field, value in fields.items()
        )
        return f'<{self.__class__.__name__} {fields_str}>'


@dataclass
class BaseClass:
    id: int
    created_at: str
    updated_at: str

    @classmethod
    def from_model(cls, model: BaseModel) -> Self:
        """
        Converts a SQLAlchemy model to a dataclass.

        :param model: The SQLAlchemy model to convert.
        :type model: BaseModel
        :return: The dataclass representation of the model.
        :rtype: Self
        """
        return cls(**{  # type: ignore
            field: value
            if not isinstance(value, datetime)
            else value.isoformat()
            for field, value in vars(model).items()
            if (not field.startswith('_'))
            and not callable(value)
            and field in cls.__dataclass_fields__.keys()
        })


type ErrorFormat = dict[str, dict[str, str]]
type ErrorResponseFormat = (
    tuple[ErrorFormat, int] | tuple[ErrorFormat, int, dict[str, str]]
)


@dataclass
class ErrorResponse:
    """Dataclass for the error response serialization."""

    @dataclass
    class ErrorDetail:
        """Dataclass for the error detail serialization."""

        status: int
        message: str

    details: ErrorDetail

    def jsonify(self) -> ErrorFormat:
        """
        Returns a JSON representation of the error response.

        :return: The JSON representation of the error response.
        :rtype: ErrorFormat
        """
        return {'details': vars(self.details)}

    @singledispatchmethod
    @classmethod
    def from_exception(cls, exception: Exception) -> ErrorResponseFormat:
        """
        Returns a tuple containing the JSON representation of the error
         and the status code.

        :param exception: The exception to convert to a JSON representation.
        :type exception: Exception
        :return: A tuple containing the JSON representation of the error
        and the status code. Can contain a third element if the exception
        has a headers attribute.
        :rtype: ErrorResponseFormat
        """
        return (
            cls(
                details=cls.ErrorDetail(
                    status=500,
                    message=str(exception),
                )
            ).jsonify()
        ), 500

    @from_exception.register
    @classmethod
    def _(cls, exception: ApplicationError) -> ErrorResponseFormat:
        exc: ErrorResponseFormat = (
            cls(
                details=cls.ErrorDetail(
                    status=exception.status_code,
                    message=str(exception.args[0]),
                )
            ).jsonify(),
            exception.status_code,
        )
        if exception.headers is not None:
            exc += (exception.headers,)
        return exc


@dataclass
class PaginatedResponse[T: BaseClass]:
    """Paginated response class"""

    page: int
    page_size: int
    total_items: int
    total_pages: int
    filters: dict[str, Any]
    items: list[T]


@dataclass
class JSONResponse[T: BaseClass]:
    """JSON response class"""

    status: str
    data: T | PaginatedResponse[T]

    def jsonify(self) -> dict[str, str | dict[str, Any]]:
        """
        Returns a JSON representation of the JSON response.

        :return: The JSON representation of the JSON response.
        :rtype: dict[str, dict[str, Any]]
        """
        return {'status': self.status, 'data': vars(self.data)}
