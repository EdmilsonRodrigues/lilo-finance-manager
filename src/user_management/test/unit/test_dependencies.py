import unittest
from collections import namedtuple
from contextlib import contextmanager
from functools import wraps
from unittest.mock import MagicMock, patch

from app.routes.dependencies import (
    UnauthorizedException,
    UnprocessableContentException,
    authentication_required,
    catch_errors,
    get_request,
    parse_body,
)


class TestDependencies(unittest.TestCase):
    @contextmanager
    def _mock_get_request(self):
        self.mock_request = MagicMock()

        def mock(func):
            @wraps(func)
            def wrapper(*args, **kwargs):
                return func(*args, request=self.mock_request, **kwargs)

            return wrapper

        from app.routes import dependencies

        dependencies.get_request = mock
        try:
            yield
        finally:
            dependencies.get_request = get_request

    def test_catch_errors(self):
        @catch_errors
        def test():
            raise Exception('Test error')

        res = test()
        self.assertEqual(
            res, ({'details': {'status': 500, 'message': 'Test error'}}, 500)
        )

    def test_catch_errors_with_application_error(self):
        @catch_errors
        def test():
            raise UnauthorizedException('Test error')

        res = test()
        self.assertEqual(
            res, ({'details': {'status': 401, 'message': 'Test error'}}, 401)
        )

    def test_catch_errors_no_error(self):
        @catch_errors
        def test():
            return 'Test'

        res = test()
        self.assertEqual(res, 'Test')

    @patch('app.routes.dependencies.User.authenticate', return_value=1)
    def test_authentication_required(self, mock_):
        with self._mock_get_request():

            @authentication_required
            def test(user_id: int):
                return user_id

            self.mock_request.headers = {'Authorization': 'Bearer token'}
            res = test()

        self.assertEqual(res, 1)

    @patch('app.routes.dependencies.User.authenticate', return_value=None)
    def test_authentication_required_no_token(self, mock_authenticate):
        with self._mock_get_request():

            @authentication_required
            def test(user_id: int):
                return user_id

            self.mock_request.headers = {}
            with self.assertRaises(UnauthorizedException):
                test()

    def test_parse_body(self):
        with self._mock_get_request():

            @parse_body('body', dict)
            def test(body):
                return body

            self.mock_request.get_json.return_value = {'name': 'John Doe'}
            res = test()

        self.assertEqual(res, {'name': 'John Doe'})

    def test_parse_body_invalid_schema(self):
        schema = namedtuple('Schema', ['parse_obj'])

        with self._mock_get_request():

            @parse_body('body', schema)
            def test(body):
                return body

            self.mock_request.get_json.return_value = {'name': 'John Doe'}
            with self.assertRaises(UnprocessableContentException):
                test()
