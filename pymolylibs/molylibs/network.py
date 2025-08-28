from datetime import datetime
import os
import jwt
import json
import asyncio
from abc import ABC, abstractmethod


from fastapi import Request
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware
from typing import List, Dict, Any, Tuple, Union, Optional
from pydantic import BaseModel, Field
from fastapi import Request, WebSocket

import molylibs.pb.message_service_pb2 as pb
from molylibs.grpc import send_progress


def verify_jwt(token: str):
    jwt_secret = os.getenv('JWT_AUTH_SECRET')
    try:
        return jwt.decode(token.replace('Bearer ', ''), jwt_secret, algorithms=['HS256'], options={"verify_aud": False})
    except jwt.ExpiredSignatureError:
        return 'Signature has expired'
    except jwt.InvalidTokenError:
        return 'Invalid token'

class JWTMiddlewear(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        token = request.headers.get('Authorization')
        if token:
            try:
                payload = verify_jwt(token)
                request.state.payload = payload
            except Exception as e:
                return JSONResponse(status_code=401, content={'message': str(e)})
        response = await call_next(request)
        return response

class JWTClaims(BaseModel):
    email: Optional[str] = Field(default='')
    first_name: Optional[str] = Field(alias='firstName', default='')
    last_name: Optional[str] = Field(alias='lastName',default='')
    org_id: Optional[str] = Field(alias='orgId', default='')
    profile_image: Optional[str] = Field(alias='profileImage', default='')
    role: Optional[str] = Field(default='')
    metadata: Optional[Dict[str, Any]] = Field(default={})
    ip: Optional[str] = Field(default='')
    iss: Optional[str] = Field(default='')
    sub: Optional[str] = Field(default='')
    aud: Optional[List[str]] = Field(default=[])
    exp: Optional[int] = Field(default=0)

JWTClaimsWithRequest = Tuple[JWTClaims, Request]

def get_jwt_from_request(request: Request) -> JWTClaimsWithRequest:
    try:
        return (JWTClaims(**request.state.payload), request)
    except Exception as e:
        return (JWTClaims(), request)
    #return (JWTClaims(**request.state.payload), request)

# Parse JWT claims from websocket headers
def get_claims_from_websocket(websocket: WebSocket) -> Union[JWTClaims, None]:
    if websocket.headers['authorization'] == "":
        return None
    payload = verify_jwt(websocket.headers['authorization'])
    return JWTClaims(**payload)

# WebSocket

# CHANNEL_ADMIN_QUESTION_PROGRESS = 'admin_question_progress'

class WebSocketPackage:
    def __init__(self, websocket: WebSocket, claims: JWTClaims):
        self.websocket = websocket
        self.claims = claims


class WebSocketConnectionManager:
    """
    Singleton class to manage websocket connections
    """
    _instance = None
    def __new__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super(WebSocketConnectionManager, cls).__new__(cls)
        return cls._instance
    
    def __init__(self):
        self.__connections: Dict[str, Dict[str, List[WebSocketPackage]]] = {}  
        pass

    def add(self, channel: str, id: str, package: WebSocketPackage):
        if channel not in self.__connections:
            self.__connections[channel] = {}
        self.__connections[channel][id] = package

    def remove(self, channel: str, id: str):
        del self.__connections[channel][id]

    def get(self, channel: str, id: str):
        return self.__connections[channel][id]
    
    def get_channel(self, channel: str):
        if channel not in self.__connections:
            return None
        return self.__connections[channel]

websocket_connections = WebSocketConnectionManager()

class WebSocketManager:
    #_instance = None

    # def __new__(cls, *args, **kwargs):
    #     if not cls._instance:
    #         cls._instance = super(WebSocketManager, cls).__new__(cls)
    #         #cls._instance.websocket_connections = []
    #     return cls._instance

    def __init__(self, channel: str, claims: JWTClaims = None, websocket: WebSocket = None):    
        if channel is None:
            raise ValueError("channel is required")
        self.channel = channel
        self.claims = claims
        self.websocket = websocket
        self.id = str(id(self))

    async def connect(self):
        if self.websocket is None:
            raise ValueError("Websocket is required")
        await self.websocket.accept()
        websocket_connections.add(self.channel, self.id, WebSocketPackage(websocket=self.websocket, claims=self.claims))
        # if self.channel not in websocket_connections:
        #     websocket_connections[self.channel] = {}
        # websocket_connections[self.channel][self.id] = WebSocketPackage(websocket=self.websocket, claims=self.claims)

    def disconnect(self):
        websocket_connections.remove(self.channel, self.id)
        #del websocket_connections[self.channel][self.id]

    async def send_message(self, message: str):
        await self.websocket.send_text(message)

    async def broadcast(self, message: str):
        if websocket_connections.get_channel(self.channel) is None:
            return
        for package in websocket_connections.get_channel(self.channel).values():
            try:
                await package.websocket.send_text(message)
            except RuntimeError as e:
                print(f"Error broadcasting message: {e}")


class ProgressAbstract:
    """
    Class to manage progress of a task
    """
    #send_progress(channel="admin", context=b"Init!!!", progress=30, id="2", progress_ctx="post_question", command=pb.ProgressCommand.UPDATE, to="admin")

    def __init__(self, channel:str, progress_ctx: str, to: str):
        self.channel = channel
        self.progress_ctx = progress_ctx
        self.to = to
        self.id = str(id(self))
        self.__send_progress(pb.ProgressCommand.ADD, context="Initializing", progress=0)

    def __send_progress(self, command:int, context:str, progress: int): 
        send_progress(channel=self.channel, context=context, progress=progress, id=self.id, progress_ctx=self.progress_ctx, command=command, to=self.to)

    def _update_progress(self, context: str, progress: int):
        print("_update_progress")
        self.__send_progress(pb.ProgressCommand.UPDATE, context=context, progress=progress)

    def _done(self):
        print("Done")
        self.__send_progress(pb.ProgressCommand.DONE, context="Done", progress=100)



