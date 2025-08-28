from os import getenv
from typing import List, Any
from urllib import request
import json

from grpc import insecure_channel

from httpx import get
from molylibs.schemas import Tree
from molylibs.pb import tree_service_pb2, tree_service_pb2_grpc
from molylibs.pb import message_service_pb2, message_service_pb2_grpc

service_name = getenv('SERVICE_NAME')

class Client:
    def __init__(self, addr: str, stub, request, service, port: str = ':5000'):
        self.addr = getenv(addr)
        self.port = port
        self.stub = stub
        self.service = service
        self.request = request

    def fetch(self) -> Any:
        with insecure_channel(self.addr + self.port) as channel:
            stub = self.stub(channel)
            service_method = getattr(stub, self.service)
            response = service_method(self.request)
            return response

def get_tree_ancestors(org_id: str, role: str, slug: str) -> List[Tree]:
    request = tree_service_pb2.TreeRequest(role=role, orgID=org_id, slug=slug)
    client = Client('TREE_SERVICE_ADDR', tree_service_pb2_grpc.TreeServiceStub, request, 'GetAncestors')
    output: List[Tree] = []
    data = json.loads(client.fetch().data)
    for tree in data:
        output.append(Tree(**tree))
    return output

def get_tree(org_id: str, role: str, slug: str) -> List[Tree]:
    request = tree_service_pb2.TreeRequest(role=role, orgID=org_id, slug=slug)
    client = Client('TREE_SERVICE_ADDR', tree_service_pb2_grpc.TreeServiceStub, request, 'GetTree')
    output: List[Tree] = []
    data = json.loads(client.fetch().data)
    return data

def get_tree_as_array(org_id: str, role: str, slug: str) -> List[Tree]:
    request = tree_service_pb2.TreeRequest(role=role, orgID=org_id, slug=slug)
    client = Client('TREE_SERVICE_ADDR', tree_service_pb2_grpc.TreeServiceStub, request, 'GetTreeAsArray')
    data = json.loads(client.fetch().data)
    return data

def is_end(org_id: str, role: str, slug: str) -> bool:
    request = tree_service_pb2.TreeRequest(role=role, orgID=org_id, slug=slug)
    client = Client('TREE_SERVICE_ADDR', tree_service_pb2_grpc.TreeServiceStub, request, 'IsEnd')
    return client.fetch().data

def broadcast_message(channel: str, message: str, message_type: message_service_pb2.MessageType) -> bool:
    request = message_service_pb2.MessageRequest(channel=channel, message=message, service=service_name)
    client = Client('MESSAGE_SERVICE_ADDR', message_service_pb2_grpc.MessageServiceStub, request, 'Broadcast')
    data = client.fetch().success
    return data

def send_progress(channel: str, progress: int, context: str, id: str, command:message_service_pb2.ProgressCommand, to: str, progress_ctx: str) -> bool:
    progress_message = message_service_pb2.ProgressMessage(progress=progress, context=context)
    request = message_service_pb2.ProgressMessageRequest(channel=channel, progress=progress_message, id=id, progressCtx=progress_ctx, service=service_name, to=to, command=command)
    client = Client('MESSAGE_SERVICE_ADDR', message_service_pb2_grpc.MessageServiceStub, request, 'SendProgress')
    data = client.fetch().success
    return data



