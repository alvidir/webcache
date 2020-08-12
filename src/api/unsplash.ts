import * as Interface from './interfaces';
import { Result } from './result';
import Unsplash from 'unsplash-js';
import { toJson } from 'unsplash-js';
import { ParsedUrlQuery } from 'querystring';

// required to get environment configuration
import environment from '../config';

const apiKey: string = environment.ApiKey;
const unsplash = new Unsplash({ accessKey: apiKey, timeout: 500 });

class UnsplashApi implements Interface.UnsplashApi {

    HandleSingleRequest(): Interface.Result {
        let result = new Result();
        unsplash.photos.getRandomPhoto({
            username: "naoufal",
            query: undefined,
            featured: undefined,
            collections: undefined,
            count: 1
        })
            .then(toJson)
            .then(json => {
                console.log(json);
            });

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