import asyncio
from datetime import datetime
import pickle
from re import sub
from time import sleep
from typing import List, Tuple, Dict, Any, Optional
from urllib import response
from click import prompt
from httpx import request
from langchain.chains import llm
from pydantic import BaseModel, Field
from abc import abstractmethod
from hashlib import md5
import os
from langchain.smith.evaluation import progress
from numpy.matlib import rand
import requests
import random

from molylibs.network import JWTClaims, WebSocketManager, websocket_connections, ProgressAbstract
from molylibs.schemas import Tree
from molylibs.grpc import get_tree_ancestors, get_tree, is_end, get_tree_as_array, broadcast_message
from molylibs.db import get_mongo_client
from starlette import status

from models.questions import insert_question, QuestionOutput, Question
from grpc_server import CTX_MATH_DIAGRAM

from langchain.schema import output_parser
from langchain_openai import ChatOpenAI
from langchain_anthropic import ChatAnthropic
from langchain.prompts import PromptTemplate
from langchain.output_parsers import ResponseSchema, StructuredOutputParser
#from langchain_core.pydantic_v1 import BaseModel, Field
from langchain.output_parsers import PydanticOutputParser
from langchain.callbacks.base import BaseCallbackHandler
from langchain.schema import LLMResult

class QuestionStorage:
    @abstractmethod
    def save(self, data: Question) -> Any:
        pass

    @abstractmethod
    def get(self, id: Any) -> Question:
        pass


class MongoStorage(QuestionStorage):
    def save(self, data: Question) -> Any:
        print("MongoMathStorage called")
        return insert_question(data)

    def get(self, id: Any) -> Question:
        pass


class StdoutStorage(QuestionStorage):
    def save(self, data: Question) -> Any:
        print("StdoutMathStorage called", data)
        return None

    def get(self, id: Any) -> Question:
        pass

class PickleStorage(QuestionStorage):
    res: str

    def __init__(self, path: str):
        self.path = path

    def save(self, data: Question) -> Any:
        with open(self.path + "/" + "-".join(["_".join(data.slugs)]), "wb") as f:
            pickle.dump(data, f)
        return None

    def get(self, id: Any) -> Question:
        pass

# Test class 
class Dummy(ProgressAbstract):

    def __init__(self):
        super().__init__(channel='post-question-progress', progress_ctx='post_question', to=None)
    async def run(self):
        sleep(1)
        print("Processing 10%")
        self._update_progress("Processing", 10)
        sleep(1)
        print("Processing 20%")
        self._update_progress("Processing", 20)
        sleep(1)
        self._done()
        pass

async def create_questions(llm: str, storage: List[QuestionStorage], claims: JWTClaims, slug: str, number_of_questions: int):
    # for w in range(number_of_questions):
    #     d = Dummy()
    #     tasks = asyncio.ensure_future(d.run())
    #     d.run()
    # await asyncio.gather(*tasks)
    # return

    jwt_claims = claims[0]
    workloads = []
    # create workload based on the number of questions and the given slug
    # if slug is not an end node, get the tree and select random nodes
    if not is_end(org_id=jwt_claims.org_id, role=jwt_claims.role, slug=slug) :
        trees = get_tree_as_array(org_id=jwt_claims.org_id, role=jwt_claims.role, slug=slug)
        if len(trees) < number_of_questions:
            number_of_questions = len(trees)
        r_trees = random.sample(trees, number_of_questions)
        for tree in r_trees :
            workloads.append([Tree(**tree[1]), Tree(**tree[2]), Tree(**tree[3])])
    else :
        [_, subject, chapter, topic] = get_tree_ancestors(
        org_id=jwt_claims.org_id, role=jwt_claims.role, slug=slug
        )
        for i in range(number_of_questions):
            workloads.append([subject, chapter, topic])

    for w in workloads:
        math = Math(llm=llm, storage=storage, claims=claims)
        #await math.create_question(w)
        tasks = asyncio.ensure_future(math.create_question(w))
    await asyncio.gather(*tasks)
    return

# The dictionary to store the progress of the question request
# question_progresses: Dict[str, Progress] = {}

class Math(ProgressAbstract):
    def __init__(self, llm: str, storage: List[QuestionStorage], claims: JWTClaims):
        # setup the progress
        super().__init__(channel='post-question-progress', progress_ctx='post_question', to=None)
        self.llm = llm
        self.storage = storage
        [self.claims, self.request] = claims
        self.__question: Question = None
        self.__storage_response: List[Any] = []

    async def create_question(self, trees: List[Tree]):
        self._update_progress("Preparing", 0)
        [subject, chapter, topic] = trees
        llm = llm_models[self.llm]
        callback = llm.callbacks[0]
        # llm = ChatAnthropic(model='claude-3-5-sonnet-20240620', callbacks=[AgentCallbackHandler()])
        # Ask LLM to generate a question
        self._update_progress("Asking", 30)
        try:
            question = get_question(llm, subject, chapter, topic)
            self.__question = Question(**question.dict())
            self.__question.is_active = True
        except Exception as e:
            self._done()
            print(e)
            self.__question = Question(**question.dict())
            self.__question.error = str(e)
            self.__question.error_prompt = callback.prompt
            self.__question.is_active = False

        # self.__question.error_prompt = callback.prompt
        # self.__question.error_result_from_ai = callback.response
        self.__question.llm = self.llm
        self.__question.slugs = [
            subject.attributes.slug,
            chapter.attributes.slug,
            topic.attributes.slug,
        ]
        self.__question.created_by = self.claims.sub
        self.__question.organization_id = self.claims.org_id
        self.__question.created_at = datetime.now()
        self.__question.view_count = 0
        self.__question.rate_count = 0
        self.__question.rate_total = 0
        self.__question.claim_total = 0
        self.__question.claims = []
        # save to storage
        # have to save before parsing diagram order to get the inserted id
        self._update_progress("Saving", 60)
        self.__save()
        # parse diagram
        try:
            self._update_progress("Diagram", 80)
            self.__parse_diagram()
        except Exception as e:
            self._done()
            for res in self.__storage_response:
                if res != None:
                  collection = get_mongo_client()["questions"]
                  collection.update_one({"_id": res.inserted_id}, {"$set": {"is_active": False, "error": str(e)}})
        self._done()
        return self

    def get_tree_details(self, slug: str) -> any:
        pass
    
    def __parse_diagram(self) -> bool:
        if self.__question.diagram == "" or self.__question.diagram == None:
            return False
        ts = datetime.now().timestamp()
        file = md5(str(ts).encode()).hexdigest()
        path = os.path.join('diagrams', f'{file}.png')
        code = self.__question.diagram.replace("#python!", "").replace("plt.title", "#plt.title").replace("plt.show()", f"plt.savefig('{path}')")
        exec(code)
        with open(path, 'rb') as f:
            headers = {
                'Authorization': self.request.headers['authorization']
            }
            file = {'file': f}
            host = os.getenv('FILE_SERVICE_ADDR')
            for s in self.__storage_response:
              if s == None:
                  continue
              data = {
                  'service': 'content_service',
                  'serviceCtx': CTX_MATH_DIAGRAM,
                  'formID': s.inserted_id
              }
              response = requests.post(f'http://{host}/upload', files=file, data=data, headers=headers)
        return True
    
    def __save(self) -> Question:
        for storage in self.storage:
            if type(storage) == type:
                s = storage()
            else:
                s = storage
            res = s.save(self.__question)
            self.__storage_response.append(res)

class CustomBaseCallbackHandler(BaseCallbackHandler):
    prompt: str
    response: str

# To debugging the response from LLM
class AgentCallbackHandler(CustomBaseCallbackHandler):
    def on_llm_start(
        self, serialized: Dict[str, Any], prompts: List[str], **kwargs: Any
    ) -> Any:
        # """Run when LLM starts running."""
        # print(f"***Prompt to LLM was:***\n{prompts[0]}")
        # print("*********")
        self.prompt = prompts[0]

    def on_llm_end(self, response: LLMResult, **kwargs: Any) -> Any:
        # """Run when LLM ends running."""
        # print(f"***LLM Response:***\n{response.generations[0][0].text}")
        # print("*********")
        self.response = response.generations[0][0].text


llm_models = {
    "openai": ChatOpenAI(model="gpt-4o-2024-08-06", callbacks=[AgentCallbackHandler()]),
    "claude": ChatAnthropic(
        model="claude-3-5-sonnet-20240620", callbacks=[AgentCallbackHandler()]
    ),
}

def get_question(llm: str, subject: str, chapter: str, topic: str) -> QuestionOutput:
    template = """
    You are a {subject} teacher. Create a multiple choice question from {subject} > chapter: {chapter} > topic : {topic}.
    This question and choices will be rendered by Python Streamlit. Please follow the format instructions below.

    - Every equations in question, description and choices should be Latex format or Python matplotlib script.
    - If Latex is included \includegraphics(graphs and diagrams) then use Python matplotlib to create it.
    - When creating a diagram, make sure to check the Python code multiple times for any issues. The diagram is a very important part.
    - If equation is Python matplotlib, then add #python! to the beginning and end of the equation.
    - Equation in LaTeX format should be with $ to the beginning and end of the equation.
    - If possible, please create a question rich in diagrams and graphs.
    - Provide 4 multiple choice options in LaTeX format/Python matplotlib script without choice number.
    - Answer should be number and answer value should be equation in LaTeX format/Python matplotlib format.
    - Explain in detail how to solve the problem by steps.
    - Verify that the answer you selected is correct to ensure accuracy.
    - If the answer is a fraction, write it as a fraction, not a decimal.
    - If question requires a diagram, provide a diagram Python matplotlib format in single line.
    - If a \\ is included in Latex, convert it to \\\\. For example "$z = \pm 1$" to "$z = \\\\pm 1$".
    - If the equation in answer or choices are fractions, express them in their simplest form as much as possible. For example 2/4 to 1/2, 13/2 to 6 1/2.
    - All python Python matplotlib script should be "import matplotlib.pyplot as plt" and end with "plt.show()"
    - When creating JSON, make it without 'properties'. Ignore anything starting with 'properties' in the example I gave.

    \n{format_instruction}
    """
    summary_parser = PydanticOutputParser(pydantic_object=QuestionOutput)
    prompt = PromptTemplate(
        template=template,
        input_variables=["subject", "chapter", "topic"],
        partial_variables={
            "format_instruction": summary_parser.get_format_instructions()
        },
    )
    chain = prompt | llm | summary_parser
    output = chain.invoke(
        input={"subject": "algebra", "chapter": chapter, "topic": topic}
    )
    return output

####################

class TestOutput(BaseModel):
    answer: str = Field(description="Your answer")
    llm: str

def get_test(llm: str, subject: str, chapter: str, topic: str) -> Question:
    template = """
    What do you think this is? "{subject} / {subject} / {chapter} / {topic}"
    \n{format_instruction}
    """
    summary_parser = PydanticOutputParser(pydantic_object=TestOutput)
    prompt = PromptTemplate(
        template=template,
        input_variables=["subject", "chapter", "topic"],
        partial_variables={
            "format_instruction": summary_parser.get_format_instructions()
        },
    )
    chain = prompt | llm | summary_parser
    output = chain.invoke(
        input={"subject": "algebra", "chapter": chapter, "topic": topic}
    )
    return output

####################

class Output(BaseModel):
    question: str = Field(
        description="Question to be asked with LaTeX/Python(matplotlib)  formatted equation"
    )
    equation: str = Field(
        description="LaTeX/Python(matplotlib) formatted equation from question",
        default="",
    )
    choices: List[str] = Field(
        description="If choice is equation, it must be LaTeX/Python(matplotlib)  formatted. No choice number"
    )
    desc: List[str] = Field(
        description="Very detail explanation of answer by steps. \ should be replaced with \\."
    )
    answer_key: str = Field(
        description="choice number of answer, such as 1 or 2 or 3 or 4"
    )
    answer_val: str = Field(
        description="LaTeX/Python(matplotlib) formatted answer", default=""
    )
    diagram: str = Field(
        description="LaTeX/Python(matplotlib) formatted diagram or graph if requires, If no diagram and graph, leave it as empty string ('')",
        default="",
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
    



        
