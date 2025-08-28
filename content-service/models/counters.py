from molylibs.db import get_mongo_client
from typing import Final

col = get_mongo_client()["counters"]

SEQUENCE_QUESTION: Final = "questions"

def get_next_sequence(name: str) -> int:
    ret = col.find_one_and_update(
        {"_id": name},
        {"$inc": {"seq": 1}},
        upsert=True,
        return_document=True,
    )
    return ret["seq"]