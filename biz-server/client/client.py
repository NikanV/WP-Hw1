import grpc
from proto import biz_pb2, biz_pb2_grpc


class BizGrpcClient:
    def __init__(self, host, port):
        self.channel = grpc.insecure_channel(f"{host}:{port}")
        self.stub = biz_pb2_grpc.BizServiceStub(self.channel)

    def get_users(self, user_id, auth_key, message_id):
        request = biz_pb2.GetUsersRequest(
            user_id=user_id,
            auth_key=auth_key,
            message_id=message_id
        )
        response = self.stub.GetUsers(request)
        return response

    def get_users_with_sql(self, user_id, auth_key, message_id):
        request = biz_pb2.GetUsersWithSQLRequest(
            user_id=user_id,
            auth_key=auth_key,
            message_id=message_id
        )
        response = self.stub.GetUsersWithSQL(request)
        return response
