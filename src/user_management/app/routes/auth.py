from collections import namedtuple

from flask import Blueprint, jsonify

from app.models.user import CreateUser, User
from app.routes.dependencies import catch_errors, parse_body

auth_bp = Blueprint('auth', __name__, url_prefix='/auth')
LoginData = namedtuple('LoginData', 'email password')


@auth_bp.route('/signup', methods=['POST'])
@catch_errors
@parse_body(name='user', schema=CreateUser)
def signup(user: CreateUser):
    """Signup a new user."""
    user = User(**vars(user))
    user.create()
    return jsonify({'message': 'User created successfully'}), 201


@auth_bp.route('/login', methods=['POST'])
@catch_errors
@parse_body(name='login_data', schema=LoginData)
def login(login_data: LoginData):
    """Login a user."""
    token, expiration_time = User.login(*login_data)
    return {
        'access_token': token,
        'token_type': 'bearer',
        'expiration_time': expiration_time,
    }
