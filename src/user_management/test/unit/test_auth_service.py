import unittest

from app.services.auth_service import AuthService


class TestAuthService(unittest.TestCase):
    def test_generate_token(self):
        token, expires_at = AuthService.generate_token(1)
        self.assertIsInstance(token, str)
        self.assertIsInstance(expires_at, str)

    def test_verify_token(self):
        token, _ = AuthService.generate_token(1)
        user_id = AuthService.verify_token(token)
        self.assertEqual(user_id, 1)
