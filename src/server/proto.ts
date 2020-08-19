import { Server, ServerUnaryCall, sendUnaryData, ServerCredentials, credentials, ChannelCredentials } from 'grpc';
import { IApiServer, ApiService } from '../proto/api_grpc_pb';
import * as pb from '../proto/api_pb';
import { GetInstance as CacheInstance } from '../cache/public';
import {ImageInfo} from '../api/interfaces';

function CopyMap(from: any, to: any){
    const keys = Object.keys(from);
    keys.forEach(function(value: string, index: number, array: string[]): void {
        to[value] = from[value];
    });
}  

function ProtobufAdapter(img: ImageInfo): pb.Image {
    const urls = img.Urls();

    const likes = img.Likes();
    const bio = img.Bio();
    const prof_img = img.ProfileImage();
    
    let response = new pb.Image();
    response.setAuthor(img.Author());
    response.setProfile(img.Profile());

    let protomap = response.getUrlsMap();
    CopyMap(urls, protomap);

    if (likes){
        response.setLikes(likes);
    }
    if (bio){
        response.setBio(bio);
    }
    if (prof_img){
        response.setProfileImage(prof_img);
    }

    return response;
}

class ApiImplementation implements IApiServer {
    public single(_: ServerUnaryCall<pb.EmptyRequest>, callback: sendUnaryData<pb.Image>): void {
        const image = CacheInstance().GetSingle();
        const value = ProtobufAdapter(image);
        callback(null, value, undefined, undefined);
    }
}

//create a server object
const server = new Server();
server.addService<IApiServer>(ApiService, new ApiImplementation());

export function GetInstance(): Server {
    return server;
}

const srv_credential = ServerCredentials.createInsecure();
export function GetServerCredential(): ServerCredentials {
    return srv_credential;
}

const ch_credential = credentials.createInsecure();
export function GetChannelCredential(): ChannelCredentials {
    return ch_credential;
}