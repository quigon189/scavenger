import os

class Config():
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'ghjgecr123'
    DUMP = os.environ.get('DUMP') or '/dump'
