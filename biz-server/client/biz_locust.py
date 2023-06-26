import random
import string
import sys

import grpc
from locust import User, between, task, TaskSet, constant

sys.path.append("../../proto")
import biz_pb2 as pb2
import biz_pb2_grpc as pb2_grpc


class GrpcUser(User):
    wait_time = between(1, 5)
    num_requests = 0

    def get_odd(self):
        num = random.randint(1, 100)
        while num % 2 == 0:
            num = random.randint(1, 100)
        return num

    def on_start(self):
        channel = grpc.insecure_channel('localhost:5062')
        self.stub = pb2_grpc.BizServiceStub(channel)

    @task
    def get_users(self):
        request = pb2.GetUsersRequest()
        request.user_id = random.choice([5263, 5303, 9649])
        request.auth_key = 'authkey'
        request.message_id = self.get_odd()

        response = self.stub.GetUsers(request)
        self.num_requests += 1

        print(f"Number of requests: {self.num_requests}")

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
        self.client = grpc.insecure_channel('localhost:5062')

    def on_stop(self):
        if self.client:
            self.client.close()

