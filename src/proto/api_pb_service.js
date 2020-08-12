// package: proto
// file: src/proto/api.proto

var src_proto_api_pb = require("../../src/proto/api_pb");
var src_proto_result_pb = require("../../src/proto/result_pb");
var src_proto_roll_pb = require("../../src/proto/roll_pb");
var src_proto_single_pb = require("../../src/proto/single_pb");
var src_proto_supply_pb = require("../../src/proto/supply_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Api = (function () {
  function Api() {}
  Api.serviceName = "proto.Api";
  return Api;
}());

Api.Single = {
  methodName: "Single",
  service: Api,
  requestStream: false,
  responseStream: false,
  requestType: src_proto_single_pb.SingleRequest,
  responseType: src_proto_result_pb.Result
};

Api.Roll = {
  methodName: "Roll",
  service: Api,
  requestStream: false,
  responseStream: false,
  requestType: src_proto_roll_pb.RollRequest,
  responseType: src_proto_result_pb.Result
};

Api.Supply = {
  methodName: "Supply",
  service: Api,
  requestStream: false,
  responseStream: false,
  requestType: src_proto_supply_pb.SupplyRequest,
  responseType: src_proto_result_pb.Result
};

exports.Api = Api;

function ApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ApiClient.prototype.single = function single(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Api.Single, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ApiClient.prototype.roll = function roll(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Api.Roll, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ApiClient.prototype.supply = function supply(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Api.Supply, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.ApiClient = ApiClient;

