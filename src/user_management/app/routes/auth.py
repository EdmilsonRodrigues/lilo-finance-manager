from flask import Blueprint, jsonify, request

from app.models.errors import UnprocessableContentException
from app.models.user import CreateUser, User
from app.routes.dependencies import catch_errors

auth_bp = Blueprint('auth', __name__, url_prefix='/auth')


@auth_bp.route('/signup', methods=['POST'])
@catch_errors
def signup():
    try:
        user = CreateUser(**request.json)
    except TypeError as exc:
        raise UnprocessableContentException('Unprocessable Content') from exc
    user = User(**vars(user))
    user.create()
    return jsonify({'message': 'User created successfully'}), 201


@auth_bp.route('/login', methods=['POST'])
@catch_errors
def login():
    try:
        token, expiration_time = User.login(**request.json)
    except TypeError as exc:
        raise UnprocessableContentException('Unprocessable Content') from exc
    return {
        'access_token': token,
        'token_type': 'bearer',
        'expiration_time': expiration_time,
    }
