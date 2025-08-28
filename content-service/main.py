import multiprocessing
import threading

import fastapi
from grpc_server import start_grpc_server
from fastapi_server import start_fastapi
from grpc_server import start_grpc_server

fastapi_process = None
grpc_process = None

if __name__ == "__main__":
    print("Starting FastAPI and gRPC servers")
    # #start_grpc_server()
    # print("Starting FastAPI and gRPC servers")
    # start_fastapi()
    # # pass
    # print("Starting FastAPI and gRPC servers")
    fastapi_process = multiprocessing.Process(target=start_fastapi)
    grpc_process = multiprocessing.Process(target=start_grpc_server)

    fastapi_process.start()
    grpc_process.start()

    fastapi_process.join()
    grpc_process.join()