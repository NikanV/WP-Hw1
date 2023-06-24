import random
from locust import User, task, between
import grpc
import sys
sys.path.append("../../proto")
import biz_pb2
import biz_pb2_grpc

class GrpcInterceptor(grpc.UnaryUnaryClientInterceptor):
    def intercept_unary_unary(self, continuation, client_call_details, request):
        print(f"Intercepting gRPC method: {client_call_details.method}")
        print(f"Request: {request}")

        response = continuation(client_call_details, request)

        print(f"Response: {response}")

        return response

class GrpcClient:
    def __init__(self, host):
        interceptors = [GrpcInterceptor()]
        channel = grpc.insecure_channel(host)
        self.channel = grpc.intercept_channel(channel, *interceptors)
        self.stub = biz_pb2_grpc.BizServiceStub(self.channel)

    def get_users(self, user_id, auth_key, message_id):
        request = biz_pb2.GetUsersRequest(
            user_id=user_id,
            auth_key=auth_key,
            message_id=message_id
        )
        return self.stub.GetUsers(request)

class GrpcLocust(User):
    host = "localhost:5062"
    wait_time = between(1, 5)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.client = GrpcClient(self.host)

    def get_odd(self):
        num = random.randint(1, 100)
        while num % 2 == 0:
            num = random.randint(1, 100)
        return num

    @task
    def get_users(self):
        user_id = random.choice([5263, 5303, 9649])
        auth_key = "authKey"
        message_id = self.get_odd()

        response = self.client.get_users(user_id, auth_key, message_id)
        print(f"Received response: {response}")
