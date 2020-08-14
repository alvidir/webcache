import * as Interface from './interfaces';

export class ImageInfo implements Interface.ImageInfo {
    urls: string;
    author: string;
    profile: string;
    profile_image?: string;
    bio?: string;
    likes?: string;

    constructor(author: string, urls: string, source: string) {
        this.author = author;
        this.urls = urls;
        this.profile = source;
    }

    Update(author: string, urls: string) {
        this.author = author;
        this.urls = urls;
    }

    Urls(): string {
        return this.urls;
    }

    Author(): string {
        return this.author;
    }
}