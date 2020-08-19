import http from 'http';
import url from 'url';
import { ParsedUrlQuery } from 'querystring';
import { GetInstance as CacheInstance } from '../../cache/public';
import { ResponseSender } from './http';

export let HandleRequest = (req: http.IncomingMessage, res: http.ServerResponse) => {
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
        case '/single':
        default:
            // return a single item in the cache
            response = CacheInstance().GetSingle();
            break;
    }

    ResponseSender(res, response);
}