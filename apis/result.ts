import * as Interface from './interfaces';

export class Result implements Interface.Result {
    result: string = "None";
    err: string = "None";

    constructor() { }

    set(result: string, err: string): void {
        this.result = result;
        this.err = err;
    }

    setResult(result: string): void {
        this.result = result;
    }

    setError(err: string): void {
        this.err = err;
    }

    missing(name: string): void {
        this.err = 'Got no value for required parameter ' + name;
    }

    build(): Result {
        return this;
    }

    Get(): string {
        return this.result;
    }

    Err(): string {
        return this.err;
    }

    Ok(): boolean {
        return this.err !== null;
    }
}