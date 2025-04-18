import unittest

from app.app import get_app
from app.sessions import db


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
