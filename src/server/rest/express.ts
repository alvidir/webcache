import { GetInstance as CacheInstance } from '../../cache/public';

const express = require('express')
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

export const SingleMethodPath = '/'

// get call has cors enabled
app.get(SingleMethodPath, cors(), (_: any, res: any) => {
    const response = CacheInstance().GetSingle();
    res.send(response);
});

export function GetExpressInstance(): any{
    return app;
}