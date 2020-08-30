import * as UnsplashApi from '../api/unsplash';
import * as Interfaces from './interfaces';

export class WallpaperCache implements Interfaces.Cache {
    current: string[];
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
        UnsplashApi.GetInstance().HandleWallpaperRequest(this.current.length, (stack: string[]): void => {
            for (let index = 0; index < stack.length && index < this.current.length; index++) {
                this.current[index] = stack[index];
            }
        
            this.last = new Date();
            console.log("Images cache has been updated at ", this.last.toDateString());
        });
    }

    private onHit() {
        const last = this.last ? this.last.getTime() : 0;
        const now = new Date().getTime();
        const diff = now - last;
        
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
    GetWallpaper(): string {
        const index = this.getRandomIndex();
        const image = this.current[index];
        this.onHit();
        return image;
    }
}