import http from 'http';
import url from 'url';
import { ParsedUrlQuery } from 'querystring';
import { GetInstance as CacheInstance } from '../cache/public';
import { GetInstance as UnsplashInstance } from '../api/unsplash';
import { Callback, ImageInfo } from '../api/interfaces';

function ResponseSender(res: http.ServerResponse, response: any | undefined) {
    res.statusCode = response ? 200 : 400;
    if (!response){
        response = '';
    }

    res.write(JSON.stringify(response))
    res.end();
}

function CallbackAdapter(res: http.ServerResponse): Callback {
    return function (imgs: ImageInfo[]): void {
        let response = imgs.length > 0 ? imgs : undefined;
        ResponseSender(res, response);
    }
}

let HandleRequest = (req: http.IncomingMessage, res: http.ServerResponse) => {
    res.setHeader('Content-Type', 'application/json');
    let source: string = req.url ? req.url : 'None';

    if (!source) {
        res.statusCode = 400;
        res.end();
        return;
    }

    let path = url.parse(source, true).pathname;
    let query: ParsedUrlQuery = url.parse(source, true).query;
    let response: any;

    switch (path) {
        case '/roll':
            // return all items in the cache
            response = CacheInstance().GetAllItems()
            break;
        case '/single':
        default:
            // return a single item in the cache
            response = CacheInstance().GetSingle();
            break;
    }

    ResponseSender(res, response);
}

//create a server object
const server = http.createServer(HandleRequest);
export function GetInstance(): http.Server {
    return server;
}