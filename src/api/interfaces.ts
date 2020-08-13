import http from 'http';
import { ParsedUrlQuery } from 'querystring';

export type Callback = (imgs: ImageInfo[]) => void;

export interface ImageInfo {
    update(author: string, urls: string): void;
    Author(): string;
    Urls(): string;
}

export interface UnsplashApi {
    HandleSingleRequest(cb: Callback): void;
    HandleRollRequest(query: ParsedUrlQuery, cb: Callback): void;
    HandleSupplyRequest(query: ParsedUrlQuery, cb: Callback): void;
}