import unittest

from app.app import get_app


class TestBasicAppFunctions(unittest.TestCase):
    def test_get_app(self):
        app = get_app()
        self.assertIsNotNone(app)
        self.assertIsNotNone(app.secret_key)
