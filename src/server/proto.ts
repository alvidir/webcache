import { Server, ServerUnaryCall, sendUnaryData, ServerCredentials } from 'grpc';
import { IApiServer, ApiService } from '../model/api_grpc_pb';
import { Result } from '../model/result_pb';
import { SingleRequest } from '../model/single_pb';
import { RollRequest } from '../model/roll_pb';
import { SupplyRequest } from '../model/supply_pb';

function singleImplementation(req: SingleRequest): Result {
    return new Result();
}

class ApiImplementation implements IApiServer {
    public single(call: ServerUnaryCall<SingleRequest>, callback: sendUnaryData<Result>): void {
        console.log("Single call from grpc");
    }

    public roll(call: ServerUnaryCall<RollRequest>, callback: sendUnaryData<Result>): void {
        console.log("Roll call from grpc");
    }

    public supply(call: ServerUnaryCall<SupplyRequest>, callback: sendUnaryData<Result>): void {
        console.log("Supply call from grpc");
    }
}

//create a server object
const server = new Server();
server.addService(ApiService, new ApiImplementation());

export function GetInstance(): Server {
    return server;
}

const credential = ServerCredentials.createInsecure();
export function GetCredential(): ServerCredentials {
    return credential;
}