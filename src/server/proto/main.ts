import * as Proto from '../../server/proto/grpc';
import environment from '../../config';

const app_host = environment.AppHost;

const proto_port = environment.ProtoServicePort;
const credential = Proto.GetServerCredential();
const prefix = app_host + ':';
const grpc_address = prefix + proto_port;

export function Run() {
    const code = Proto.GetInstance().bind(grpc_address, credential);
    if (code != 0) {
        console.log('Listenning as gRPC server on address', grpc_address);
        Proto.GetInstance().start();
    } else {
        console.error("Starting the gRPC server has failed");
    }
}