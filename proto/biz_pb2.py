# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: biz.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\tbiz.proto\"H\n\x0fGetUsersRequest\x12\x0f\n\x07user_id\x18\x01 \x01(\x03\x12\x10\n\x08\x61uth_key\x18\x02 \x01(\t\x12\x12\n\nmessage_id\x18\x03 \x01(\x03\"O\n\x16GetUsersWithSQLRequest\x12\x0f\n\x07user_id\x18\x01 \x01(\t\x12\x10\n\x08\x61uth_key\x18\x02 \x01(\t\x12\x12\n\nmessage_id\x18\x03 \x01(\x03\"<\n\x10GetUsersResponse\x12\x14\n\x05users\x18\x01 \x03(\x0b\x32\x05.USER\x12\x12\n\nmessage_id\x18\x02 \x01(\x03\"]\n\x04USER\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0e\n\x06\x66\x61mily\x18\x02 \x01(\t\x12\n\n\x02id\x18\x03 \x01(\x03\x12\x0b\n\x03\x61ge\x18\x04 \x01(\x03\x12\x0b\n\x03sex\x18\x05 \x01(\t\x12\x11\n\tcreatedAt\x18\x06 \x01(\t2\x80\x01\n\nBizService\x12\x31\n\x08GetUsers\x12\x10.GetUsersRequest\x1a\x11.GetUsersResponse\"\x00\x12?\n\x0fGetUsersWithSQL\x12\x17.GetUsersWithSQLRequest\x1a\x11.GetUsersResponse\"\x00\x42\x04Z\x02./b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'biz_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\002./'
  _GETUSERSREQUEST._serialized_start=13
  _GETUSERSREQUEST._serialized_end=85
  _GETUSERSWITHSQLREQUEST._serialized_start=87
  _GETUSERSWITHSQLREQUEST._serialized_end=166
  _GETUSERSRESPONSE._serialized_start=168
  _GETUSERSRESPONSE._serialized_end=228
  _USER._serialized_start=230
  _USER._serialized_end=323
  _BIZSERVICE._serialized_start=326
  _BIZSERVICE._serialized_end=454
# @@protoc_insertion_point(module_scope)
