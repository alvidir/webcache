import * as Interfaces from '../api/interfaces';

export interface Cache {
    GetSingle(): Interfaces.ImageInfo;
    GetByTag(tag: number): Interfaces.ImageInfo | undefined;
    GetSubset(n: number): Interfaces.ImageInfo[];
    GetAllItems(): Interfaces.ImageInfo[];
}