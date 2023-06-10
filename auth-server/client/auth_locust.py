from locust import User, task, between
import grpc
from proto import auth_pb2, auth_pb2_grpc


class GrpcUser(User):
    host = "localhost"
    port = 2884
    wait_time = between(1, 2)

    def on_start(self):
        self.client = AuthenticatorGrpcClient(self.host, self.port)

    @task
    def grpc_request_task(self):
        nonce = "nonce"
        server_nonce = "server_nonce"
        message_id = 123
        a = 789
        auth_key = 123

        # Make gRPC request using the client
        response = self.client.request_pq(nonce, message_id)
        print(response)

        response = self.client.request_dh_params(nonce, server_nonce, message_id, a)
        print(response)

        response = self.client.auth_check(nonce, server_nonce, message_id, auth_key)
        print(response)


class AuthenticatorGrpcClient:
    def __init__(self, host, port):
        self.channel = grpc.insecure_channel(f"{host}:{port}")
        self.stub = auth_pb2_grpc.AuthenticatorStub(self.channel)

    def request_pq(self, nonce, message_id):
        request = auth_pb2.PQRequest(
            nonce=nonce,
            message_id=message_id
        )
        response = self.stub.RequestPQ(request)
        return response

    def request_dh_params(self, nonce, server_nonce, message_id, a):
        request = auth_pb2.DHRequest(
            nonce=nonce,
            server_nonce=server_nonce,
            message_id=message_id,
            a=a
        )
        response = self.stub.RequestDHParams(request)
        return response

    def auth_check(self, nonce, server_nonce, message_id, auth_key):
        request = auth_pb2.ACRequest(
            nonce=nonce,
            server_nonce=server_nonce,
            message_id=message_id,
            auth_key=auth_key
        )
        response = self.stub.AuthCheck(request)
        return response
