import { Server, ServerUnaryCall, sendUnaryData, ServerCredentials, credentials, ChannelCredentials } from 'grpc';
import { IApiServer, ApiService } from '../../proto/api_grpc_pb';
import * as api_pb from '../../proto/api_pb';
import { GetInstance as CacheInstance } from '../../cache/public';
import { ProtobufAdapter } from './adapters';

class ApiImplementation implements IApiServer {
    public single(_: ServerUnaryCall<api_pb.EmptyRequest>, callback: sendUnaryData<api_pb.Image>): void {
        const image = CacheInstance().GetSingle();
        const value = ProtobufAdapter(image);
        callback(null, value);
    }
}

//create a server object
const server = new Server();
server.addService<IApiServer>(ApiService, new ApiImplementation());

export function GetInstance(): Server {
    return server;
}

const srv_credential = ServerCredentials.createInsecure();
export function GetServerCredential(): ServerCredentials {
    return srv_credential;
}

const ch_credential = credentials.createInsecure();
export function GetChannelCredential(): ChannelCredentials {
    return ch_credential;
}