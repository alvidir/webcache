// package: proto
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class Result extends jspb.Message { 
    getCode(): number;
    setCode(value: number): Result;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Result.AsObject;
    static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Result;
    static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
}

export namespace Result {
    export type AsObject = {
        code: number,
    }
}

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

export class SingleRequest extends jspb.Message { 
    getCode(): number;
    setCode(value: number): SingleRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SingleRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SingleRequest): SingleRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SingleRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SingleRequest;
    static deserializeBinaryFromReader(message: SingleRequest, reader: jspb.BinaryReader): SingleRequest;
}

export namespace SingleRequest {
    export type AsObject = {
        code: number,
    }
}

export class SupplyRequest extends jspb.Message { 
    getCode(): number;
    setCode(value: number): SupplyRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SupplyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SupplyRequest): SupplyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SupplyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SupplyRequest;
    static deserializeBinaryFromReader(message: SupplyRequest, reader: jspb.BinaryReader): SupplyRequest;
}

export namespace SupplyRequest {
    export type AsObject = {
        code: number,
    }
}
