import json
from typing import Optional

from molylibs.redis import redis_client
from molylibs.network import JWTClaims
from molylibs.grpc import get_tree_ancestors
from models.questions import get_random_questions, QuestionOutput, FETCHED_BY_ID, FETCHED_BY_SLUG, FETCHED_BY_EXAMPLE

redis_key = 'math-service:questions'
example_size = 5
example_ttl = 300

class Activity:
    def __init__(self, claims: JWTClaims, id: str, slug: str):
        self.__claims = claims
        self.__id = None if id == 'undefined' else id
        self.__slug = None if slug == 'undefined' else slug
        self.__redis_key = f'{redis_key}:{slug}'
    
    def fetch(self):
        if self.__id == None:
            # Guest
            if not self.__claims.email:
                return self.__sample()
            # Member
            return self.__random()
        return self.__fetch()

    def __fetch(self):
        print('fetch')
        pass

    def __sample(self):
        #try to fetch from redis
        questions = redis_client.get(self.__redis_key)
        if questions:
            return json.loads(questions)
        questions = get_random_questions(example_size, self.__slug)
        output = []
        for q in get_random_questions(example_size, self.__slug):
            trees = get_tree_ancestors(org_id=q.organization_id, role='guest', slug=q.slugs[-1])
            q.path = [{'slug':t.attributes.slug, 'label':t.attributes.label} for t in trees]
            q.fetched_by = FETCHED_BY_EXAMPLE
            output.append(q.to_dict())
        #questions = [q.to_dict() for q in questions]
        redis_client.set(self.__redis_key, json.dumps(output), ex=example_ttl)
        return output

    def __random(self):
        print('random')
        pass