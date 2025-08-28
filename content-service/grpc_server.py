import grpc
import os
from concurrent import futures
import molylibs.pb.file_client_pb2_grpc as file_client_pb2_grpc
import molylibs.pb.file_client_pb2 as file_client_pb2
from molylibs.db import get_mongo_client
from bson import ObjectId

CTX_MATH_DIAGRAM = "math_diagrams"

class FileService(file_client_pb2_grpc.FileClientServicer):
    def GetFileRules(self, request, context):
        print("GetFileRules invoked")
        path = os.path.join(request.orgID, "math", "diagrams", "")
        return file_client_pb2.FileRulesResponse(s3Dir=path, maxFileSize=1000000, allowedFileTypes=["*"], allowedContentTypes=[])  

    def SaveFileData(self, request, context):
        print("SaveFileData invoked")
        with open("grpc_server.log", "a") as log_file:
            print("SaveFileData invoked", file=log_file)
            print("request", request, file=log_file)
            print('request.fileInfos[0]', request.fileInfos[0], file=log_file)
            # print('request.fileInfos.s3Key', request.fileInfos.s3Key, file=log_file)
            # print('request["fileInfos"]["s3Key"]', request["fileInfos"]["s3Key"], file=log_file)
            try:
                collection = get_mongo_client()["questions"]
                object_id = ObjectId(request.formID)
                collection.update_one({"_id": object_id}, {"$set": {"diagram_image": request.fileInfos[0].s3Key}})
            except Exception as e:
                print("error", str(e), file=log_file)
        return file_client_pb2.FileSaveDataResponse(output=b"success")

def start_grpc_server():
    print("starting grpc server")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    file_client_pb2_grpc.add_FileClientServicer_to_server(FileService(), server)
    server.add_insecure_port('[::]:5001')
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    start_grpc_server()