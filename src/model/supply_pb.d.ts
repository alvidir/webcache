// package: model
// file: supply.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

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
