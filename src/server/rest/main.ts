import * as Api from '../../server/rest/express';
import environment from '../../config';

const api_port = environment.RestServicePort;

export function Run() {
    const rest_server = Api.GetExpressInstance().listen(api_port);
    //const rest_server = ApiRest.GetHttpInstance().listen(api_port);

    console.log('Listenning as API server on port', api_port);
}