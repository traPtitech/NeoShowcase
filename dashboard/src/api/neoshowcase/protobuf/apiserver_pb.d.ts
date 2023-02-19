import * as jspb from 'google-protobuf'

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';


export class App extends jspb.Message {
  getId(): string;
  setId(value: string): App;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): App.AsObject;
  static toObject(includeInstance: boolean, msg: App): App.AsObject;
  static serializeBinaryToWriter(message: App, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): App;
  static deserializeBinaryFromReader(message: App, reader: jspb.BinaryReader): App;
}

export namespace App {
  export type AsObject = {
    id: string,
  }
}

export class GetAppsResponse extends jspb.Message {
  getAppsList(): Array<App>;
  setAppsList(value: Array<App>): GetAppsResponse;
  clearAppsList(): GetAppsResponse;
  addApps(value?: App, index?: number): App;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAppsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAppsResponse): GetAppsResponse.AsObject;
  static serializeBinaryToWriter(message: GetAppsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAppsResponse;
  static deserializeBinaryFromReader(message: GetAppsResponse, reader: jspb.BinaryReader): GetAppsResponse;
}

export namespace GetAppsResponse {
  export type AsObject = {
    appsList: Array<App.AsObject>,
  }
}

