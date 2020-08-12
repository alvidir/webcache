import * as Interface from './interfaces';
import { Result } from './result';
import Unsplash from 'unsplash-js';
import { ParsedUrlQuery } from 'querystring';

import url from 'url';
const apiKey: string = process.env.API_KEY || "None";
const unsplash = new Unsplash({ accessKey: apiKey });

class UnsplashApi implements Interface.UnsplashApi {

    HandleSingleRequest(): Interface.Result {
        let result = new Result();
        result.setResult('single');
        return result.build();
    }

    HandleRollRequest(query: ParsedUrlQuery): Interface.Result {
        let result = new Result();
        result.setResult("None");

        if (!query['hola']) {
            result.missing('hola');
        }

        return result.build();
    }

    HandleSupplyRequest(query: ParsedUrlQuery): Interface.Result {
        let result = new Result();
        result.setResult("None");

        if (query.name?.length == 0) {
            result.setError('No parameters set for supply call');
        }

        return result.build();
    }
}

const instance: Interface.UnsplashApi = new UnsplashApi();
export function GetInstance(): Interface.UnsplashApi {
    return instance;
}