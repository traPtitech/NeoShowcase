import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import {
  type CreateRepositoryAuth,
  type CreateRepositoryRequest,
  type Repository,
  Repository_AuthMethod,
  type UpdateRepositoryRequest,
  type UpdateRepositoryRequest_UpdateOwners,
} from '/@/api/neoshowcase/protobuf/gateway_pb'

// --- create repository

const repositoryAuthSchema = v.pipe(
  v.variant('method', [
    v.object({
      method: v.literal('none'),
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
  ]),
  v.transform((input): PartialMessage<CreateRepositoryAuth> => {
    switch (input.method) {
      case 'none': {
        return {
          auth: {
            case: input.method,
            value: {},
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
  }),
)

const createRepositorySchema = v.pipe(
  v.object({
    name: v.pipe(v.string(), v.nonEmpty('Enter Repository Name')),
    url: v.pipe(v.string(), v.nonEmpty('Enter Repository URL')),
    auth: repositoryAuthSchema,
  }),
  // Basic認証の場合、URLはhttpsでなければならない
  v.forward(
    v.partialCheck(
      [['url'], ['auth']],
      (input) => {
        if (input.auth.auth?.case === 'basic') return input.url.startsWith('https')

        // 認証方法が basic 以外の場合は常にvalid
        return true
      },
      'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
    ),
    ['url'],
  ),
  v.transform((input): PartialMessage<CreateRepositoryRequest> => {
    return {
      name: input.name,
      url: input.url,
      auth: input.auth,
    }
  }),
)

export const createRepositoryFormInitialValues = (): CreateOrUpdateRepositoryInput => ({
  type: 'create',
  form: {
    name: '',
    url: '',
    auth: {
      method: 'none',
    },
  },
})

// --- update repository

const ownersSchema = v.pipe(
  v.array(v.string()),
  v.transform(
    (input): PartialMessage<UpdateRepositoryRequest_UpdateOwners> => ({
      ownerIds: input,
    }),
  ),
)

export const updateRepositorySchema = v.pipe(
  v.object({
    id: v.string(),
    name: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository Name'))),
    url: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository URL'))),
    auth: v.optional(repositoryAuthSchema),
    ownerIds: v.optional(ownersSchema),
  }),
  // Basic認証の場合、URLはhttpsで始まる必要がある
  v.forward(
    v.partialCheck(
      [['url'], ['auth']],
      (input) => {
        if (input.auth?.auth?.case === 'basic') return input.url?.startsWith('https') ?? false

        // 認証方法が basic 以外の場合は常にvalid
        return true
      },
      'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
    ),
    ['url'],
  ),
  v.transform((input): PartialMessage<UpdateRepositoryRequest> => {
    return {
      id: input.id,
      name: input.name,
      url: input.url,
      auth: input.auth,
      ownerIds: input.ownerIds,
    }
  }),
)

/** protobuf message -> valobot schema input */
const authMethodToAuthConfig = (method: Repository_AuthMethod): v.InferInput<typeof repositoryAuthSchema> => {
  switch (method) {
    case Repository_AuthMethod.NONE: {
      return {
        method: 'none',
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
            keyId: undefined,
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

export const updateRepositoryFormInitialValues = (input: Repository): CreateOrUpdateRepositoryInput => {
  return {
    type: 'update',
    form: {
      id: input.id,
      name: input.name,
      url: input.url,
      auth: authMethodToAuthConfig(input.authMethod),
      ownerIds: input.ownerIds,
    },
  }
}

export const createOrUpdateRepositorySchema = v.variant('type', [
  v.object({
    type: v.literal('create'),
    form: createRepositorySchema,
  }),
  v.object({
    type: v.literal('update'),
    form: updateRepositorySchema,
  }),
])

export type CreateOrUpdateRepositoryInput = v.InferInput<typeof createOrUpdateRepositorySchema>
export type CreateOrUpdateRepositoryOutput = v.InferOutput<typeof createOrUpdateRepositorySchema>

export const handleSubmitRepositoryForm = (
  input: CreateOrUpdateRepositoryInput,
  handler: (output: CreateOrUpdateRepositoryOutput['form']) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateRepositorySchema, input)
  return handler(result.form)
}
