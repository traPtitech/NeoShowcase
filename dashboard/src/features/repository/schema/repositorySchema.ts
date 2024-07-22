import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import {
  type CreateRepositoryAuth,
  type CreateRepositoryRequest,
  type Repository,
  Repository_AuthMethod,
  type UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'

// --- create repository

const repositoryAuthSchema = v.variant('method', [
  v.object({
    method: v.literal('none'),
    value: v.object({
      none: v.object({}),
    }),
  }),
  v.object({
    method: v.literal('basic'),
    value: v.object({
      basic: v.object({
        username: v.pipe(v.string(), v.nonEmpty('Enter UserName')),
        password: v.pipe(v.string(), v.nonEmpty('Enter Password')),
      }),
    }),
  }),
  v.object({
    method: v.literal('ssh'),
    value: v.object({
      ssh: v.object({
        keyId: v.optional(v.string()), // undefinedの場合はNeoShowcase全体共通の公開鍵が使用される
      }),
    }),
  }),
])

const createRepositorySchema = v.pipe(
  v.object({
    type: v.literal('create'),
    name: v.pipe(v.string(), v.nonEmpty('Enter Repository Name')),
    url: v.pipe(v.string(), v.nonEmpty('Enter Repository URL')),
    auth: repositoryAuthSchema,
  }),
  // Basic認証の場合、URLはhttpsでなければならない
  v.forward(
    v.partialCheck(
      [['url'], ['auth', 'method']],
      (input) => {
        if (input.auth.method === 'basic') return input.url.startsWith('https')

        // 認証方法が basic 以外の場合は常にvalid
        return true
      },
      'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
    ),
    ['url'],
  ),
)
type CreateRepositorySchema = v.InferInput<typeof createRepositorySchema>

export const createRepositoryFormInitialValues = (): CreateOrUpdateRepositorySchema =>
  ({
    type: 'create',
    name: '',
    url: '',
    auth: {
      method: 'none',
      value: {
        none: {},
      },
    },
  }) satisfies CreateRepositorySchema

/** valobot schema -> protobuf message */
const repositoryAuthSchemaToMessage = (
  input: v.InferInput<typeof repositoryAuthSchema>,
): PartialMessage<CreateRepositoryAuth> => {
  switch (input.method) {
    case 'none': {
      return {
        auth: {
          case: input.method,
          value: input.value.none,
        },
      }
    }
    case 'basic': {
      return {
        auth: {
          case: input.method,
          value: input.value.basic,
        },
      }
    }
    case 'ssh': {
      return {
        auth: {
          case: input.method,
          value: input.value.ssh,
        },
      }
    }
  }
}

/** valobot schema -> protobuf message */
export const convertCreateRepositoryInput = (
  input: CreateOrUpdateRepositorySchema,
): PartialMessage<CreateRepositoryRequest> => {
  if (input.type !== 'create')
    throw new Error("The type of input passed to convertCreateRepositoryInput must be 'create'")

  return {
    ...input,
    auth: repositoryAuthSchemaToMessage(input.auth),
  }
}

// --- update repository

const ownersSchema = v.array(v.string())

export const updateRepositorySchema = v.pipe(
  v.object({
    type: v.literal('update'),
    id: v.string(),
    name: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository Name'))),
    url: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository URL'))),
    auth: v.optional(repositoryAuthSchema),
    ownerIds: v.optional(ownersSchema),
  }),
  // Basic認証の場合、URLはhttpsでなければならない
  v.forward(
    v.partialCheck(
      [['url'], ['auth', 'method']],
      (input) => {
        if (input.auth?.method === 'basic') return input.url?.startsWith('https') ?? false

        // 認証方法が basic 以外の場合は常にvalid
        return true
      },
      'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
    ),
    ['url'],
  ),
)

type UpdateRepositorySchema = v.InferInput<typeof updateRepositorySchema>

/** protobuf message -> valobot schema */
const authMethodToAuthConfig = (method: Repository_AuthMethod): v.InferInput<typeof repositoryAuthSchema> => {
  switch (method) {
    case Repository_AuthMethod.NONE: {
      return {
        method: 'none',
        value: {
          none: {},
        },
      }
    }
    case Repository_AuthMethod.BASIC: {
      return {
        method: 'basic',
        value: {
          basic: {
            username: '',
            password: '',
          },
        },
      }
    }
    case Repository_AuthMethod.SSH: {
      return {
        method: 'ssh',
        value: {
          ssh: {
            keyId: '',
          },
        },
      }
    }
    default: {
      const _unreachable: never = method
      throw new Error('unknown repository auth method')
    }
  }
}

export const updateRepositoryFormInitialValues = (input: Repository): CreateOrUpdateRepositorySchema => {
  return {
    type: 'update',
    id: input.id,
    name: input.name,
    url: input.url,
    auth: authMethodToAuthConfig(input.authMethod),
    ownerIds: input.ownerIds,
  } satisfies UpdateRepositorySchema
}

/** valobot schema -> protobuf message */
export const convertUpdateRepositoryInput = (
  input: CreateOrUpdateRepositorySchema,
): PartialMessage<UpdateRepositoryRequest> => {
  if (input.type !== 'update')
    throw new Error("The type of input passed to convertCreateRepositoryInput must be 'create'")

  return {
    ...input,
    auth: input.auth ? repositoryAuthSchemaToMessage(input.auth) : undefined,
    ownerIds: input.ownerIds
      ? {
          ownerIds: input.ownerIds,
        }
      : undefined,
  }
}

export const createOrUpdateRepositorySchema = v.variant('type', [createRepositorySchema, updateRepositorySchema])

export type CreateOrUpdateRepositorySchema = v.InferInput<typeof createOrUpdateRepositorySchema>
