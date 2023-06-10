import grpc
from proto import auth_pb2, auth_pb2_grpc


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
