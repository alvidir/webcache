// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var api_pb = require('./api_pb.js');

function serialize_proto_EmptyRequest(arg) {
  if (!(arg instanceof api_pb.EmptyRequest)) {
    throw new Error('Expected argument of type proto.EmptyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_EmptyRequest(buffer_arg) {
  return api_pb.EmptyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_proto_Image(arg) {
  if (!(arg instanceof api_pb.Image)) {
    throw new Error('Expected argument of type proto.Image');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_Image(buffer_arg) {
  return api_pb.Image.deserializeBinary(new Uint8Array(buffer_arg));
}


// API Service
var ApiService = exports.ApiService = {
  single: {
    path: '/proto.Api/Single',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.EmptyRequest,
    responseType: api_pb.Image,
    requestSerialize: serialize_proto_EmptyRequest,
    requestDeserialize: deserialize_proto_EmptyRequest,
    responseSerialize: serialize_proto_Image,
    responseDeserialize: deserialize_proto_Image,
  },
};

exports.ApiClient = grpc.makeGenericClientConstructor(ApiService);
