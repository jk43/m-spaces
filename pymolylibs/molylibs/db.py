import os
from pymongo import MongoClient

def get_mongo_client() -> MongoClient:
    password = os.getenv('MONGO_PASSWORD')
    user = os.getenv('MONGO_USERNAME')
    host = os.getenv('MONGO_HOST')
    port = os.getenv('MONGO_PORT')
    db = os.getenv('MONGO_DATABASE')
    client = MongoClient(f'mongodb://{user}:{password}@{host}:{port}')
    return client[db]