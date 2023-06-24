import random
from locust import User, task, between
import grpc
import sys
from typing import Any, Callable
from locust.exception import LocustError
import grpc.experimental.gevent as grpc_gevent
sys.path.append("../../proto")
import auth_pb2 as authenticator_pb2
import auth_pb2_grpc as authenticator_pb2_grpc

class GrpcInterceptor(grpc.UnaryUnaryClientInterceptor):
    def intercept_unary_unary(self, continuation, client_call_details, request):
        print(f"Intercepting gRPC method: {client_call_details.method}")
        print(f"Request: {request}")

        response = continuation(client_call_details, request)

        print(f"Response: {response}")

        return response

class GrpcClient:
    def __init__(self, host):
        self.host = host

    def init(self):
        interceptors = [GrpcInterceptor()]
        channel = grpc.insecure_channel(self.host)
        self.channel = grpc.intercept_channel(channel, *interceptors)
        self.stub = authenticator_pb2_grpc.AuthenticatorStub(self.channel)

    def request_pq(self, nonce, message_id):
        request = authenticator_pb2.PQRequest(
            nonce=nonce,
            message_id=message_id
        )
        return self.stub.RequestPQ(request)

    def request_dh_params(self, nonce, server_nonce, message_id, a):
        request = authenticator_pb2.DHRequest(
            nonce=nonce,
            server_nonce=server_nonce,
            message_id=message_id,
            a=a
        )
        return self.stub.RequestDHParams(request)

    def auth_check(self, message_id, auth_key):
        request = authenticator_pb2.ACRequest(
            message_id=message_id,
            auth_key=auth_key
        )
        return self.stub.AuthCheck(request)

class GrpcLocust(User):
    host = "localhost:5052"
    wait_time = between(1, 5)

    def on_start(self):
        self.client = GrpcClient(self.host)
        self.client.init()

    def generate_nonce(self, max_width=20):
        return ''.join(random.choice('abcdefghijklmnopqrstuvwxyz') for _ in range(max_width))

    def generate_message_id(self):
        return random.randint(2, 1000000) * 2

    @task
    def request_pq(self):
        nonce = self.generate_nonce()
        message_id = self.generate_message_id()

        response = self.client.request_pq(nonce, message_id)
        print(f"Received response: {response}")

    @task
    def request_dh_params(self):
        nonce = self.generate_nonce()
        server_nonce = self.generate_nonce()
        message_id = self.generate_message_id()
        a = random.randint(1, 100)

        response = self.client.request_dh_params(nonce, server_nonce, message_id, a)
        print(f"Received response: {response}")

    @task
    def auth_check(self):
        message_id = self.generate_message_id()
        auth_key = random.randint(1, 1000)

        response = self.client.auth_check(message_id, auth_key)
        print(f"Received response: {response}")

grpc_gevent.init_gevent()
