import {ImageInfo} from '../../api/interfaces';
import * as api_pb from '../../proto/api_pb';

function CopyImageUrlMap(from: any, to: api_pb.Image){
    const keys = Object.keys(from);
    keys.forEach(function(value: string, index: number, array: string[]): void {
        to.getUrlsMap().set(value, from[value]);
    });
}  

export function ProtobufAdapter(img: ImageInfo): api_pb.Image {
    const urls = img.Urls();

    const likes = img.Likes();
    const bio = img.Bio();
    const prof_img = img.ProfileImage();
    
    let response = new api_pb.Image();
    response.setAuthor(img.Author());
    response.setProfile(img.Profile());

    CopyImageUrlMap(urls, response);

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