// package: proto
// file: src/proto/api.proto

import * as src_proto_api_pb from "../../src/proto/api_pb";
import * as src_proto_result_pb from "../../src/proto/result_pb";
import * as src_proto_roll_pb from "../../src/proto/roll_pb";
import * as src_proto_single_pb from "../../src/proto/single_pb";
import * as src_proto_supply_pb from "../../src/proto/supply_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ApiSingle = {
  readonly methodName: string;
  readonly service: typeof Api;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof src_proto_single_pb.SingleRequest;
  readonly responseType: typeof src_proto_result_pb.Result;
};

type ApiRoll = {
  readonly methodName: string;
  readonly service: typeof Api;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof src_proto_roll_pb.RollRequest;
  readonly responseType: typeof src_proto_result_pb.Result;
};

type ApiSupply = {
  readonly methodName: string;
  readonly service: typeof Api;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof src_proto_supply_pb.SupplyRequest;
  readonly responseType: typeof src_proto_result_pb.Result;
};

export class Api {
  static readonly serviceName: string;
  static readonly Single: ApiSingle;
  static readonly Roll: ApiRoll;
  static readonly Supply: ApiSupply;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  single(
    requestMessage: src_proto_single_pb.SingleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
  single(
    requestMessage: src_proto_single_pb.SingleRequest,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
  roll(
    requestMessage: src_proto_roll_pb.RollRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
  roll(
    requestMessage: src_proto_roll_pb.RollRequest,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
  supply(
    requestMessage: src_proto_supply_pb.SupplyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
  supply(
    requestMessage: src_proto_supply_pb.SupplyRequest,
    callback: (error: ServiceError|null, responseMessage: src_proto_result_pb.Result|null) => void
  ): UnaryResponse;
}

