import * as Api from '../server/rest/express';
import { ApiClient } from '../proto/api_grpc_pb';
import { EmptyRequest } from '../proto/api_pb';
import * as Proto from '../server/proto/grpc';
import environment from '../config';

const proto_port = environment.ProtoServicePort;
const prefix = environment.AppHost + ':';
const grpc_address = prefix + proto_port;

const ch_credential = Proto.GetChannelCredential();
const client = new ApiClient(grpc_address, ch_credential);
const req = new EmptyRequest();

export function callback(): void {
    client.single(req, (err, response) => {
        if (err) {
            console.log("Got an error while testing gRPC server: shutting down");
            Proto.GetInstance().forceShutdown();
            Api.GetExpressInstance().close();
        } else {
            console.log('Grpc test successfully passed');
        }
    });
}