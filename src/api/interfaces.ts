export type Callback = (stack: string[]) => void;

export interface UnsplashApi {
    HandleWallpaperRequest(n: number, cb: Callback): void;
}