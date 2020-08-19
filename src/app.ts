import * as ApiRest from './server/rest';
import * as Proto from './server/proto';
import * as CacheFactory from './cache/factory';
import * as Cache from './cache/public';
import environment from './config';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

const cache_size = +environment.CacheSize;
const cache_dead = +environment.CacheDeadline * 60000; // cache deadline must be in seconds
const cache = CacheFactory.NewRandomImageCache(cache_size, cache_dead);
Cache.InitGlobalCache(cache);

const app_host = environment.AppHost;

// gRPC server
const proto_port = environment.ProtoServicePort;
const credential = Proto.GetServerCredential();
const prefix = app_host + ':';
const grpc_address = prefix + proto_port;
const code = Proto.GetInstance().bind(grpc_address, credential);
if (code != 0) {
    console.log('Listenning as gRPC server on address', grpc_address);
    Proto.GetInstance().start();
} else {
    console.error("Starting the gRPC server has failed");
}

// Rest server
const api_port = environment.RestServicePort;
const api_address = api_port;
const rest_server = ApiRest.GetInstance().listen(api_address)

console.log('Listenning as API server on address', prefix + api_address);

// gRPC testing
import { ApiClient } from './proto/api_grpc_pb';
import { EmptyRequest } from './proto/api_pb';
const ch_credential = Proto.GetChannelCredential();
const client = new ApiClient(grpc_address, ch_credential);
const req = new EmptyRequest();

function callback(): void {
    client.single(req, (err, _) => {
        if (err) {
            console.log("Got an error while testing gRPC server: shutting down");
            Proto.GetInstance().forceShutdown();
            rest_server.close();
        }
    });
}

setTimeout(callback, 1500);