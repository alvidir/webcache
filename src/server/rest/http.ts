import http from 'http';
import { HandleRequest } from './handlers';

export function ResponseSender(res: http.ServerResponse, response: any | undefined) {
    res.statusCode = response ? 200 : 400;
    if (!response){
        response = '';
    }

    const format = JSON.stringify(response, null, 4);
    console.log(format);
    res.end(format);
}

// create a http server object
const server = http.createServer(HandleRequest);
export function GetHttpInstance(): http.Server {
    return server;
}