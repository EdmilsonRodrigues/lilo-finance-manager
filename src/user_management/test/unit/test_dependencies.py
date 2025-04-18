import unittest

from app.models.errors import ApplicationError
from app.routes.dependencies import catch_errors


class TestDependencies(unittest.TestCase):
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
            raise ApplicationError('Test error')

        res = test()
        self.assertEqual(
            res, ({'details': {'status': 400, 'message': 'Test error'}}, 400)
        )

    def test_catch_errors_no_error(self):
        @catch_errors
        def test():
            return 'Test'

        res = test()
        self.assertEqual(res, 'Test')
