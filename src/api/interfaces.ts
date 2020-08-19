import http from 'http';
import { ParsedUrlQuery } from 'querystring';

export type Callback = (imgs: ImageInfo[]) => void;

export interface ImageInfo {
    Update(author: string, urls: any): void;
    Author(): string;
    Urls(): any;
    Profile(): string;
    ProfileImage(): string | undefined;
    Bio(): string | undefined;
    Likes(): number | undefined;
}

export interface UnsplashApi {
    HandleRandomRequest(n: number, cb: Callback): void;
    HandleQueryRequest(query: ParsedUrlQuery, cb: Callback): void;
}