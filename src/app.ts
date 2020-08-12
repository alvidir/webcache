import * as ApiRest from './server/rest'
import * as Proto from './server/proto'
import environment from './config';

const app_host = environment.AppHost;
const api_port = environment.RestServicePort;
ApiRest.GetInstance().listen(api_port)

console.log('Listenning as API Rest on port', api_port);

const proto_port = environment.ProtoServicePort;
const credential = Proto.GetCredential();
Proto.GetInstance().bind(app_host + ':' + proto_port, credential);

console.log('Listenning as gRPC server on port', proto_port);