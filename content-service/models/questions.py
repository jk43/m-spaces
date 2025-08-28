from os import error
from pymongo.results import InsertOneResult
#from langchain_core.pydantic_v1 import BaseModel as LCBaseModel, Field as LCField
from typing import List, Tuple, Dict, Any, Optional, Union, Set
from datetime import datetime
from bson import ObjectId 
from pydantic import BaseModel, Field

from molylibs.db import get_mongo_client
from molylibs.schemas import MongoObjectIdPydanticConfig

from models.counters import get_next_sequence, SEQUENCE_QUESTION

class Claim(BaseModel):
    user: str
    claim: str
    created_at: datetime
# For langchain
# Due to different pydantic version on langchain, we need to use different BaseModel class. LCBaseModel is used in langchain(pydandtic 1)
class QuestionOutput(BaseModel):
    question: str = Field(
        description="Question to be asked with LaTeX/Python(matplotlib)  formatted equation",
        default=None,
    )
    equation: str = Field(
        description="LaTeX/Python(matplotlib)  formatted equation from question",
        default=None,
    )
    answer_key: str = Field(
        description="choice number of answer, such as 1 or 2 or 3 or 4", default=None
    )
    answer_val: str = Field(
        description="LaTeX/Python(matplotlib)  formatted answer", default=None
    )
    diagram: str = Field(
        description="LaTeX/Python(matplotlib)  formatted diagram or graph if requires, If no diagram and graph, leave it as empty string ('')",
        default=None,
    )
    choices: List[str] = Field(
        description="If choice is equation, it must be LaTeX/Python(matplotlib)  formatted. No choice number",
        default=None,
    )
    desc: List[str] = Field(
        description="Very detail explanation of answer by steps. \ should be replaced with \\.",
        default=None,
    )
    def to_dict(self) -> Dict[str, str]:
        return {
            "desc": self.desc,
            "question": self.question,
            "equation": self.equation,
            "answer_key": self.answer_key,
            "answer_val": self.answer_val,
            "diagram": self.diagram,
            "choices": self.choices,
        }
# For langchain
# Save output to model.    
class Question(QuestionOutput):
    seq_id: Optional[int] = None
    llm: Optional[str] = None
    slugs: Optional[List[str]] = None
    diagram_image: Optional[str] = None
    created_by: Optional[str] = None
    created_at: Optional[datetime] = None
    error: Optional[str] = None
    error_prompt: Optional[str] = None
    error_result_from_ai: Optional[str] = None
    view_count: Optional[int] = None
    is_active: Optional[bool] = None
    rate_count: Optional[int] = None
    rate_total: Optional[float] = None
    claim_total: Optional[int] = None
    claims: Optional[List[Claim]] = None
    organization_id : Optional[str] = None

# Pydantic 2 model
class P2BaseQuestion(BaseModel):
    id: Optional[ObjectId] = Field(alias='_id')
    seq_id: Optional[int] = None
    llm: Optional[str] = None
    slugs: Optional[List[str]] = None
    diagram_image: Optional[str] = None
    created_by: Optional[str] = None
    created_at: Optional[datetime] = None
    error: Optional[str] = None
    error_prompt: Optional[str] = None
    error_result_from_ai: Optional[str] = None
    view_count: Optional[int] = None
    is_active: Optional[bool] = None
    rate_count: Optional[int] = None
    rate_total: Optional[float] = None
    claim_total: Optional[int] = None
    claims: Optional[List[Claim]] = None
    organization_id : Optional[str] = None

    question: str = None
    equation: str = None
    answer_key: str = None
    answer_val: str = None
    diagram: str = None
    choices: List[str] = []
    desc: List[str] = []

    class Config(MongoObjectIdPydanticConfig):
        pass

FETCHED_BY_ID = 'id'
FETCHED_BY_SLUG = 'slug'
FETCHED_BY_EXAMPLE = 'example'

# Model to be used in GET 
class ActivityQuestion(P2BaseQuestion):
    path: List[str] = [] #path from tree
    fetched_by: str = FETCHED_BY_ID

    def to_dict(self) -> Dict[str, str]:
        return {
            "id" : str(self.id),
            "seq_id": self.seq_id,
            "desc": self.desc,
            "question": self.question,
            "equation": self.equation,
            "answer_key": self.answer_key,
            "answer_val": self.answer_val,
            "diagram_image": self.diagram_image,
            "choices": self.choices,
            "slugs": self.slugs,
            "path": self.path,
            "fetched_by": self.fetched_by
        }
    
# Model to be used in GET /admin/questions
class QuestionListResponse(P2BaseQuestion):
    # id: Optional[ObjectId] = Field(alias='_id')
    # is_active: bool
    # created_at: datetime
    # llm: str
    # rate_count: int
    # rate_total: float
    # slugs: List[str]
    # view_count: int
    # error: Union[str, None]

    # # to decode ObjectId to string
    # class Config(MongoObjectIdPydanticConfig):
    #     pass

    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": str(self.id),
            "is_active": self.is_active,
            "created_at": self.created_at,
            "llm": self.llm,
            "rate_count": self.rate_count,
            "rate_total": self.rate_total,
            # Problem with Quasar 2.0 q-table, it only accept object with key-value pair
            #"slugs": self.slugs, 
            "slugs": {"subject": self.slugs[0], "chapter": self.slugs[1], "topic": self.slugs[2]},
            "view_count": self.view_count,
            "error": self.error,
        }

col = get_mongo_client()["questions"]

def insert_question(question:Question) -> InsertOneResult:
    question.seq_id = get_next_sequence(SEQUENCE_QUESTION)
    return col.insert_one(question.dict())

def get_questions(page: int = 0, per_page: int = 30, search: Dict[str, Any] | None = {}, fields: List[str] = [], order: int = 1, order_by: str = "created_at") -> Tuple[Any, int]:
    if search is None:
        search = {}
    try:
        total_count = col.count_documents(search)
        docs = col.find(search).skip((page - 1) * per_page).sort(order_by, order).limit(per_page)
        return (list(docs), total_count)
    except Exception as e:
        return str(e) 
    
def get_question_by_id(id: str) -> Optional[Dict[str, Any]]:
    question = col.find_one({"_id": ObjectId(id)})
    if question:
        question["_id"] = str(question["_id"])
    return question

def get_random_questions(size: int, slug: str = None) -> List[ActivityQuestion]:
    if slug:
        questions = col.aggregate([
            {"$match": {"slugs": slug}},
            {"$sample": {"size": size}}
        ])
    else:
        questions = col.aggregate([{"$sample": {"size": size}}])
    # for q in questions:
    #     a = ActivityQuestion(**q)
    #     print(a)
    return [ActivityQuestion(**q) for q in questions]