// package: proto
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as api_pb from "./api_pb";

interface IApiService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    single: IApiService_ISingle;
}

interface IApiService_ISingle extends grpc.MethodDefinition<api_pb.EmptyRequest, api_pb.Image> {
    path: string; // "/proto.Api/Single"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.EmptyRequest>;
    requestDeserialize: grpc.deserialize<api_pb.EmptyRequest>;
    responseSerialize: grpc.serialize<api_pb.Image>;
    responseDeserialize: grpc.deserialize<api_pb.Image>;
}

export const ApiService: IApiService;

export interface IApiServer {
    single: grpc.handleUnaryCall<api_pb.EmptyRequest, api_pb.Image>;
}

export interface IApiClient {
    single(request: api_pb.EmptyRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
    single(request: api_pb.EmptyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
    single(request: api_pb.EmptyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
}

export class ApiClient extends grpc.Client implements IApiClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public single(request: api_pb.EmptyRequest, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
    public single(request: api_pb.EmptyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
    public single(request: api_pb.EmptyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Image) => void): grpc.ClientUnaryCall;
}
