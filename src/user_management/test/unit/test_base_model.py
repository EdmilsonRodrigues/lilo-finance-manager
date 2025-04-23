import inspect
import unittest
from contextlib import contextmanager
from unittest.mock import patch

from app.models.base import (
    ApplicationError,
    BaseClass,
    BaseModel,
    ErrorResponse,
    NotFoundException,
)


class TestBaseModel(unittest.TestCase):
    model = BaseModel(id=54, created_at='2023-01-01', updated_at='2023-01-01')

    @contextmanager
    def mock_query(self):
        self.mock_query_instance = unittest.mock.MagicMock()
        original = inspect.getattr_static(BaseModel, 'query')
        BaseModel.query = self.mock_query_instance
        yield
        BaseModel.query = original

    def test_instanciate(self):
        self.assertEqual(self.model.id, 54)
        self.assertEqual(self.model.created_at, '2023-01-01')
        self.assertEqual(self.model.updated_at, '2023-01-01')

    @patch('app.models.base.db.session.add')
    @patch('app.models.base.db.session.commit')
    def test_create(self, mock_commit, mock_add):
        mock_add.return_value = None
        mock_commit.return_value = None

        res = self.model.create()

        mock_add.assert_called_once_with(self.model)
        mock_commit.assert_called_once_with()
        self.assertIs(res, self.model)

    @patch('app.models.base.db.session.commit')
    @patch('app.models.base.BaseModel.get_one')
    def test_update(self, mock_get_one, mock_commit):
        mock_get_one.return_value = self.model
        with self.mock_query():
            mock_filter_by = self.mock_query_instance.filter_by
            mock_filter_by.return_value.update.return_value = 1

            res = self.model.update(id=1, fields={'name': 'John'})

            mock_commit.assert_called_once_with()
            self.mock_query_instance.filter_by.assert_called_once_with(id=1)
            self.mock_query_instance.filter_by.return_value.update.assert_called_once_with({
                'name': 'John'
            })
            mock_get_one.assert_called_once_with(1)
            self.assertIs(res, self.model)

    def test_update_no_user_found(self):
        with self.mock_query():
            mock_filter_by = self.mock_query_instance.filter_by
            mock_filter_by.return_value.update.return_value = 0

            self.assertRaises(
                NotFoundException,
                self.model.update,
                id=1,
                fields={'name': 'John'},
            )

            mock_filter_by.assert_called_once_with(id=1)
            mock_filter_by.return_value.update.assert_called_once_with({
                'name': 'John'
            })

    @patch('app.models.base.db.session.commit')
    def test_delete(self, mock_commit):
        with self.mock_query():
            mock_filter_by = self.mock_query_instance.filter_by
            mock_filter_by.return_value.delete.return_value = 1

            self.model.delete(id=1)

            mock_filter_by.assert_called_once_with(id=1)
            mock_filter_by.return_value.delete.assert_called_once_with()
            mock_commit.assert_called_once_with()

    def test_delete_no_user_found(self):
        with self.mock_query():
            mock_filter_by = self.mock_query_instance.filter_by
            mock_filter_by.return_value.delete.return_value = 0

            self.assertRaises(NotFoundException, self.model.delete, id=1)

            mock_filter_by.assert_called_once_with(id=1)
            mock_filter_by.return_value.delete.assert_called_once_with()

    def test_get_one(self):
        with self.mock_query():
            mock_get = self.mock_query_instance.get
            mock_get.return_value = self.model

            res = self.model.get_one(id=1)

            mock_get.assert_called_once_with(1)
            self.assertIs(res, self.model)

    def test_get_one_no_user_found(self):
        with self.mock_query():
            mock_get = self.mock_query_instance.get
            mock_get.return_value = None

            self.assertRaises(NotFoundException, self.model.get_one, id=1)
            mock_get.assert_called_once_with(1)

    def test_get_many(self):
        with self.mock_query():
            mock_filter_by = self.mock_query_instance.filter_by
            mock_filter_by.return_value.all.return_value = [self.model]
            mock_filter_by.return_value.count.return_value = 1
            res = self.model.get_many()
            mock_filter_by.assert_called_once_with()
            self.assertEqual(res, [self.model])

    def test_repr(self):
        res = repr(self.model)
        self.assertEqual(
            res,
            "<BaseModel id=54, created_at='2023-01-01',"
            " updated_at='2023-01-01'>",
        )


class TestBaseClass(unittest.TestCase):
    def test_initialize(self):
        base_class = BaseClass(
            id=1, created_at='2023-01-01', updated_at='2023-01-01'
        )
        self.assertEqual(base_class.id, 1)
        self.assertEqual(base_class.created_at, '2023-01-01')
        self.assertEqual(base_class.updated_at, '2023-01-01')

    def test_from_model(self):
        base_model = BaseModel(
            id=1, created_at='2023-01-01', updated_at='2023-01-01'
        )
        base_class = BaseClass.from_model(base_model)
        self.assertEqual(base_class.id, 1)
        self.assertEqual(base_class.created_at, '2023-01-01')
        self.assertEqual(base_class.updated_at, '2023-01-01')


class TestErrorResponse(unittest.TestCase):
    def test_initialize(self):
        error_response = ErrorResponse(
            details=ErrorResponse.ErrorDetail(**{
                'status': 500,
                'message': 'Internal Server Error',
            })
        )
        self.assertEqual(
            error_response.details,
            ErrorResponse.ErrorDetail(
                status=500,
                message='Internal Server Error',
            ),
        )

    def test_from_known_exception(self):
        exceptions = [ApplicationError, *ApplicationError.__subclasses__()]
        for exception in exceptions:
            error_response = ErrorResponse.from_exception(
                exception('Test message')
            )
            self.assertEqual(
                error_response,
                (
                    {
                        'details': {
                            'status': exception.status_code,
                            'message': 'Test message',
                        }
                    },
                    exception.status_code,
                ),
            )

    def test_from_known_exception_with_headers(self):
        exceptions = [ApplicationError, *ApplicationError.__subclasses__()]
        for exception in exceptions:
            error_response = ErrorResponse.from_exception(
                exception('Test message', headers={'X-Test': 'Test'})
            )
            self.assertEqual(
                error_response,
                (
                    {
                        'details': {
                            'status': exception.status_code,
                            'message': 'Test message',
                        },
                    },
                    exception.status_code,
                    {'X-Test': 'Test'},
                ),
            )

    def test_unexpeted_exception(self):
        error_response = ErrorResponse.from_exception(
            Exception('Test message')
        )
        self.assertEqual(
            error_response,
            (
                {
                    'details': {
                        'status': 500,
                        'message': 'Test message',
                    }
                },
                500,
            ),
        )
