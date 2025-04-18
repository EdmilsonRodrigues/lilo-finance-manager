from dataclasses import dataclass
from functools import singledispatchmethod
from typing import Any, Self

from app.models.errors import ApplicationError, NotFoundException
from app.sessions import db


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
        updated = cls.query.filter_by(id=id).update(fields)
        if not updated:
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
        deleted = cls.query.filter_by(id=id).delete()
        if not deleted:
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
        obj = cls.query.get(id)
        if obj is None:
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
    def from_model(cls, model: BaseModel):
        return cls(**{
            field: value
            for field, value in vars(model).items()
            if (not field.startswith('_'))
            and not callable(value)
            and field in cls.__dataclass_fields__.keys()
        })


@dataclass
class ErrorResponse:
    @dataclass
    class ErrorDetail:
        status: int
        message: str

    details: ErrorDetail

    @singledispatchmethod
    @classmethod
    def from_exception(cls, exception: Exception):
        return vars(
            cls(
                details=vars(
                    cls.ErrorDetail(
                        status=500,
                        message=str(exception),
                    )
                )
            )
        ), 500

    @from_exception.register
    @classmethod
    def _(cls, exception: ApplicationError):
        return vars(
            cls(
                details=vars(
                    cls.ErrorDetail(
                        status=exception.status_code,
                        message=exception.args[0],
                    )
                )
            )
        ), exception.status_code


@dataclass
class PaginatedResponse[T: BaseClass]:
    page: int
    page_size: int
    total_items: int
    total_pages: int
    filters: dict[str, Any]
    items: list[T]


@dataclass
class JSONResponse[T: BaseClass]:
    status: str
    data: T | PaginatedResponse[T]
