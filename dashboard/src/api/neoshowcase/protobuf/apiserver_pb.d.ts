import * as jspb from 'google-protobuf'

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as neoshowcase_protobuf_null_pb from '../../neoshowcase/protobuf/null_pb';


export class Repository extends jspb.Message {
  getId(): string;
  setId(value: string): Repository;

  getName(): string;
  setName(value: string): Repository;

  getUrl(): string;
  setUrl(value: string): Repository;

  getAuthMethod(): string;
  setAuthMethod(value: string): Repository;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Repository.AsObject;
  static toObject(includeInstance: boolean, msg: Repository): Repository.AsObject;
  static serializeBinaryToWriter(message: Repository, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Repository;
  static deserializeBinaryFromReader(message: Repository, reader: jspb.BinaryReader): Repository;
}

export namespace Repository {
  export type AsObject = {
    id: string,
    name: string,
    url: string,
    authMethod: string,
  }
}

export class CreateRepositoryAuthBasic extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): CreateRepositoryAuthBasic;

  getPassword(): string;
  setPassword(value: string): CreateRepositoryAuthBasic;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRepositoryAuthBasic.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRepositoryAuthBasic): CreateRepositoryAuthBasic.AsObject;
  static serializeBinaryToWriter(message: CreateRepositoryAuthBasic, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRepositoryAuthBasic;
  static deserializeBinaryFromReader(message: CreateRepositoryAuthBasic, reader: jspb.BinaryReader): CreateRepositoryAuthBasic;
}

export namespace CreateRepositoryAuthBasic {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class CreateRepositoryAuthSSH extends jspb.Message {
  getSshKey(): string;
  setSshKey(value: string): CreateRepositoryAuthSSH;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRepositoryAuthSSH.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRepositoryAuthSSH): CreateRepositoryAuthSSH.AsObject;
  static serializeBinaryToWriter(message: CreateRepositoryAuthSSH, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRepositoryAuthSSH;
  static deserializeBinaryFromReader(message: CreateRepositoryAuthSSH, reader: jspb.BinaryReader): CreateRepositoryAuthSSH;
}

export namespace CreateRepositoryAuthSSH {
  export type AsObject = {
    sshKey: string,
  }
}

export class CreateRepositoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateRepositoryRequest;

  getUrl(): string;
  setUrl(value: string): CreateRepositoryRequest;

  getNone(): google_protobuf_empty_pb.Empty | undefined;
  setNone(value?: google_protobuf_empty_pb.Empty): CreateRepositoryRequest;
  hasNone(): boolean;
  clearNone(): CreateRepositoryRequest;

  getBasic(): CreateRepositoryAuthBasic | undefined;
  setBasic(value?: CreateRepositoryAuthBasic): CreateRepositoryRequest;
  hasBasic(): boolean;
  clearBasic(): CreateRepositoryRequest;

  getSsh(): CreateRepositoryAuthSSH | undefined;
  setSsh(value?: CreateRepositoryAuthSSH): CreateRepositoryRequest;
  hasSsh(): boolean;
  clearSsh(): CreateRepositoryRequest;

  getAuthCase(): CreateRepositoryRequest.AuthCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRepositoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRepositoryRequest): CreateRepositoryRequest.AsObject;
  static serializeBinaryToWriter(message: CreateRepositoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRepositoryRequest;
  static deserializeBinaryFromReader(message: CreateRepositoryRequest, reader: jspb.BinaryReader): CreateRepositoryRequest;
}

export namespace CreateRepositoryRequest {
  export type AsObject = {
    name: string,
    url: string,
    none?: google_protobuf_empty_pb.Empty.AsObject,
    basic?: CreateRepositoryAuthBasic.AsObject,
    ssh?: CreateRepositoryAuthSSH.AsObject,
  }

  export enum AuthCase { 
    AUTH_NOT_SET = 0,
    NONE = 3,
    BASIC = 4,
    SSH = 5,
  }
}

export class ApplicationConfig extends jspb.Message {
  getUseMariadb(): boolean;
  setUseMariadb(value: boolean): ApplicationConfig;

  getUseMongodb(): boolean;
  setUseMongodb(value: boolean): ApplicationConfig;

  getBaseImage(): string;
  setBaseImage(value: string): ApplicationConfig;

  getDockerfileName(): string;
  setDockerfileName(value: string): ApplicationConfig;

  getArtifactPath(): string;
  setArtifactPath(value: string): ApplicationConfig;

  getBuildCmd(): string;
  setBuildCmd(value: string): ApplicationConfig;

  getEntrypointCmd(): string;
  setEntrypointCmd(value: string): ApplicationConfig;

  getAuthentication(): AuthenticationType;
  setAuthentication(value: AuthenticationType): ApplicationConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationConfig): ApplicationConfig.AsObject;
  static serializeBinaryToWriter(message: ApplicationConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationConfig;
  static deserializeBinaryFromReader(message: ApplicationConfig, reader: jspb.BinaryReader): ApplicationConfig;
}

export namespace ApplicationConfig {
  export type AsObject = {
    useMariadb: boolean,
    useMongodb: boolean,
    baseImage: string,
    dockerfileName: string,
    artifactPath: string,
    buildCmd: string,
    entrypointCmd: string,
    authentication: AuthenticationType,
  }
}

export class Website extends jspb.Message {
  getId(): string;
  setId(value: string): Website;

  getFqdn(): string;
  setFqdn(value: string): Website;

  getPathPrefix(): string;
  setPathPrefix(value: string): Website;

  getStripPrefix(): boolean;
  setStripPrefix(value: boolean): Website;

  getHttps(): boolean;
  setHttps(value: boolean): Website;

  getHttpPort(): number;
  setHttpPort(value: number): Website;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Website.AsObject;
  static toObject(includeInstance: boolean, msg: Website): Website.AsObject;
  static serializeBinaryToWriter(message: Website, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Website;
  static deserializeBinaryFromReader(message: Website, reader: jspb.BinaryReader): Website;
}

export namespace Website {
  export type AsObject = {
    id: string,
    fqdn: string,
    pathPrefix: string,
    stripPrefix: boolean,
    https: boolean,
    httpPort: number,
  }
}

export class CreateWebsiteRequest extends jspb.Message {
  getFqdn(): string;
  setFqdn(value: string): CreateWebsiteRequest;

  getPathPrefix(): string;
  setPathPrefix(value: string): CreateWebsiteRequest;

  getStripPrefix(): boolean;
  setStripPrefix(value: boolean): CreateWebsiteRequest;

  getHttps(): boolean;
  setHttps(value: boolean): CreateWebsiteRequest;

  getHttpPort(): number;
  setHttpPort(value: number): CreateWebsiteRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebsiteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebsiteRequest): CreateWebsiteRequest.AsObject;
  static serializeBinaryToWriter(message: CreateWebsiteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebsiteRequest;
  static deserializeBinaryFromReader(message: CreateWebsiteRequest, reader: jspb.BinaryReader): CreateWebsiteRequest;
}

export namespace CreateWebsiteRequest {
  export type AsObject = {
    fqdn: string,
    pathPrefix: string,
    stripPrefix: boolean,
    https: boolean,
    httpPort: number,
  }
}

export class Application extends jspb.Message {
  getId(): string;
  setId(value: string): Application;

  getName(): string;
  setName(value: string): Application;

  getRepositoryId(): string;
  setRepositoryId(value: string): Application;

  getBranchName(): string;
  setBranchName(value: string): Application;

  getBuildType(): BuildType;
  setBuildType(value: BuildType): Application;

  getState(): ApplicationState;
  setState(value: ApplicationState): Application;

  getCurrentCommit(): string;
  setCurrentCommit(value: string): Application;

  getWantCommit(): string;
  setWantCommit(value: string): Application;

  getConfig(): ApplicationConfig | undefined;
  setConfig(value?: ApplicationConfig): Application;
  hasConfig(): boolean;
  clearConfig(): Application;

  getWebsitesList(): Array<Website>;
  setWebsitesList(value: Array<Website>): Application;
  clearWebsitesList(): Application;
  addWebsites(value?: Website, index?: number): Website;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): Application;
  clearOwnerIdsList(): Application;
  addOwnerIds(value: string, index?: number): Application;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Application.AsObject;
  static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
  static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Application;
  static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
}

export namespace Application {
  export type AsObject = {
    id: string,
    name: string,
    repositoryId: string,
    branchName: string,
    buildType: BuildType,
    state: ApplicationState,
    currentCommit: string,
    wantCommit: string,
    config?: ApplicationConfig.AsObject,
    websitesList: Array<Website.AsObject>,
    ownerIdsList: Array<string>,
  }
}

export class AvailableDomains extends jspb.Message {
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): AvailableDomains;
  clearDomainsList(): AvailableDomains;
  addDomains(value: string, index?: number): AvailableDomains;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AvailableDomains.AsObject;
  static toObject(includeInstance: boolean, msg: AvailableDomains): AvailableDomains.AsObject;
  static serializeBinaryToWriter(message: AvailableDomains, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AvailableDomains;
  static deserializeBinaryFromReader(message: AvailableDomains, reader: jspb.BinaryReader): AvailableDomains;
}

export namespace AvailableDomains {
  export type AsObject = {
    domainsList: Array<string>,
  }
}

export class AddAvailableDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): AddAvailableDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAvailableDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAvailableDomainRequest): AddAvailableDomainRequest.AsObject;
  static serializeBinaryToWriter(message: AddAvailableDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAvailableDomainRequest;
  static deserializeBinaryFromReader(message: AddAvailableDomainRequest, reader: jspb.BinaryReader): AddAvailableDomainRequest;
}

export namespace AddAvailableDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class ApplicationEnvironmentVariable extends jspb.Message {
  getKey(): string;
  setKey(value: string): ApplicationEnvironmentVariable;

  getValue(): string;
  setValue(value: string): ApplicationEnvironmentVariable;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationEnvironmentVariable.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationEnvironmentVariable): ApplicationEnvironmentVariable.AsObject;
  static serializeBinaryToWriter(message: ApplicationEnvironmentVariable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationEnvironmentVariable;
  static deserializeBinaryFromReader(message: ApplicationEnvironmentVariable, reader: jspb.BinaryReader): ApplicationEnvironmentVariable;
}

export namespace ApplicationEnvironmentVariable {
  export type AsObject = {
    key: string,
    value: string,
  }
}

export class ApplicationEnvironmentVariables extends jspb.Message {
  getVariablesList(): Array<ApplicationEnvironmentVariable>;
  setVariablesList(value: Array<ApplicationEnvironmentVariable>): ApplicationEnvironmentVariables;
  clearVariablesList(): ApplicationEnvironmentVariables;
  addVariables(value?: ApplicationEnvironmentVariable, index?: number): ApplicationEnvironmentVariable;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationEnvironmentVariables.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationEnvironmentVariables): ApplicationEnvironmentVariables.AsObject;
  static serializeBinaryToWriter(message: ApplicationEnvironmentVariables, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationEnvironmentVariables;
  static deserializeBinaryFromReader(message: ApplicationEnvironmentVariables, reader: jspb.BinaryReader): ApplicationEnvironmentVariables;
}

export namespace ApplicationEnvironmentVariables {
  export type AsObject = {
    variablesList: Array<ApplicationEnvironmentVariable.AsObject>,
  }
}

export class MariaDbKey extends jspb.Message {
  getHost(): string;
  setHost(value: string): MariaDbKey;

  getDatabase(): string;
  setDatabase(value: string): MariaDbKey;

  getUser(): string;
  setUser(value: string): MariaDbKey;

  getPassword(): string;
  setPassword(value: string): MariaDbKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MariaDbKey.AsObject;
  static toObject(includeInstance: boolean, msg: MariaDbKey): MariaDbKey.AsObject;
  static serializeBinaryToWriter(message: MariaDbKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MariaDbKey;
  static deserializeBinaryFromReader(message: MariaDbKey, reader: jspb.BinaryReader): MariaDbKey;
}

export namespace MariaDbKey {
  export type AsObject = {
    host: string,
    database: string,
    user: string,
    password: string,
  }
}

export class MongoKey extends jspb.Message {
  getHost(): string;
  setHost(value: string): MongoKey;

  getDatabase(): string;
  setDatabase(value: string): MongoKey;

  getUser(): string;
  setUser(value: string): MongoKey;

  getPassword(): string;
  setPassword(value: string): MongoKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MongoKey.AsObject;
  static toObject(includeInstance: boolean, msg: MongoKey): MongoKey.AsObject;
  static serializeBinaryToWriter(message: MongoKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MongoKey;
  static deserializeBinaryFromReader(message: MongoKey, reader: jspb.BinaryReader): MongoKey;
}

export namespace MongoKey {
  export type AsObject = {
    host: string,
    database: string,
    user: string,
    password: string,
  }
}

export class ApplicationKeys extends jspb.Message {
  getMariadbkey(): MariaDbKey | undefined;
  setMariadbkey(value?: MariaDbKey): ApplicationKeys;
  hasMariadbkey(): boolean;
  clearMariadbkey(): ApplicationKeys;

  getMongokey(): MongoKey | undefined;
  setMongokey(value?: MongoKey): ApplicationKeys;
  hasMongokey(): boolean;
  clearMongokey(): ApplicationKeys;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationKeys.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationKeys): ApplicationKeys.AsObject;
  static serializeBinaryToWriter(message: ApplicationKeys, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationKeys;
  static deserializeBinaryFromReader(message: ApplicationKeys, reader: jspb.BinaryReader): ApplicationKeys;
}

export namespace ApplicationKeys {
  export type AsObject = {
    mariadbkey?: MariaDbKey.AsObject,
    mongokey?: MongoKey.AsObject,
  }
}

export class ApplicationBuildArtifact extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): ApplicationBuildArtifact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationBuildArtifact.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationBuildArtifact): ApplicationBuildArtifact.AsObject;
  static serializeBinaryToWriter(message: ApplicationBuildArtifact, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationBuildArtifact;
  static deserializeBinaryFromReader(message: ApplicationBuildArtifact, reader: jspb.BinaryReader): ApplicationBuildArtifact;
}

export namespace ApplicationBuildArtifact {
  export type AsObject = {
    url: string,
  }
}

export class ApplicationOutput extends jspb.Message {
  getOutput(): string;
  setOutput(value: string): ApplicationOutput;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationOutput.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationOutput): ApplicationOutput.AsObject;
  static serializeBinaryToWriter(message: ApplicationOutput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationOutput;
  static deserializeBinaryFromReader(message: ApplicationOutput, reader: jspb.BinaryReader): ApplicationOutput;
}

export namespace ApplicationOutput {
  export type AsObject = {
    output: string,
  }
}

export class Build extends jspb.Message {
  getId(): string;
  setId(value: string): Build;

  getCommit(): string;
  setCommit(value: string): Build;

  getStatus(): Build.BuildStatus;
  setStatus(value: Build.BuildStatus): Build;

  getStartedAt(): neoshowcase_protobuf_null_pb.NullTimestamp | undefined;
  setStartedAt(value?: neoshowcase_protobuf_null_pb.NullTimestamp): Build;
  hasStartedAt(): boolean;
  clearStartedAt(): Build;

  getUpdatedAt(): neoshowcase_protobuf_null_pb.NullTimestamp | undefined;
  setUpdatedAt(value?: neoshowcase_protobuf_null_pb.NullTimestamp): Build;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Build;

  getFinishedAt(): neoshowcase_protobuf_null_pb.NullTimestamp | undefined;
  setFinishedAt(value?: neoshowcase_protobuf_null_pb.NullTimestamp): Build;
  hasFinishedAt(): boolean;
  clearFinishedAt(): Build;

  getRetriable(): boolean;
  setRetriable(value: boolean): Build;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Build.AsObject;
  static toObject(includeInstance: boolean, msg: Build): Build.AsObject;
  static serializeBinaryToWriter(message: Build, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Build;
  static deserializeBinaryFromReader(message: Build, reader: jspb.BinaryReader): Build;
}

export namespace Build {
  export type AsObject = {
    id: string,
    commit: string,
    status: Build.BuildStatus,
    startedAt?: neoshowcase_protobuf_null_pb.NullTimestamp.AsObject,
    updatedAt?: neoshowcase_protobuf_null_pb.NullTimestamp.AsObject,
    finishedAt?: neoshowcase_protobuf_null_pb.NullTimestamp.AsObject,
    retriable: boolean,
  }

  export enum BuildStatus { 
    BUILDING = 0,
    SUCCEEDED = 1,
    FAILED = 2,
    CANCELLED = 3,
    QUEUED = 4,
    SKIPPED = 5,
  }
}

export class BuildLog extends jspb.Message {
  getOutput(): string;
  setOutput(value: string): BuildLog;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BuildLog.AsObject;
  static toObject(includeInstance: boolean, msg: BuildLog): BuildLog.AsObject;
  static serializeBinaryToWriter(message: BuildLog, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BuildLog;
  static deserializeBinaryFromReader(message: BuildLog, reader: jspb.BinaryReader): BuildLog;
}

export namespace BuildLog {
  export type AsObject = {
    output: string,
  }
}

export class GetRepositoriesResponse extends jspb.Message {
  getRepositoriesList(): Array<Repository>;
  setRepositoriesList(value: Array<Repository>): GetRepositoriesResponse;
  clearRepositoriesList(): GetRepositoriesResponse;
  addRepositories(value?: Repository, index?: number): Repository;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRepositoriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRepositoriesResponse): GetRepositoriesResponse.AsObject;
  static serializeBinaryToWriter(message: GetRepositoriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRepositoriesResponse;
  static deserializeBinaryFromReader(message: GetRepositoriesResponse, reader: jspb.BinaryReader): GetRepositoriesResponse;
}

export namespace GetRepositoriesResponse {
  export type AsObject = {
    repositoriesList: Array<Repository.AsObject>,
  }
}

export class GetApplicationsResponse extends jspb.Message {
  getApplicationsList(): Array<Application>;
  setApplicationsList(value: Array<Application>): GetApplicationsResponse;
  clearApplicationsList(): GetApplicationsResponse;
  addApplications(value?: Application, index?: number): Application;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationsResponse): GetApplicationsResponse.AsObject;
  static serializeBinaryToWriter(message: GetApplicationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationsResponse;
  static deserializeBinaryFromReader(message: GetApplicationsResponse, reader: jspb.BinaryReader): GetApplicationsResponse;
}

export namespace GetApplicationsResponse {
  export type AsObject = {
    applicationsList: Array<Application.AsObject>,
  }
}

export class CreateApplicationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateApplicationRequest;

  getRepositoryId(): string;
  setRepositoryId(value: string): CreateApplicationRequest;

  getBranchName(): string;
  setBranchName(value: string): CreateApplicationRequest;

  getBuildType(): BuildType;
  setBuildType(value: BuildType): CreateApplicationRequest;

  getConfig(): ApplicationConfig | undefined;
  setConfig(value?: ApplicationConfig): CreateApplicationRequest;
  hasConfig(): boolean;
  clearConfig(): CreateApplicationRequest;

  getWebsitesList(): Array<CreateWebsiteRequest>;
  setWebsitesList(value: Array<CreateWebsiteRequest>): CreateApplicationRequest;
  clearWebsitesList(): CreateApplicationRequest;
  addWebsites(value?: CreateWebsiteRequest, index?: number): CreateWebsiteRequest;

  getStartOnCreate(): boolean;
  setStartOnCreate(value: boolean): CreateApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateApplicationRequest): CreateApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: CreateApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateApplicationRequest;
  static deserializeBinaryFromReader(message: CreateApplicationRequest, reader: jspb.BinaryReader): CreateApplicationRequest;
}

export namespace CreateApplicationRequest {
  export type AsObject = {
    name: string,
    repositoryId: string,
    branchName: string,
    buildType: BuildType,
    config?: ApplicationConfig.AsObject,
    websitesList: Array<CreateWebsiteRequest.AsObject>,
    startOnCreate: boolean,
  }
}

export class ApplicationIdRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ApplicationIdRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationIdRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationIdRequest): ApplicationIdRequest.AsObject;
  static serializeBinaryToWriter(message: ApplicationIdRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationIdRequest;
  static deserializeBinaryFromReader(message: ApplicationIdRequest, reader: jspb.BinaryReader): ApplicationIdRequest;
}

export namespace ApplicationIdRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetApplicationBuildsResponse extends jspb.Message {
  getBuildsList(): Array<Build>;
  setBuildsList(value: Array<Build>): GetApplicationBuildsResponse;
  clearBuildsList(): GetApplicationBuildsResponse;
  addBuilds(value?: Build, index?: number): Build;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationBuildsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationBuildsResponse): GetApplicationBuildsResponse.AsObject;
  static serializeBinaryToWriter(message: GetApplicationBuildsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationBuildsResponse;
  static deserializeBinaryFromReader(message: GetApplicationBuildsResponse, reader: jspb.BinaryReader): GetApplicationBuildsResponse;
}

export namespace GetApplicationBuildsResponse {
  export type AsObject = {
    buildsList: Array<Build.AsObject>,
  }
}

export class GetApplicationBuildRequest extends jspb.Message {
  getBuildId(): string;
  setBuildId(value: string): GetApplicationBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationBuildRequest): GetApplicationBuildRequest.AsObject;
  static serializeBinaryToWriter(message: GetApplicationBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationBuildRequest;
  static deserializeBinaryFromReader(message: GetApplicationBuildRequest, reader: jspb.BinaryReader): GetApplicationBuildRequest;
}

export namespace GetApplicationBuildRequest {
  export type AsObject = {
    buildId: string,
  }
}

export class GetApplicationBuildLogRequest extends jspb.Message {
  getBuildId(): string;
  setBuildId(value: string): GetApplicationBuildLogRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationBuildLogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationBuildLogRequest): GetApplicationBuildLogRequest.AsObject;
  static serializeBinaryToWriter(message: GetApplicationBuildLogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationBuildLogRequest;
  static deserializeBinaryFromReader(message: GetApplicationBuildLogRequest, reader: jspb.BinaryReader): GetApplicationBuildLogRequest;
}

export namespace GetApplicationBuildLogRequest {
  export type AsObject = {
    buildId: string,
  }
}

export class SetApplicationEnvironmentVariableRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): SetApplicationEnvironmentVariableRequest;

  getKey(): string;
  setKey(value: string): SetApplicationEnvironmentVariableRequest;

  getValue(): string;
  setValue(value: string): SetApplicationEnvironmentVariableRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetApplicationEnvironmentVariableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetApplicationEnvironmentVariableRequest): SetApplicationEnvironmentVariableRequest.AsObject;
  static serializeBinaryToWriter(message: SetApplicationEnvironmentVariableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetApplicationEnvironmentVariableRequest;
  static deserializeBinaryFromReader(message: SetApplicationEnvironmentVariableRequest, reader: jspb.BinaryReader): SetApplicationEnvironmentVariableRequest;
}

export namespace SetApplicationEnvironmentVariableRequest {
  export type AsObject = {
    applicationId: string,
    key: string,
    value: string,
  }
}

export class CancelBuildRequest extends jspb.Message {
  getBuildId(): string;
  setBuildId(value: string): CancelBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CancelBuildRequest): CancelBuildRequest.AsObject;
  static serializeBinaryToWriter(message: CancelBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelBuildRequest;
  static deserializeBinaryFromReader(message: CancelBuildRequest, reader: jspb.BinaryReader): CancelBuildRequest;
}

export namespace CancelBuildRequest {
  export type AsObject = {
    buildId: string,
  }
}

export class RetryCommitBuildRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): RetryCommitBuildRequest;

  getCommit(): string;
  setCommit(value: string): RetryCommitBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryCommitBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RetryCommitBuildRequest): RetryCommitBuildRequest.AsObject;
  static serializeBinaryToWriter(message: RetryCommitBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryCommitBuildRequest;
  static deserializeBinaryFromReader(message: RetryCommitBuildRequest, reader: jspb.BinaryReader): RetryCommitBuildRequest;
}

export namespace RetryCommitBuildRequest {
  export type AsObject = {
    applicationId: string,
    commit: string,
  }
}

export enum BuildType { 
  RUNTIME = 0,
  STATIC = 1,
}
export enum ApplicationState { 
  IDLE = 0,
  DEPLOYING = 1,
  RUNNING = 2,
  ERRORED = 3,
}
export enum AuthenticationType { 
  OFF = 0,
  SOFT = 1,
  HARD = 2,
}
