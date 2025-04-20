from flask import Blueprint, jsonify

from app.models.base import JSONResponse
from app.models.user import User, UserResponse, get_patch_fields
from app.routes.dependencies import (
    authentication_required,
    catch_errors,
    parse_body,
)

users_bp = Blueprint('users', __name__, url_prefix='/users')


@users_bp.route('/me', methods=['GET'])
@catch_errors
@authentication_required
def get_me(user_id: int):
    """Get the current user."""
    user = User.get_one(user_id)
    return JSONResponse(
        status='success', data=UserResponse.from_model(user)
    ).jsonify()


@users_bp.route('/me', methods=['PATCH'])
@catch_errors
@authentication_required
@parse_body(name='update_fields', schema=get_patch_fields)
def update_me(user_id: int, update_fields: dict):
    """Update the current user."""
    user = User.update(user_id, update_fields)
    return JSONResponse(
        status='success', data=UserResponse.from_model(user)
    ).jsonify()


@users_bp.route('/me', methods=['DELETE'])
@catch_errors
@authentication_required
def delete_me(user_id: int):
    """Delete the current user."""
    User.delete(user_id)
    return jsonify(), 204
