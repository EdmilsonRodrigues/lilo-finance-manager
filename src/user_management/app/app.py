from flask import Flask

from app.config import SECRET_KEY
from app.sessions import db


def get_app(db_uri: str, test: bool = False):
    app = Flask(__name__)
    app.secret_key = SECRET_KEY
    app.config['SQLALCHEMY_DATABASE_URI'] = db_uri
    app.config['TESTING'] = test
    db.init_app(app)

    from app.models import user  # noqa F401

    with app.app_context():
        db.create_all()

    return app


if __name__ == '__main__':
    from app.config import DATABASE_URI

    if DATABASE_URI is None:
        raise ValueError('DATABASE_URI is not set')

    app = get_app(DATABASE_URI)
    app.run(debug=True)
