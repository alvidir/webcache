import http from 'http';
import url from 'url';

import * as UnsplashApi from '../api/unsplash';
import * as Interfaces from '../api/interfaces';
import { ParsedUrlQuery } from 'querystring';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

function CallbackAdapter(res: http.ServerResponse): Interfaces.Callback {
    return function (imgs: Interfaces.ImageInfo[]): void {
        res.statusCode = imgs.length > 0 ? 200 : 400;
        res.write(JSON.stringify(imgs));
        res.end();
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
    const callback = CallbackAdapter(res);

    switch (path) {
        case '/roll':
            UnsplashApi.GetInstance().HandleRollRequest(query, callback);
            break;

        case '/supply':
            UnsplashApi.GetInstance().HandleSupplyRequest(query, callback);
            break;

        case '/single':
        default:
            UnsplashApi.GetInstance().HandleSingleRequest(callback);
    }
}

//create a server object
const server = http.createServer(HandleRequest);
export function GetInstance(): http.Server {
    return server;
}