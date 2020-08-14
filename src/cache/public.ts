import * as Interfaces from './interfaces';

let GlobalCache: Interfaces.Cache;
export function InitGlobalCache(cache: Interfaces.Cache) {
    if (!GlobalCache) {
        GlobalCache = cache;
    }
}

export function GetInstance(): Interfaces.Cache {
    return GlobalCache;
}