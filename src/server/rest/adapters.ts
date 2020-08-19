import http from 'http';
import { Callback, ImageInfo } from '../../api/interfaces';
import { ResponseSender } from './http';

function CallbackAdapter(res: http.ServerResponse): Callback {
    return function (imgs: ImageInfo[]): void {
        let response = imgs.length > 0 ? imgs : undefined;
        ResponseSender(res, response);
    }
}