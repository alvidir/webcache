import { Server, ServerUnaryCall, sendUnaryData, ServerCredentials } from 'grpc';
import { IApiServer, ApiService } from '../proto/api_grpc_pb';
import * as pb from '../proto/api_pb';

class ApiImplementation implements IApiServer {
    public single(call: ServerUnaryCall<pb.SingleRequest>, callback: sendUnaryData<pb.Result>): void {
        console.log("Single call from grpc");
    }

    public roll(call: ServerUnaryCall<pb.RollRequest>, callback: sendUnaryData<pb.Result>): void {
        console.log("Roll call from grpc");
    }
}

//create a server object
const server = new Server();
server.addService<IApiServer>(ApiService, new ApiImplementation());

export function GetInstance(): Server {
    return server;
}

const credential = ServerCredentials.createInsecure();
export function GetCredential(): ServerCredentials {
    return credential;
}