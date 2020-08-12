import http from 'http';
import url from 'url';


import * as UnsplashApi from './api/unsplash';
import * as Interfaces from './api/interfaces';

// fetch is required by the unsplash api
global.fetch = require('node-fetch');

import { ParsedUrlQuery } from 'querystring';
const port = process.env.SERVICE_PORT || 3001;

let HandleRequest = (req: http.IncomingMessage, res: http.ServerResponse) => {
    let result: Interfaces.Result;

    let source: string = req.url ? req.url : 'None';
    if (!source) {
        res.statusCode = 400;
        res.end();
        return;
    }

    let path = url.parse(source, true).pathname;
    let query: ParsedUrlQuery = url.parse(source, true).query;

    switch (path) {
        case '/roll':
            result = UnsplashApi.GetInstance().HandleRollRequest(query);
            break;

        case '/supply':
            result = UnsplashApi.GetInstance().HandleSupplyRequest(query);
            break;

        case '/single':
        default:
            result = UnsplashApi.GetInstance().HandleSingleRequest();
    }

    res.setHeader('Content-Type', 'application/json');
    res.statusCode = result.Ok() ? 200 : 400;
    res.write('{"result": ' + JSON.stringify(result) + '}');

    res.end(); //end the response
}

//create a server object
const server = http.createServer(HandleRequest);
server.listen(port);

console.log('Listenning on port', port);