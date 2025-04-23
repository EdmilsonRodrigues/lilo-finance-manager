import os
from pathlib import Path

from dotenv import load_dotenv

env_path = Path('env')
load_dotenv(dotenv_path=env_path)

VERSION = '0.1.0'
JWT_EXPIRATION_TIME = 3600


SECRET_KEY = os.getenv('LFM_SECRET_KEY', os.urandom(32))

DATABASE_URI = os.getenv(
    'LFM_USER_MANAGEMENT_SQLALCHEMY_DATABASE_URI', 'sqlite:///:memory:'
)

type UnsetType = object
Unset = object()
