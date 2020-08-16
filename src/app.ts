import * as ApiRest from './server/rest';
import * as Proto from './server/proto';
import * as CacheFactory from './cache/factory';
import * as Cache from './cache/public';
import environment from './config';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

const cache_size = +environment.CacheSize;
const cache_dead = +environment.CacheDeadline * 1000; // cache deadline must be in seconds
const cache = CacheFactory.NewRandomImageCache(cache_size, cache_dead);
Cache.InitGlobalCache(cache);

const app_host = environment.AppHost;

// gRPC server
const proto_port = environment.ProtoServicePort;
const credential = Proto.GetCredential();
const address = app_host + ':' + proto_port;
const code = Proto.GetInstance().bind(address, credential);
if (code != 0) {
    console.log('Listenning as gRPC server on address', address);
    Proto.GetInstance().start();
} else {
    console.error("Starting the gRPC server has failed");
}

// Rest server
// const api_port = environment.RestServicePort;
// ApiRest.GetInstance().listen(api_port)
// 
// console.log('Listenning as API Rest on port', api_port);
