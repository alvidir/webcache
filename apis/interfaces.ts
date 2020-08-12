import { ParsedUrlQuery } from 'querystring';

export interface Result {
    Get(): string;
    Err(): string;
    Ok(): boolean;
}

export interface UnsplashApi {
    HandleSingleRequest(): Result;
    HandleRollRequest(query: ParsedUrlQuery): Result;
    HandleSupplyRequest(query: ParsedUrlQuery): Result;
}