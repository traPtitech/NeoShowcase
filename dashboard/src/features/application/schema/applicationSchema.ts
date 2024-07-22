import type { PartialMessage } from '@bufbuild/protobuf'
import { match } from 'ts-pattern'
import * as v from 'valibot'
import {
  type Application,
  type ApplicationConfig,
  AuthenticationType,
  type AvailableDomain,
  type CreateApplicationRequest,
  type CreateWebsiteRequest,
  PortPublicationProtocol,
  type UpdateApplicationRequest,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'

const runtimeConfigSchema = v.object({
  useMariadb: v.boolean(),
  useMongodb: v.boolean(),
  entrypoint: v.string(),
  command: v.string(),
})
const staticConfigSchema = v.object({
  artifactPath: v.pipe(v.string(), v.nonEmpty('Enter Artifact Path')),
  spa: v.boolean(),
})

const buildpackConfigSchema = v.object({
  context: v.string(),
})
const cmdConfigSchema = v.object({
  baseImage: v.string(),
  buildCmd: v.string(),
})
const dockerfileConfigSchema = v.object({
  dockerfileName: v.pipe(v.string(), v.nonEmpty('Enter Dockerfile Name')),
  context: v.string(),
})

const applicationConfigSchema = v.object({
  deployConfig: v.pipe(
    v.optional(
      v.variant('type', [
        v.object({
          type: v.literal('runtime'),
          value: v.object({
            runtime: runtimeConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('static'),
          value: v.object({
            static: staticConfigSchema,
          }),
        }),
      ]),
    ),
    // アプリ作成時には最初undefinedになっているが、submit時にはundefinedで無い必要がある
    v.check((input) => !!input, 'Select Deploy Type'),
  ),
  buildConfig: v.pipe(
    v.optional(
      v.variant('type', [
        v.object({
          type: v.literal('buildpack'),
          value: v.object({
            buildpack: buildpackConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('cmd'),
          value: v.object({
            cmd: cmdConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('dockerfile'),
          value: v.object({
            dockerfile: dockerfileConfigSchema,
          }),
        }),
      ]),
    ),
    // アプリ作成時には最初undefinedになっているが、submit時にはundefinedで無い必要がある
    v.check((input) => !!input, 'Select Build Type'),
  ),
})

type ApplicationConfigInput = v.InferInput<typeof applicationConfigSchema>

const createWebsiteSchema = v.object({
  state: v.union([v.literal('noChange'), v.literal('readyToChange'), v.literal('readyToDelete'), v.literal('added')]),
  subdomain: v.string(),
  domain: v.string(),
  pathPrefix: v.string(),
  stripPrefix: v.boolean(),
  https: v.boolean(),
  h2c: v.boolean(),
  httpPort: v.pipe(v.number(), v.integer()),
  authentication: v.enum(AuthenticationType),
})

const portPublicationSchema = v.object({
  internetPort: v.pipe(v.number(), v.integer()),
  applicationPort: v.pipe(v.number(), v.integer()),
  protocol: v.enum(PortPublicationProtocol),
})

// --- create application

const createApplicationSchema = v.object({
  type: v.literal('create'),
  name: v.pipe(v.string(), v.nonEmpty('Enter Application Name')),
  repositoryId: v.string(),
  refName: v.pipe(v.string(), v.nonEmpty('Enter Branch Name')),
  config: applicationConfigSchema,
  websites: v.array(createWebsiteSchema),
  portPublications: v.array(portPublicationSchema),
  startOnCreate: v.boolean(),
})

type CreateApplicationSchema = v.InferInput<typeof createApplicationSchema>

export const createApplicationFormInitialValues = (): CreateOrUpdateApplicationSchema =>
  ({
    type: 'create',
    name: '',
    repositoryId: '',
    refName: '',
    config: {},
    websites: [],
    portPublications: [],
    startOnCreate: false,
  }) satisfies CreateApplicationSchema

/** valobot schema -> protobuf message */
const configSchemaToMessage = (
  input: v.InferInput<typeof applicationConfigSchema>,
): PartialMessage<ApplicationConfig> => {
  return match([input.deployConfig, input.buildConfig])
    .returnType<PartialMessage<ApplicationConfig>>()
    .with(
      [
        {
          type: 'runtime',
        },
        {
          type: 'buildpack',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeBuildpack',
            value: {
              ...buildConfig.value.buildpack,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      },
    )
    .with(
      [
        {
          type: 'runtime',
        },
        {
          type: 'cmd',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeCmd',
            value: {
              ...buildConfig.value.cmd,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      },
    )
    .with(
      [
        {
          type: 'runtime',
        },
        {
          type: 'dockerfile',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeDockerfile',
            value: {
              ...buildConfig.value.dockerfile,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      },
    )
    .with(
      [
        {
          type: 'static',
        },
        {
          type: 'buildpack',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticBuildpack',
            value: {
              ...buildConfig.value.buildpack,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      },
    )
    .with(
      [
        {
          type: 'static',
        },
        {
          type: 'cmd',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticCmd',
            value: {
              ...buildConfig.value.cmd,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      },
    )
    .with(
      [
        {
          type: 'static',
        },
        {
          type: 'dockerfile',
        },
      ],
      ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticDockerfile',
            value: {
              ...buildConfig.value.dockerfile,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      },
    )
    .otherwise(() => ({
      buildConfig: { case: undefined },
    }))
}

const createWebsiteSchemaToMessage = (
  input: v.InferInput<typeof createWebsiteSchema>,
): PartialMessage<CreateWebsiteRequest> => {
  const { domain, subdomain, ...rest } = input

  // wildcard domainならsubdomainとdomainを結合
  const fqdn = input.domain.startsWith('*')
    ? `${input.subdomain}${input.domain.replace(/\*/g, '')}`
    : // non-wildcard domainならdomainをそのまま使う
      input.domain

  return {
    ...rest,
    fqdn,
  }
}

/** valobot schema -> protobuf message */
export const convertCreateApplicationInput = (
  input: CreateOrUpdateApplicationSchema,
): PartialMessage<CreateApplicationRequest> => {
  if (input.type !== 'create')
    throw new Error("The type of input passed to convertCreateApplicationInput must be 'create'")

  const { type: _type, config, websites, ...rest } = input

  return {
    ...rest,
    config: configSchemaToMessage(input.config),
    websites: input.websites.map((w) => createWebsiteSchemaToMessage(w)),
  }
}

// --- update application

const ownersSchema = v.array(v.string())

export const updateApplicationSchema = v.object({
  type: v.literal('update'),
  id: v.string(),
  name: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Application Name'))),
  repositoryId: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository ID'))),
  refName: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Branch Name'))),
  config: v.optional(applicationConfigSchema),
  websites: v.array(createWebsiteSchema),
  portPublications: v.array(portPublicationSchema),
  ownerIds: v.optional(ownersSchema),
})

type UpdateApplicationSchema = v.InferInput<typeof updateApplicationSchema>

/** protobuf message -> valobot schema */
const configMessageToSchema = (config: ApplicationConfig): ApplicationConfigInput => {
  let deployConfig: ApplicationConfigInput['deployConfig']
  const _case = config.buildConfig.case
  switch (_case) {
    case 'runtimeBuildpack':
    case 'runtimeDockerfile':
    case 'runtimeCmd': {
      deployConfig = {
        type: 'runtime',
        value: {
          runtime: config.buildConfig.value.runtimeConfig ?? {
            command: '',
            entrypoint: '',
            useMariadb: false,
            useMongodb: false,
          },
        },
      }
      break
    }
    case 'staticBuildpack':
    case 'staticDockerfile':
    case 'staticCmd': {
      deployConfig = {
        type: 'static',
        value: {
          static: config.buildConfig.value.staticConfig ?? {
            spa: false,
            artifactPath: '',
          },
        },
      }
      break
    }
    case undefined: {
      break
    }
    default: {
      const _unreachable: never = _case
      throw new Error('unknown application build config case')
    }
  }

  let buildConfig: ApplicationConfigInput['buildConfig']
  switch (_case) {
    case 'runtimeBuildpack':
    case 'staticBuildpack': {
      buildConfig = {
        type: 'buildpack',
        value: {
          buildpack: config.buildConfig.value,
        },
      }
      break
    }
    case 'runtimeCmd':
    case 'staticCmd': {
      buildConfig = {
        type: 'cmd',
        value: {
          cmd: config.buildConfig.value,
        },
      }
      break
    }
    case 'runtimeDockerfile':
    case 'staticDockerfile': {
      buildConfig = {
        type: 'dockerfile',
        value: {
          dockerfile: config.buildConfig.value,
        },
      }
      break
    }
    case undefined: {
      break
    }
    default: {
      const _unreachable: never = _case
      throw new Error('unknown application build config case')
    }
  }

  return {
    deployConfig,
    buildConfig,
  }
}

const extractSubdomain = (
  fqdn: string,
  availableDomains: AvailableDomain[],
): {
  subdomain: string
  domain: string
} => {
  const nonWildcardDomains = availableDomains.filter((d) => !d.domain.startsWith('*'))
  const wildcardDomains = availableDomains.filter((d) => d.domain.startsWith('*'))

  const matchNonWildcardDomain = nonWildcardDomains.find((d) => fqdn === d.domain)
  if (matchNonWildcardDomain !== undefined) {
    return {
      subdomain: '',
      domain: matchNonWildcardDomain.domain,
    }
  }

  const matchDomain = wildcardDomains.find((d) => fqdn.endsWith(d.domain.replace(/\*/g, '')))
  if (matchDomain === undefined) {
    const fallbackDomain = availableDomains.at(0)
    if (fallbackDomain === undefined) throw new Error('No domain available')
    return {
      subdomain: '',
      domain: fallbackDomain.domain,
    }
  }
  return {
    subdomain: fqdn.slice(0, -matchDomain.domain.length + 1),
    domain: matchDomain.domain,
  }
}

const websiteMessageToSchema = (website: Website): v.InferInput<typeof createWebsiteSchema> => {
  const availableDomains = systemInfo()?.domains ?? []

  const { domain, subdomain } = extractSubdomain(website.fqdn, availableDomains)

  return {
    state: 'noChange',
    domain,
    subdomain,
    pathPrefix: website.pathPrefix,
    stripPrefix: website.stripPrefix,
    https: website.https,
    h2c: website.h2c,
    httpPort: website.httpPort,
    authentication: website.authentication,
  }
}

export const updateApplicationFormInitialValues = (input: Application): CreateOrUpdateApplicationSchema => {
  return {
    type: 'update',
    id: input.id,
    name: input.name,
    repositoryId: input.repositoryId,
    refName: input.refName,
    config: input.config ? configMessageToSchema(input.config) : undefined,
    websites: input.websites.map((w) => websiteMessageToSchema(w)),
    portPublications: input.portPublications,
    ownerIds: input.ownerIds,
  } satisfies UpdateApplicationSchema
}

/** valobot schema -> protobuf message */
export const convertUpdateApplicationInput = (
  input: CreateOrUpdateApplicationSchema,
): PartialMessage<UpdateApplicationRequest> => {
  if (input.type !== 'update')
    throw new Error("The type of input passed to convertUpdateApplicationInput must be 'create'")

  return {
    id: input.id,
    name: input.name,
    repositoryId: input.repositoryId,
    refName: input.refName,
    config: input.config ? configSchemaToMessage(input.config) : undefined,
    websites: {
      websites: input.websites.map((w) => createWebsiteSchemaToMessage(w)),
    },
    portPublications: { portPublications: input.portPublications },
    ownerIds: {
      ownerIds: input.ownerIds,
    },
  }
}

export const createOrUpdateApplicationSchema = v.variant('type', [createApplicationSchema, updateApplicationSchema])

export type CreateOrUpdateApplicationSchema = v.InferInput<typeof createOrUpdateApplicationSchema>
