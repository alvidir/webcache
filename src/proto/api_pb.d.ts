import * as jspb from "google-protobuf"

export class Image extends jspb.Message {
  getUrlsMap(): jspb.Map<string, string>;
  clearUrlsMap(): Image;

  getAuthor(): string;
  setAuthor(value: string): Image;

  getProfile(): string;
  setProfile(value: string): Image;

  getProfileImage(): string;
  setProfileImage(value: string): Image;

  getBio(): string;
  setBio(value: string): Image;

  getLikes(): number;
  setLikes(value: number): Image;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Image.AsObject;
  static toObject(includeInstance: boolean, msg: Image): Image.AsObject;
  static serializeBinaryToWriter(message: Image, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Image;
  static deserializeBinaryFromReader(message: Image, reader: jspb.BinaryReader): Image;
}

export namespace Image {
  export type AsObject = {
    urlsMap: Array<[string, string]>,
    author: string,
    profile: string,
    profileImage: string,
    bio: string,
    likes: number,
  }
}

export class EmptyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmptyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EmptyRequest): EmptyRequest.AsObject;
  static serializeBinaryToWriter(message: EmptyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmptyRequest;
  static deserializeBinaryFromReader(message: EmptyRequest, reader: jspb.BinaryReader): EmptyRequest;
}

export namespace EmptyRequest {
  export type AsObject = {
  }
}

