import { GetInstance as CacheInstance } from '../../cache/public';

const express = require('express')
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

export const WallpaperPath = '/wallpaper'

// get call has cors enabled
app.get(WallpaperPath, cors(), (_: any, res: any) => {
    const response = CacheInstance().GetWallpaper();
    res.send(response);
});

export function GetExpressInstance(): any{
    return app;
}