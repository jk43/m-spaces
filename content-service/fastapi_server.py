from hashlib import pbkdf2_hmac
from urllib import response
from fastapi import FastAPI, Request, HTTPException, WebSocket, WebSocketDisconnect, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from langchain.chains import llm
from matplotlib.pyplot import cla
from numpy import number
from pydantic import BaseModel, Field
from typing import List, Dict, Generic, Optional, TypeVar, Any
from grpc import insecure_channel
from datetime import datetime

from molylibs.schemas import HttpSuccessResponse, HttpSuccessResponseWithTotalCount, HttpErrorResponse
from molylibs.network import JWTMiddlewear, get_jwt_from_request, get_claims_from_websocket, WebSocketManager
from models import questions
from service.math import create_questions, MongoStorage, PickleStorage, StdoutStorage
from models.questions import get_questions, get_question_by_id, QuestionListResponse
from models.counters import get_next_sequence, SEQUENCE_QUESTION

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  
    allow_credentials=True,
    allow_methods=["*"],  
    allow_headers=["*"], 
)
app.add_middleware(JWTMiddlewear)

# base payload model
T = TypeVar('T')

class Payload(BaseModel, Generic[T]):
    data: T

class PostMathPayload(BaseModel):
    tree: List[str]
    llm: str
    numberOfQuestions: int

@app.post("/admin/math")
async def post_math(json: Payload[PostMathPayload], request: Request, background_tasks: BackgroundTasks) -> HttpSuccessResponse:
    claims = get_jwt_from_request(request)
    #math = Math(llm=json.data.llm, storage=[StdoutStorage, MongoStorage], claims=claims)
    background_tasks.add_task(create_questions, llm=json.data.llm, storage=[StdoutStorage, MongoStorage], claims=claims, slug=json.data.tree[-1], number_of_questions=json.data.numberOfQuestions)
    #await create_questions(llm=json.data.llm, storage=[StdoutStorage, MongoStorage], claims=claims, slug=json.data.tree[-1], number_of_questions=json.data.numberOfQuestions)
    return HttpSuccessResponse(data=True)

def success_response(data):
    return {
        "result": "success",
        "data": data
    }

from molylibs.grpc import broadcast_message, send_progress
from service.math import Dummy
import molylibs.pb.message_service_pb2 as pb

@app.get("/admin/questions")
async def get_admin_questions(request: Request, page: int = 1, rowsPerPage: int = 15, sortBy:str= "createdAt", descending:bool=True) -> HttpSuccessResponseWithTotalCount:
    get_next_sequence(SEQUENCE_QUESTION)
    # print("pb.ProgressCommand.ADD,", pb.ProgressCommand.ADD, pb.MessageType.TEXT)
    # send_progress(channel="admin", context=b"Init!!!", progress=30, id="2", progress_ctx="post_question", command=pb.ProgressCommand.UPDATE, to="admin")
    # d = Dummy()
    # d.run()
    questions, total = get_questions(page=page, per_page=rowsPerPage, order= 1 if descending else -1, order_by=sortBy)
    output = []
    for q in questions:
        res = QuestionListResponse(**q)
        model = res.to_dict()
        output.append(model)
    return HttpSuccessResponseWithTotalCount(data=output, total=total)

@app.get("/admin/question/{id}")
async def get_admin_question(id: str):
    question = get_question_by_id(id)
    return HttpSuccessResponse(data=question)

from molylibs.redis import redis_client
from models.questions import get_random_questions
from service.activity import Activity
@app.get("/question/{slug}/{id}")
async def get_question(slug: str, id: str, request: Request):
    # claims = get_jwt_from_request(request)
    # print(claims.email)
    question = Activity(claims=get_jwt_from_request(request)[0], id=id, slug=slug) 
    #redis_client.set("test", "test")
    # test = redis_client.get("test")
    # print(test)
    #question = get_question_by_id(id)
    #claims = get_jwt_from_request(request)
    # questions = get_random_questions(3)
    # redis_client.set("math-service:questions", json.dumps([q.to_dict() for q in questions]), ex=100)

    # print(json.dumps([q.to_dict() for q in questions]))
    return HttpSuccessResponse(data=question.fetch())

@app.post("problem")
async def post_problem():
    pass

def start_fastapi():
    import uvicorn
    uvicorn.run("fastapi_server:app", host="0.0.0.0", port=80, reload=True, workers=1)