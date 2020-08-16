import * as UnsplashApi from '../api/unsplash';
import * as Interfaces from './interfaces';
import { ImageInfo } from '../api/interfaces';

export class RandomImageCache implements Interfaces.Cache {
    current: ImageInfo[];
    deadline?: number;
    last?: Date;

    constructor(size: number) {
        if (size < 1) {
            throw new Error("The minium size of any caché is 1 item");
        }

        this.current = new Array(size);
        // initialize cache images
        this.update();
    }

    private update() {
        UnsplashApi.GetInstance().HandleRandomRequest(this.current.length, (imgs: ImageInfo[]): void => {
            for (let index = 0; index < imgs.length && index < this.current.length; index++) {
                console.log(JSON.stringify(imgs[index]));
                this.current[index] = imgs[index];
            }
        
            this.last = new Date();
            console.log("Images cache has been updated at ", this.last.toDateString());
        });
    }

    private onHit() {
        const last = this.last ? this.last.getTime() : 0;
        const diff = new Date().getTime() - last;

        if (diff < 0 || this.deadline && diff > this.deadline) {
            this.update()
        }
    }

    private getRandomIndex(): number {
        const min = Math.ceil(0);
        const max = Math.floor(this.current.length);
        return Math.floor(Math.random() * (max - min)) + min;
    }

    private getRandomSlice(n: number): number[] {
        let result: number[] = new Array();
        while (result.length < n) {
            let cand = this.getRandomIndex();
            if (result.indexOf(cand) < 0) {
                result.push(cand);
            }
        }

        return result;
    }

    SetDeadline(deadline: number) {
        this.deadline = deadline;
    }

    // GetSingle returns a random image in the cache
    GetSingle(): ImageInfo {
        const index = this.getRandomIndex();
        const image = this.current[index];
        this.onHit();
        return image;
    }

    // GetSingle returns the item stored in the provided index
    GetByTag(tag: number): ImageInfo | undefined {
        if (tag >= this.current.length) {
            return;
        }

        const image = this.current[tag];
        this.onHit();
        return image;
    }

    // GetSubset returns a random subset of n elements as maximum
    GetSubset(n: number): ImageInfo[] {
        let image: ImageInfo[] = new Array(n);
        let indexes = this.getRandomSlice(n);
        indexes.forEach((value: number, index: number, _: number[]): void => {
            image[index] = this.current[value];
        }, this);

        this.onHit();
        return image;
    }

    // GetAllItems returns all the items stored in the cache
    GetAllItems(): ImageInfo[] {
        const image = this.current;
        this.onHit();
        return image;
    }

}