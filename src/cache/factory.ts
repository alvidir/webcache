import { WallpaperCache } from './cache';
import { Cache } from './interfaces';

export function NewWallpaperCache(size: number, timeout: number): Cache {
    let instace = new WallpaperCache(size);
    instace.SetDeadline(timeout);
    return instace;
}