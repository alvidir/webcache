//import { Server, ServerUnaryCall, sendUnaryData, ServerCredentials } from 'grpc';
//import { ApiService, IApiService } from '../proto/proto-ts/api_grpc_pb';
//import { Result } from '../proto/proto-ts/result_pb';
//import { SingleRequest } from '../proto/proto-ts/single_pb';
//import { RollRequest } from '../proto/proto-ts/roll_pb';
//import { SupplyRequest } from '../proto/proto-ts/supply_pb';
//
//function singleImplementation(req: SingleRequest): Result {
//    return new Result();
//}
//
//class ApiImplementation implements IApiService {
//    public single(call: ServerUnaryCall<SingleRequest>, callback: sendUnaryData<Result>): void {
//        console.log("Single call from grpc");
//    }
//
//    public roll(call: ServerUnaryCall<RollRequest>, callback: sendUnaryData<Result>): void {
//        console.log("Roll call from grpc");
//    }
//
//    public supply(call: ServerUnaryCall<SupplyRequest>, callback: sendUnaryData<Result>): void {
//        console.log("Supply call from grpc");
//    }
//}
//
////create a server object
//const server = new Server();
//server.addService(ApiService, new ApiImplementation());
//
//export function GetInstance(): Server {
//    return server;
//}
//
//const credential = ServerCredentials.createInsecure();
//export function GetCredential(): ServerCredentials {
//    return credential;
//}