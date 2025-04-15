import os
from pathlib import Path

from dotenv import load_dotenv

env_path = Path('env')
load_dotenv(dotenv_path=env_path)

SECRET_KEY = os.getenv('MELO_SECRET_KEY', os.urandom(32))
