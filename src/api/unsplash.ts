import Unsplash from 'unsplash-js';
import * as Interface from './interfaces';
import { toJson } from 'unsplash-js';

// required to get environment configuration
import environment from '../config';

const apiKey: string = environment.ApiKey;
const timeout = + environment.ApiTimeout;
const unsplash = new Unsplash({ accessKey: apiKey, timeout: timeout });

class UnsplashApi implements Interface.UnsplashApi {

    private ParseSourceJSON(json: any): string[] {
        const size = json.length
        let stack: string[] = new Array(size);

        for (let index = 0; index < size || 0; index++) {
            stack[index] = json[index];
        }

        return stack;
    }

    private HandleError(err: Error) {
        console.log(err);
    }

    async HandleWallpaperRequest(n: number, callback: Interface.Callback) {
        unsplash.photos.getRandomPhoto({
            username: undefined,
            query: undefined,
            featured: undefined,
            collections: undefined,
            count: n
        })
            .then(toJson)
            .then(this.ParseSourceJSON)
            .then(callback)
            .catch(this.HandleError);
    }
}

const instance: Interface.UnsplashApi = new UnsplashApi();
export function GetInstance(): Interface.UnsplashApi {
    return instance;
}