// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var result_pb = require('./result_pb.js');
var roll_pb = require('./roll_pb.js');
var single_pb = require('./single_pb.js');
var supply_pb = require('./supply_pb.js');

function serialize_model_Result(arg) {
  if (!(arg instanceof result_pb.Result)) {
    throw new Error('Expected argument of type model.Result');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_model_Result(buffer_arg) {
  return result_pb.Result.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_model_RollRequest(arg) {
  if (!(arg instanceof roll_pb.RollRequest)) {
    throw new Error('Expected argument of type model.RollRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_model_RollRequest(buffer_arg) {
  return roll_pb.RollRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_model_SingleRequest(arg) {
  if (!(arg instanceof single_pb.SingleRequest)) {
    throw new Error('Expected argument of type model.SingleRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_model_SingleRequest(buffer_arg) {
  return single_pb.SingleRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_model_SupplyRequest(arg) {
  if (!(arg instanceof supply_pb.SupplyRequest)) {
    throw new Error('Expected argument of type model.SupplyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_model_SupplyRequest(buffer_arg) {
  return supply_pb.SupplyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var ApiService = exports.ApiService = {
  single: {
    path: '/model.Api/Single',
    requestStream: false,
    responseStream: false,
    requestType: single_pb.SingleRequest,
    responseType: result_pb.Result,
    requestSerialize: serialize_model_SingleRequest,
    requestDeserialize: deserialize_model_SingleRequest,
    responseSerialize: serialize_model_Result,
    responseDeserialize: deserialize_model_Result,
  },
  roll: {
    path: '/model.Api/Roll',
    requestStream: false,
    responseStream: false,
    requestType: roll_pb.RollRequest,
    responseType: result_pb.Result,
    requestSerialize: serialize_model_RollRequest,
    requestDeserialize: deserialize_model_RollRequest,
    responseSerialize: serialize_model_Result,
    responseDeserialize: deserialize_model_Result,
  },
  supply: {
    path: '/model.Api/Supply',
    requestStream: false,
    responseStream: false,
    requestType: supply_pb.SupplyRequest,
    responseType: result_pb.Result,
    requestSerialize: serialize_model_SupplyRequest,
    requestDeserialize: deserialize_model_SupplyRequest,
    responseSerialize: serialize_model_Result,
    responseDeserialize: deserialize_model_Result,
  },
};

exports.ApiClient = grpc.makeGenericClientConstructor(ApiService);
