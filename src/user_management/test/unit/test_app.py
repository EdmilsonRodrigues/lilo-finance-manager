import unittest

from app.app import get_app
from app.sessions import db


class TestBasicAppFunctions(unittest.TestCase):
    def setUp(self):
        self.app = get_app('sqlite:///:memory:', test=True)
        self.client = self.app.test_client()

    def tearDown(self):
        with self.app.app_context():
            db.session.remove()
            db.drop_all()

    def test_check_if_tables_are_created(self):
        with self.app.app_context():
            self.assertTrue(
                db.engine.dialect.has_table(db.engine.connect(), 'user')
            )
