var http = require('http');
var url = require('url');

var api = require('./apis/unsplash');

const port = process.env.SERVICE_PORT || 3001;

let HandleRequest = (req, res) => {
    let result = 'None';

    let path = url.parse(req.url, true).pathname;
    let qu = url.parse(req.url, true).query;

    switch (path) {
        case '/roll':
            result = api.HandleRollRequest(path);
            break;

        case '/supply':
            result = api.HandleSupplyRequest(path);
            break;

        case '/single':
        default:
            result = api.HandleSingleRequest();
    }

    res.setHeader('Content-Type', 'application/json');
    res.statusCode = 200;
    res.write('{"result": ' + result + '}');

    res.end(); //end the response
}

//create a server object
const server = http.createServer(HandleRequest);
server.listen(port);

console.log('Listenning on port', port);