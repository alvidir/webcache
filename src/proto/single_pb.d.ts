// package: proto
// file: src/proto/single.proto

import * as jspb from "google-protobuf";

export class SingleRequest extends jspb.Message {
  getCode(): number;
  setCode(value: number): void;

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

