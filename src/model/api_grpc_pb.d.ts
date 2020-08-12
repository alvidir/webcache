// package: model
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as api_pb from "./api_pb";
import * as result_pb from "./result_pb";
import * as roll_pb from "./roll_pb";
import * as single_pb from "./single_pb";
import * as supply_pb from "./supply_pb";

interface IApiService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    single: IApiService_ISingle;
    roll: IApiService_IRoll;
    supply: IApiService_ISupply;
}

interface IApiService_ISingle extends grpc.MethodDefinition<single_pb.SingleRequest, result_pb.Result> {
    path: string; // "/model.Api/Single"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<single_pb.SingleRequest>;
    requestDeserialize: grpc.deserialize<single_pb.SingleRequest>;
    responseSerialize: grpc.serialize<result_pb.Result>;
    responseDeserialize: grpc.deserialize<result_pb.Result>;
}
interface IApiService_IRoll extends grpc.MethodDefinition<roll_pb.RollRequest, result_pb.Result> {
    path: string; // "/model.Api/Roll"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<roll_pb.RollRequest>;
    requestDeserialize: grpc.deserialize<roll_pb.RollRequest>;
    responseSerialize: grpc.serialize<result_pb.Result>;
    responseDeserialize: grpc.deserialize<result_pb.Result>;
}
interface IApiService_ISupply extends grpc.MethodDefinition<supply_pb.SupplyRequest, result_pb.Result> {
    path: string; // "/model.Api/Supply"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<supply_pb.SupplyRequest>;
    requestDeserialize: grpc.deserialize<supply_pb.SupplyRequest>;
    responseSerialize: grpc.serialize<result_pb.Result>;
    responseDeserialize: grpc.deserialize<result_pb.Result>;
}

export const ApiService: IApiService;

export interface IApiServer {
    single: grpc.handleUnaryCall<single_pb.SingleRequest, result_pb.Result>;
    roll: grpc.handleUnaryCall<roll_pb.RollRequest, result_pb.Result>;
    supply: grpc.handleUnaryCall<supply_pb.SupplyRequest, result_pb.Result>;
}

export interface IApiClient {
    single(request: single_pb.SingleRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    single(request: single_pb.SingleRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    single(request: single_pb.SingleRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: roll_pb.RollRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: roll_pb.RollRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    roll(request: roll_pb.RollRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: supply_pb.SupplyRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: supply_pb.SupplyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    supply(request: supply_pb.SupplyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
}

export class ApiClient extends grpc.Client implements IApiClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public single(request: single_pb.SingleRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public single(request: single_pb.SingleRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public single(request: single_pb.SingleRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: roll_pb.RollRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: roll_pb.RollRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public roll(request: roll_pb.RollRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: supply_pb.SupplyRequest, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: supply_pb.SupplyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
    public supply(request: supply_pb.SupplyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: result_pb.Result) => void): grpc.ClientUnaryCall;
}
