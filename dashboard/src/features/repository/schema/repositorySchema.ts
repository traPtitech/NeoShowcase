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
  v.transform((input): CreateRepositoryAuth => {
    switch (input.method) {
      case 'none': {
        return {
          $typeName: 'neoshowcase.protobuf.CreateRepositoryAuth',
          auth: {
            case: input.method,
            value: { $typeName: 'google.protobuf.Empty' },
          },
        }
      }
      case 'basic': {
        return {
          $typeName: 'neoshowcase.protobuf.CreateRepositoryAuth',
          auth: {
            case: input.method,
            value: {
              $typeName: 'neoshowcase.protobuf.CreateRepositoryAuthBasic',
              ...input.value.basic,
            },
          },
        }
      }
      case 'ssh': {
        return {
          $typeName: 'neoshowcase.protobuf.CreateRepositoryAuth',
          auth: {
            case: input.method,
            value: {
              $typeName: 'neoshowcase.protobuf.CreateRepositoryAuthSSH',
              keyId: input.value.ssh.keyId ?? '',
            },
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
  v.transform(
    (input): CreateRepositoryRequest => ({
      $typeName: 'neoshowcase.protobuf.CreateRepositoryRequest',
      name: input.name,
      url: input.url,
      auth: input.auth,
    }),
  ),
)

type CreateRepositoryOutput = v.InferOutput<typeof createRepositorySchema>

export const getInitialValueOfCreateRepoForm = (): CreateOrUpdateRepositoryInput => ({
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
    (input): UpdateRepositoryRequest_UpdateOwners => ({
      $typeName: 'neoshowcase.protobuf.UpdateRepositoryRequest.UpdateOwners',
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
  v.transform(
    (input): UpdateRepositoryRequest => ({
      $typeName: 'neoshowcase.protobuf.UpdateRepositoryRequest',
      id: input.id,
      name: input.name,
      url: input.url,
      auth: input.auth,
      ownerIds: input.ownerIds,
    }),
  ),
)

type UpdateRepositoryOutput = v.InferOutput<typeof updateRepositorySchema>

/** protobuf message -> valobot schema input */
const authMethodToAuthConfig = (method: Repository_AuthMethod): v.InferInput<typeof repositoryAuthSchema> => {
  // ts-patternでnumber型のenumの網羅性チェックができないためswitchを使用
  // https://zenn.dev/tacrew/articles/c58ab324aee960#number型のenumにおいて網羅性チェックが不十分
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

export const getInitialValueOfUpdateRepoForm = (input: Repository): CreateOrUpdateRepositoryInput => {
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

export const handleSubmitCreateRepositoryForm = (
  input: CreateOrUpdateRepositoryInput,
  handler: (output: CreateRepositoryOutput) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateRepositorySchema, input)
  if (result.type !== 'create')
    throw new Error('The type of input passed to handleSubmitCreateRepositoryForm must be "create"')
  return handler(result.form)
}

export const handleSubmitUpdateRepositoryForm = (
  input: CreateOrUpdateRepositoryInput,
  handler: (output: UpdateRepositoryOutput) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateRepositorySchema, input)
  if (result.type !== 'update')
    throw new Error('The type of input passed to handleSubmitCreateRepositoryForm must be "update"')
  return handler(result.form)
}
