import redis
import os

(host, port) = os.getenv('REDIS_ADDR').split(':')

redis_client = redis.Redis(host=host, port=port, db=0)