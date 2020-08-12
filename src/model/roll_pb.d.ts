// package: model
// file: roll.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class RollRequest extends jspb.Message { 
    getCode(): number;
    setCode(value: number): RollRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RollRequest.AsObject;
    static toObject(includeInstance: boolean, msg: RollRequest): RollRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RollRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RollRequest;
    static deserializeBinaryFromReader(message: RollRequest, reader: jspb.BinaryReader): RollRequest;
}

export namespace RollRequest {
    export type AsObject = {
        code: number,
    }
}
