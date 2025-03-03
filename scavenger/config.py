import os

class Config():
    SECRET_KEY = os.environ.get('SECTRE_KEY') or 'ghjgecr123'
    DEBUG = os.environ.get('DEBUG') or False
