/**
 * @fileoverview gRPC-Web generated client stub for proto
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.proto = require('./api_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.ApiClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.ApiPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.EmptyRequest,
 *   !proto.proto.Image>}
 */
const methodDescriptor_Api_Single = new grpc.web.MethodDescriptor(
  '/proto.Api/Single',
  grpc.web.MethodType.UNARY,
  proto.proto.EmptyRequest,
  proto.proto.Image,
  /**
   * @param {!proto.proto.EmptyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.Image.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.EmptyRequest,
 *   !proto.proto.Image>}
 */
const methodInfo_Api_Single = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.Image,
  /**
   * @param {!proto.proto.EmptyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.Image.deserializeBinary
);


/**
 * @param {!proto.proto.EmptyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.Image)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.Image>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ApiClient.prototype.single =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.Api/Single',
      request,
      metadata || {},
      methodDescriptor_Api_Single,
      callback);
};


/**
 * @param {!proto.proto.EmptyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.Image>}
 *     A native promise that resolves to the response
 */
proto.proto.ApiPromiseClient.prototype.single =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.Api/Single',
      request,
      metadata || {},
      methodDescriptor_Api_Single);
};


module.exports = proto.proto;

