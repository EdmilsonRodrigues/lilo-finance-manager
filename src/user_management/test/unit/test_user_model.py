import unittest

from app.models.user import (
    CreateUser,
    PatchUser,
    PatchUserEmail,
    PatchUserPassword,
    UnprocessableContentException,
    User,
    UserResponse,
    get_patch_fields,
)


class TestUserModel(unittest.TestCase):
    def test_user_instantiation(self):
        user = User(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        self.assertEqual(user.email, 'test@gmail.com')
        self.assertEqual(user.password, 'password')
        self.assertEqual(user.full_name, 'Test User')
        self.assertEqual(user.role, 'user')

    def test_user_password_hashing(self):
        user = User(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        hashed_password = user.hash_password('password').decode('utf-8')
        self.assertNotEqual(user.password, hashed_password)

    def test_compare_password(self):
        user = User(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        user.password = user.hash_password('password').decode('utf-8')
        self.assertTrue(user.check_password('password'))
        self.assertFalse(user.check_password('wrong_password'))


class TestUserResponse(unittest.TestCase):
    def test_user_response(self):
        user_response = UserResponse(
            id=54,
            created_at='2023-01-01 00:00:00',
            updated_at='2023-01-01 00:00:00',
            email='test@gmail.com',
            full_name='Test User',
            role='user',
        )
        self.assertEqual(user_response.id, 54)
        self.assertEqual(user_response.email, 'test@gmail.com')
        self.assertEqual(user_response.full_name, 'Test User')
        self.assertEqual(user_response.role, 'user')
        self.assertEqual(user_response.created_at, '2023-01-01 00:00:00')
        self.assertEqual(user_response.updated_at, '2023-01-01 00:00:00')

    def test_from_user_model(self):
        user = User(
            id=54,
            created_at='2023-01-01 00:00:00',
            updated_at='2023-01-01 00:00:00',
            email='test@gmail.com',
            password='password',
            full_name='Test User',
            role='user',
        )
        user_response = UserResponse.from_model(user)
        self.assertEqual(user_response.id, user.id)
        self.assertEqual(user_response.email, user.email)
        self.assertEqual(user_response.full_name, user.full_name)
        self.assertEqual(user_response.role, user.role)
        self.assertEqual(
            user_response.created_at,
            user.created_at,
        )
        self.assertEqual(
            user_response.updated_at,
            user.updated_at,
        )

    def test_bad_email(self):
        with self.assertRaises(UnprocessableContentException):
            UserResponse(
                id=54,
                created_at='2023-01-01 00:00:00',
                updated_at='2023-01-01 00:00:00',
                email='test@example',
                full_name='Test User',
                role='user',
            )


class TestUserCreateModel(unittest.TestCase):
    def test_user_create_model(self):
        user_create = CreateUser(
            email='test@gmail.com',
            password='password',
            full_name='Test User',
        )
        self.assertEqual(user_create.email, 'test@gmail.com')
        self.assertNotEqual(user_create.password, 'password')
        self.assertEqual(user_create.full_name, 'Test User')
        self.assertEqual(user_create.role, 'user')

    def test_user_create_model_try_to_pass_role(self):
        with self.assertRaises(TypeError):
            CreateUser(
                email='test@gmail.com',
                password='password',
                full_name='Test User',
                role='admin',
            )

    def test_user_create_model_bad_email(self):
        with self.assertRaises(UnprocessableContentException):
            CreateUser(
                email='test@example',
                password='password',
                full_name='Test User',
            )

    def test_user_create_model_bad_password(self):
        with self.assertRaises(UnprocessableContentException):
            CreateUser(
                email='test@gmail.com',
                password=4583,
                full_name='Test User',
            )


class TestUserPatchModels(unittest.TestCase):
    def test_user_patch_model(self):
        user_patch = PatchUser(
            full_name='Test User',
        )
        self.assertEqual(user_patch.full_name, 'Test User')

    def test_user_patch_model_email(self):
        user_patch = PatchUserEmail(
            email='test@gmail.com',
        )
        self.assertEqual(user_patch.email, 'test@gmail.com')

    def test_user_patch_model_email_bad_email(self):
        with self.assertRaises(UnprocessableContentException):
            PatchUserEmail(
                email='test@example',
            )

    def test_user_patch_model_password(self):
        user_patch = PatchUserPassword(
            old_password='password',
            new_password='newpassword',
        )
        self.assertEqual(user_patch.old_password, 'password')
        self.assertNotEqual(user_patch.new_password, 'newpassword')

    def test_user_patch_model_password_bad_password(self):
        with self.assertRaises(UnprocessableContentException):
            PatchUserPassword(
                old_password='password',
                new_password=4583,
            )

    def test_user_patch_model_password_same_password(self):
        with self.assertRaises(UnprocessableContentException):
            PatchUserPassword(
                old_password='password',
                new_password='password',
            )

    def test_gen_patch_fields(self):
        user_patch = {
            'full_name': 'Test User',
        }
        self.assertEqual(
            get_patch_fields(user_patch),
            {'full_name': 'Test User'},
        )
        user_patch = {
            'email': 'test@gmail.com',
        }
        self.assertEqual(
            get_patch_fields(user_patch),
            {'email': 'test@gmail.com'},
        )
        user_patch = {
            'old_password': 'password',
            'new_password': 'newpassword',
        }
        self.assertEqual(
            get_patch_fields(user_patch)['old_password'], 'password'
        )
        self.assertNotEqual(
            get_patch_fields(user_patch)['new_password'],
            'newpassword',
        )
        self.assertEqual(
            get_patch_fields(user_patch).keys(),
            {'old_password', 'new_password'},
        )
        user_patch = {
            'full_name': 'Test User',
            'role': 'admin',
        }
        self.assertRaises(
            UnprocessableContentException, get_patch_fields, user_patch
        )
        user_patch = {}
        self.assertEqual(get_patch_fields(user_patch), {})
        user_patch = 'not a dict'
        self.assertRaises(
            UnprocessableContentException, get_patch_fields, user_patch
        )
