import random
import string
import sys

import grpc
from locust import User, between, task, TaskSet, constant

sys.path.append("/home/alireza/Developer/goland/WP-Hw1/proto")
import auth_pb2 as pb2
import auth_pb2_grpc as pb2_grpc


class GrpcUser(User):
    wait_time = between(1, 5)

    def on_start(self):
        channel = grpc.insecure_channel('localhost:5052')  # Replace with your gRPC server address and port
        self.stub = pb2_grpc.AuthenticatorStub(channel)

    @task
    def request_pq(self):
        request = pb2.PQRequest()
        request.nonce = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
        request.message_id = random.randint(1, 100)

        response = self.stub.RequestPQ(request)

        print(response)

    @task
    def request_dh_params(self):
        request = pb2.DHRequest()
        request.nonce = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
        request.server_nonce = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
        request.message_id = random.randint(1, 100)
        request.a = random.randint(1, 100)

        response = self.stub.RequestDHParams(request)

        print(response)

    @task
    def auth_check(self):
        request = pb2.ACRequest()
        request.nonce = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
        request.server_nonce = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
        request.message_id = random.randint(1, 100)
        request.auth_key = random.randint(1, 100)

        response = self.stub.AuthCheck(request)

        print(response)


class GrpcUserTasks(TaskSet):
    tasks = {GrpcUser: 1}


class GrpcUserTestRunner(TaskSet):
    task_set = GrpcUserTasks
    wait_time = constant(0)


class GrpcUserLocust(User):
    host = ''
    task_set = GrpcUserTestRunner

    def __init__(self):
        super().__init__()
        self.client = None

    def on_start(self):
        self.client = grpc.insecure_channel('localhost:5052')  # Replace with your gRPC server address and port

    def on_stop(self):
        if self.client:
            self.client.close()
