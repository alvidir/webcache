import { RandomImageCache } from './cache';
import { Cache } from './interfaces';

export function NewRandomImageCache(size: number, timeout: number): Cache {
    let instace = new RandomImageCache(size);
    instace.SetDeadline(timeout);
    return instace;
}