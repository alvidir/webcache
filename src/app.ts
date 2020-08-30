import * as CacheFactory from './cache/factory';
import * as Cache from './cache/public';
import * as http_srv from './server/rest/main';
import environment from './config';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

// seting up cache
const cache_size = +environment.CacheSize;
const cache_dead = +environment.CacheDeadline * 1000; // cache deadline must be in seconds
const cache = CacheFactory.NewWallpaperCache(cache_size, cache_dead);
Cache.InitGlobalCache(cache);

// Rest server
http_srv.Run();