import * as CacheFactory from './cache/factory';
import * as Cache from './cache/public';
import * as grpc_srv from './server/proto/main';
import * as http_srv from './server/rest/main';
import * as grpc_tst from './test/grpc';
import environment from './config';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

// seting up cache
const cache_size = +environment.CacheSize;
const cache_dead = +environment.CacheDeadline * 1000; // cache deadline must be in seconds
const cache = CacheFactory.NewRandomImageCache(cache_size, cache_dead);
Cache.InitGlobalCache(cache);

// gRPC server
grpc_srv.Run();
// Rest server
http_srv.Run();

setTimeout(grpc_tst.callback, 1500);