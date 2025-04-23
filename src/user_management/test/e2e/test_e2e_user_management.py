import unittest
import warnings

from testcontainers.postgres import PostgresContainer

from app.app import db, get_app


class E2EUserManagement(unittest.TestCase):
    def setUp(self):
        warnings.filterwarnings('ignore', category=ResourceWarning)
        self.postgres = PostgresContainer()
        self.postgres.start()
        self.app = get_app(self.postgres.get_connection_url(), test=True)
        self.client = self.app.test_client()

    def tearDown(self):
        with self.app.app_context():
            db.drop_all()
            db.session.remove()
            db.engine.dispose()
        self.postgres.stop()

    def test_e2e_user_management(self):
        # Create a user
        user_data = {
            'email': 'test@gmail.com',
            'password': 'password123',
            'full_name': 'John Doe',
        }
        response = self.client.post('/api/v1/auth/signup', json=user_data)
        self.assertEqual(response.status_code, 201)

        # Login the user
        login_data = {'email': 'test@gmail.com', 'password': 'password123'}
        response = self.client.post('/api/v1/auth/login', json=login_data)
        self.assertEqual(response.status_code, 200)
        token = response.json['access_token']
        self.assertIsInstance(token, str)

        # Get the user's information
        response = self.client.get(
            '/api/v1/users/me', headers={'Authorization': f'Bearer {token}'}
        )
        self.assertEqual(response.status_code, 200)
        user_info = response.json
        self.assertEqual(user_info['data']['email'], 'test@gmail.com')
        self.assertEqual(user_info['data']['full_name'], 'John Doe')
        self.assertEqual(user_info['data']['role'], 'user')

        # Update the user's email
        update_data = {'email': 'newemail@gmail.com'}
        response = self.client.patch(
            '/api/v1/users/me',
            json=update_data,
            headers={'Authorization': f'Bearer {token}'},
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json['data']['email'], 'newemail@gmail.com')

        # Update the user's password
        update_data = {
            'new_password': 'newpassword',
            'old_password': 'password123',
        }
        response = self.client.patch(
            '/api/v1/users/me',
            json=update_data,
            headers={'Authorization': f'Bearer {token}'},
        )
        self.assertEqual(response.status_code, 200)

        # Update the user's name
        update_data = {'full_name': 'Jane Doe'}
        response = self.client.patch(
            '/api/v1/users/me',
            json=update_data,
            headers={'Authorization': f'Bearer {token}'},
        )
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json['data']['full_name'], 'Jane Doe')
        self.assertEqual(response.json['data']['email'], 'newemail@gmail.com')

        # Delete the user
        response = self.client.delete(
            '/api/v1/users/me', headers={'Authorization': f'Bearer {token}'}
        )
        self.assertEqual(response.status_code, 204)
