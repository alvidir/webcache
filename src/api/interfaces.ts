import http from 'http';
import { ParsedUrlQuery } from 'querystring';

export type Callback = (imgs: ImageInfo[]) => void;

export interface ImageInfo {
    Update(author: string, urls: string): void;
    Author(): string;
    Urls(): string;
}

export interface UnsplashApi {
    HandleRandomRequest(n: number, cb: Callback): void;
    HandleQueryRequest(query: ParsedUrlQuery, cb: Callback): void;
}