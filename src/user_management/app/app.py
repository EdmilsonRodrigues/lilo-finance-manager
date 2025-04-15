from flask import Flask

from app.config import SECRET_KEY


def get_app():
    app = Flask(__name__)
    app.secret_key = SECRET_KEY
    return app
