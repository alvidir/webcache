import http from 'http';
import Unsplash from 'unsplash-js';
import * as Interface from './interfaces';
import { ImageInfo } from './image';
import { toJson } from 'unsplash-js';
import { ParsedUrlQuery } from 'querystring';

// required to get environment configuration
import environment from '../config';

const apiKey: string = environment.ApiKey;
const timeout = + environment.ApiTimeout;
const unsplash = new Unsplash({ accessKey: apiKey, timeout: timeout });

class UnsplashApi implements Interface.UnsplashApi {

    private ParseSourceJSON(json: any): Interface.ImageInfo[] {
        let stack: Interface.ImageInfo[] = new Array();
        for (let index = 0; index < json.length || 0; index++) {
            const data = json[index];
            const author = data['user']['username'];
            const source = data['user']['links']['html'];
            const urls = data['urls'];

            let image = new ImageInfo(author, urls, source);
            image.profile_image = data['user']['profile_image']['medium'];
            image.bio = data['user']['bio'];
            image.likes = data['likes'];

            stack.push(image);
        }

        return stack;
    }

    private HandleError(err: Error) {
        console.log(err);
    }

    async HandleRandomRequest(n: number, callback: Interface.Callback) {
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

    async HandleQueryRequest(query: ParsedUrlQuery, callback: Interface.Callback) {
    }
}

const instance: Interface.UnsplashApi = new UnsplashApi();
export function GetInstance(): Interface.UnsplashApi {
    return instance;
}