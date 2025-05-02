import unittest

from app.app import VERSION, db, get_app


class TestBasicAppFunctions(unittest.TestCase):
    def setUp(self):
        self.app = get_app('sqlite:///:memory:', test=True)
        self.client = self.app.test_client()

    def tearDown(self):
        with self.app.app_context():
            db.drop_all()
            db.session.remove()
            db.engine.dispose()

    def test_check_if_tables_are_created(self):
        with self.app.app_context():
            self.assertTrue(
                db.engine.dialect.has_table(db.engine.connect(), 'users')
            )
            self.assertFalse(
                db.engine.dialect.has_table(db.engine.connect(), 'base_model')
            )

    def test_health_check(self):
        response = self.client.get('/api/v1/')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json, {'message': VERSION})
