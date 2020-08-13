// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var api_pb = require('./api_pb.js');

function serialize_proto_Result(arg) {
  if (!(arg instanceof api_pb.Result)) {
    throw new Error('Expected argument of type proto.Result');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_Result(buffer_arg) {
  return api_pb.Result.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_proto_RollRequest(arg) {
  if (!(arg instanceof api_pb.RollRequest)) {
    throw new Error('Expected argument of type proto.RollRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_RollRequest(buffer_arg) {
  return api_pb.RollRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_proto_SingleRequest(arg) {
  if (!(arg instanceof api_pb.SingleRequest)) {
    throw new Error('Expected argument of type proto.SingleRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_SingleRequest(buffer_arg) {
  return api_pb.SingleRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_proto_SupplyRequest(arg) {
  if (!(arg instanceof api_pb.SupplyRequest)) {
    throw new Error('Expected argument of type proto.SupplyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_proto_SupplyRequest(buffer_arg) {
  return api_pb.SupplyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


// API Service
var ApiService = exports.ApiService = {
  single: {
    path: '/proto.Api/Single',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.SingleRequest,
    responseType: api_pb.Result,
    requestSerialize: serialize_proto_SingleRequest,
    requestDeserialize: deserialize_proto_SingleRequest,
    responseSerialize: serialize_proto_Result,
    responseDeserialize: deserialize_proto_Result,
  },
  roll: {
    path: '/proto.Api/Roll',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.RollRequest,
    responseType: api_pb.Result,
    requestSerialize: serialize_proto_RollRequest,
    requestDeserialize: deserialize_proto_RollRequest,
    responseSerialize: serialize_proto_Result,
    responseDeserialize: deserialize_proto_Result,
  },
  supply: {
    path: '/proto.Api/Supply',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.SupplyRequest,
    responseType: api_pb.Result,
    requestSerialize: serialize_proto_SupplyRequest,
    requestDeserialize: deserialize_proto_SupplyRequest,
    responseSerialize: serialize_proto_Result,
    responseDeserialize: deserialize_proto_Result,
  },
};

exports.ApiClient = grpc.makeGenericClientConstructor(ApiService);
