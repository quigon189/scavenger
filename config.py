import os


class Config():
    SECRET_KEY = os.environ.get('SECTRE_KEY') or 'ghjgecr123'
    DUMP = os.environ.get('DUMP') or os.path.join(os.path.abspath(os.path.dirname(__file__)), "DUMP")
