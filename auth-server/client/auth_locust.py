import random
import string
import sys
import interceptor
from locust import task

sys.path.append("../../proto")
import auth_pb2 as pb2
import auth_pb2_grpc as pb2_grpc


class AuthServiceGrpcUser(interceptor.GrpcUser):
    host = "localhost:5052"
    stub_class = pb2_grpc.AuthenticatorStub

    @task
    def requestPQ(self):
        self.stub.RequestPQ(pb2.PQRequest(message_id=4, nonce="abcdabcdabcdabcdabcd"))

    # @task
    # def requestDHParams(self):
    #     self.stub.RequestDHParams(pb2.DHRequest(message_id=4, nonce="abcdabcdabcdabcdabcd", server_nonce="efghefghefghefghefgh", a=4))
        
    @task
    def authCheck(self):
        self.stub.AuthCheck(pb2.ACRequest(message_id=4, auth_key=15))
