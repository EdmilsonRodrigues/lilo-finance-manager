import unittest

from app.app import db, get_app
from app.models.user import User
from app.services.auth_service import AuthService


class TestUserRoutes(unittest.TestCase):
    def setUp(self):
        self.app = get_app('sqlite:///:memory:', test=True)
        self.client = self.app.test_client()
        self.user = User(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        self.user.password = User.hash_password(self.user.password).decode(
            'utf-8'
        )
        with self.app.app_context():
            self.user.create()
            self.token, _ = AuthService.generate_token(self.user.id)

    def tearDown(self):
        with self.app.app_context():
            db.drop_all()
            db.session.remove()
            db.engine.dispose()

    def test_update_user_email(self):
        response = self.client.patch(
            '/api/v1/users/me',
            json={
                'email': 'test_updated@gmail.com',
                'full_name': 'Test User Updated',
            },
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(
            response.json,
            {
                'status': 'success',
                'data': {
                    'id': self.user.id,
                    'email': 'test_updated@gmail.com',
                    'full_name': 'Test User',
                    'role': 'user',
                    'created_at': self.user.created_at.isoformat(),
                    'updated_at': self.user.updated_at.isoformat(),
                },
            },
        )

    def test_update_user_password(self):
        response = self.client.patch(
            '/api/v1/users/me',
            json={
                'old_password': 'password',
                'new_password': 'new_password',
            },
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(
            response.json,
            {
                'status': 'success',
                'data': {
                    'id': self.user.id,
                    'email': 'test@gmail.com',
                    'full_name': 'Test User',
                    'role': 'user',
                    'created_at': self.user.created_at.isoformat(),
                    'updated_at': response.json['data']['updated_at'],
                },
            },
        )

    def test_update_user_name(self):
        response = self.client.patch(
            '/api/v1/users/me',
            json={
                'full_name': 'Test User Updated',
            },
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(
            response.json,
            {
                'status': 'success',
                'data': {
                    'id': self.user.id,
                    'email': 'test@gmail.com',
                    'full_name': 'Test User Updated',
                    'role': 'user',
                    'created_at': self.user.created_at.isoformat(),
                    'updated_at': self.user.updated_at.isoformat(),
                },
            },
        )

    def test_fail_update_user_name_tried_update_role(self):
        response = self.client.patch(
            '/api/v1/users/me',
            json={
                'full_name': 'Test User Updated',
                'role': 'admin',
            },
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 422)
        self.assertEqual(
            response.json,
            {'details': {'status': 422, 'message': 'Unprocessable Content'}},
        )

    def test_delete_user(self):
        response = self.client.delete(
            '/api/v1/users/me',
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 204)

    def test_fail_delete_user(self):
        with self.app.app_context():
            self.user.delete(self.user.id)
        response = self.client.delete(
            '/api/v1/users/me',
            headers={'Authorization': f'Bearer {self.token}'},
        )
        self.assertEqual(response.status_code, 404)
        self.assertEqual(
            response.json,
            {'details': {'status': 404, 'message': 'User not found'}},
        )
