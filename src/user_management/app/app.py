import atexit
import logging

from flask import Blueprint, Flask, jsonify

from app.config import SECRET_KEY, VERSION
from app.routes.auth import auth_bp
from app.routes.users import users_bp
from app.sessions import db

logger = logging.getLogger()
logger.setLevel(logging.INFO)


def on_exit():
    logger.info('User management service going down')


atexit.register(on_exit)


def get_app(db_uri: str, test: bool = False) -> Flask:
    """
    Function to create and configure the Flask application.

    :param db_uri: The URI of the database to connect to.
    :type db_uri: str
    :param test: Whether the application is running in test mode.
    :type test: bool
    :return: The configured Flask application.
    """
    logger.info('User management service going up')
    app = Flask(__name__)
    app.secret_key = SECRET_KEY
    app.config['SQLALCHEMY_DATABASE_URI'] = db_uri
    app.config['SQLALCHEMY_ENGINE_OPTIONS'] = {'pool_pre_ping': True}
    app.config['SQLALCHEMY_POOL_RECYCLE'] = 3600
    app.config['TESTING'] = test
    db.init_app(app)

    from app.models import user  # noqa F401

    with app.app_context():
        db.create_all()

    app.register_blueprint(api_bp)

    return app


api_bp = Blueprint('api', __name__, url_prefix='/api/v1')
api_bp.register_blueprint(auth_bp)
api_bp.register_blueprint(users_bp)


@api_bp.route('/', methods=['GET'])
def health_check():
    """Health check endpoint."""
    return jsonify({'message': VERSION}), 200


if __name__ == '__main__':  # pragma: no cover
    from app.config import DATABASE_URI

    if DATABASE_URI is None:
        raise ValueError('DATABASE_URI is not set')

    app = get_app(DATABASE_URI)
    app.run(debug=True)
