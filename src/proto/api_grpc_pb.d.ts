// package: proto
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as api_pb from "./api_pb";

interface IApiService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    single: IApiService_ISingle;
    roll: IApiService_IRoll;
    supply: IApiService_ISupply;
}

interface IApiService_ISingle extends grpc.MethodDefinition<api_pb.SingleRequest, api_pb.Result> {
    path: string; // "/proto.Api/Single"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.SingleRequest>;
    requestDeserialize: grpc.deserialize<api_pb.SingleRequest>;
    responseSerialize: grpc.serialize<api_pb.Result>;
    responseDeserialize: grpc.deserialize<api_pb.Result>;
}
interface IApiService_IRoll extends grpc.MethodDefinition<api_pb.RollRequest, api_pb.Result> {
    path: string; // "/proto.Api/Roll"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.RollRequest>;
    requestDeserialize: grpc.deserialize<api_pb.RollRequest>;
    responseSerialize: grpc.serialize<api_pb.Result>;
    responseDeserialize: grpc.deserialize<api_pb.Result>;
}
interface IApiService_ISupply extends grpc.MethodDefinition<api_pb.SupplyRequest, api_pb.Result> {
    path: string; // "/proto.Api/Supply"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.SupplyRequest>;
    requestDeserialize: grpc.deserialize<api_pb.SupplyRequest>;
    responseSerialize: grpc.serialize<api_pb.Result>;
    responseDeserialize: grpc.deserialize<api_pb.Result>;
}

export const ApiService: IApiService;

export interface IApiServer {
    single: grpc.handleUnaryCall<api_pb.SingleRequest, api_pb.Result>;
    roll: grpc.handleUnaryCall<api_pb.RollRequest, api_pb.Result>;
    supply: grpc.handleUnaryCall<api_pb.SupplyRequest, api_pb.Result>;
}

export interface IApiClient {
    single(request: api_pb.SingleRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    single(request: api_pb.SingleRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    single(request: api_pb.SingleRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: api_pb.RollRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: api_pb.RollRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: api_pb.RollRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: api_pb.SupplyRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: api_pb.SupplyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: api_pb.SupplyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
}

export class ApiClient extends grpc.Client implements IApiClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public single(request: api_pb.SingleRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public single(request: api_pb.SingleRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public single(request: api_pb.SingleRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: api_pb.RollRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: api_pb.RollRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: api_pb.RollRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: api_pb.SupplyRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: api_pb.SupplyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: api_pb.SupplyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Result) => void): grpc.ClientUnaryCall;
}
