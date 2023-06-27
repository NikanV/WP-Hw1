import sys
import interceptor
from locust import task

sys.path.append("../../proto")
import biz_pb2 as pb2
import biz_pb2_grpc as pb2_grpc


class BizServiceGrpcUser(interceptor.GrpcUser):
    host = "localhost:5062"
    stub_class = pb2_grpc.BizServiceStub

    @task
    def getUsers(self):
        self.stub.GetUsers(pb2.GetUsersRequest(user_id=5303, message_id=4, auth_key="15"))

    @task
    def getUsersWithSQL(self):
        self.stub.GetUsersWithSQL(pb2.GetUsersWithSQLRequest(user_id="5303", message_id=4, auth_key="15"))