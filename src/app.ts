import * as ApiRest from './server/rest';
import * as Proto from './server/proto';
import * as CacheFactory from './cache/factory';
import * as Cache from './cache/public';
import environment from './config';

const cache_size = +environment.CacheSize;
const cache_dead = +environment.CacheDeadline * 1000; // cache deadline must be in seconds
const cache = CacheFactory.NewRandomImageCache(cache_size, cache_dead);
Cache.InitGlobalCache(cache);

const app_host = environment.AppHost;

//const api_port = environment.RestServicePort;
//ApiRest.GetInstance().listen(api_port)
//
//console.log('Listenning as API Rest on port', api_port);

const proto_port = environment.ProtoServicePort;
const credential = Proto.GetCredential();
Proto.GetInstance().bind(app_host + ':' + proto_port, credential);

console.log('Listenning as gRPC server on port', proto_port);