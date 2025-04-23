import unittest

from app.app import db, get_app
from app.routes.auth import User


class TestAuthRoutes(unittest.TestCase):
    def setUp(self):
        self.app = get_app('sqlite:///:memory:', test=True)
        self.client = self.app.test_client()

    def tearDown(self):
        with self.app.app_context():
            db.drop_all()
            db.session.remove()
            db.engine.dispose()

    def test_signup(self):
        response = self.client.post(
            '/api/v1/auth/signup',
            json={
                'email': 'test@gmail.com',
                'password': 'password',
                'full_name': 'Test User',
            },
        )
        self.assertEqual(response.status_code, 201)
        self.assertEqual(
            response.json, {'message': 'User created successfully'}
        )

    def test_fail_signup(self):
        response = self.client.post(
            '/api/v1/auth/signup',
            json={
                'email': 'test@example.com',
                'password': 'password',
                'full_name': 'Test User',
                'role': 'admin',
            },
        )
        self.assertEqual(response.status_code, 422)
        self.assertEqual(
            response.json,
            {'details': {'status': 422, 'message': 'Unprocessable Content'}},
        )

    def _insert_user_in_db(self):
        user = User(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        user.password = user.hash_password(user.password).decode('utf-8')
        with self.app.app_context():
            user.create()

    def test_login(self):
        self._insert_user_in_db()
        response = self.client.post(
            '/api/v1/auth/login',
            json={
                'email': 'test@gmail.com',
                'password': 'password',
            },
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(
            response.json.keys(),
            {'access_token', 'token_type', 'expiration_time'},
        )
        self.assertEqual(response.json['token_type'], 'bearer')
        self.assertIsInstance(response.json['expiration_time'], str)
        self.assertIsInstance(response.json['access_token'], str)

    def test_fail_login(self):
        response = self.client.post(
            '/api/v1/auth/login',
            json={
                'email': 'test@gmail.com',
                'password': 'password',
            },
        )
        self.assertEqual(response.status_code, 401)
        self.assertEqual(
            response.json,
            {'details': {'status': 401, 'message': 'Invalid credentials'}},
        )

    def test_fail_login_with_invalid_password(self):
        self._insert_user_in_db()
        response = self.client.post(
            '/api/v1/auth/login',
            json={
                'email': 'test@gmail.com',
                'password': 'wrong_password',
            },
        )
        self.assertEqual(response.status_code, 401)
        self.assertEqual(
            response.json,
            {'details': {'status': 401, 'message': 'Invalid credentials'}},
        )

    def test_fail_login_wrong_payload(self):
        response = self.client.post(
            '/api/v1/auth/login',
            json={
                'email': 'test@gmail.com',
                'password': 'password',
                'full_name': 'Test User',
            },
        )
        self.assertEqual(response.status_code, 422)
        self.assertEqual(
            response.json,
            {'details': {'status': 422, 'message': 'Unprocessable Content'}},
        )
