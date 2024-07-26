import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import {
  type Application,
  AuthenticationType,
  type AvailableDomain,
  type CreateApplicationRequest,
  type CreateWebsiteRequest,
  PortPublicationProtocol,
  type UpdateApplicationRequest,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { applicationConfigSchema, configMessageToSchema, configSchemaToMessage } from './applicationConfigSchema'

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

type CreateApplicationInput = v.InferInput<typeof createApplicationSchema>

export const createApplicationFormInitialValues = (): CreateOrUpdateApplicationInput =>
  ({
    type: 'create',
    name: '',
    repositoryId: '',
    refName: '',
    config: {},
    websites: [],
    portPublications: [],
    startOnCreate: false,
  }) satisfies CreateApplicationInput

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
  input: CreateOrUpdateApplicationInput,
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
  websites: v.optional(v.array(createWebsiteSchema)),
  portPublications: v.optional(v.array(portPublicationSchema)),
  ownerIds: v.optional(ownersSchema),
})

type UpdateApplicationInput = v.InferInput<typeof updateApplicationSchema>

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

export const updateApplicationFormInitialValues = (input: Application): CreateOrUpdateApplicationInput => {
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
  } satisfies UpdateApplicationInput
}

/** valobot schema -> protobuf message */
export const convertUpdateApplicationInput = (
  input: CreateOrUpdateApplicationInput,
): PartialMessage<UpdateApplicationRequest> => {
  if (input.type !== 'update')
    throw new Error("The type of input passed to convertUpdateApplicationInput must be 'create'")

  return {
    id: input.id,
    name: input.name,
    repositoryId: input.repositoryId,
    refName: input.refName,
    config: input.config ? configSchemaToMessage(input.config) : undefined,
    websites: input.websites
      ? {
          websites: input.websites?.map((w) => createWebsiteSchemaToMessage(w)),
        }
      : undefined,
    portPublications: input.portPublications ? { portPublications: input.portPublications } : undefined,
    ownerIds: input.ownerIds
      ? {
          ownerIds: input.ownerIds,
        }
      : undefined,
  }
}

export const createOrUpdateApplicationSchema = v.variant('type', [createApplicationSchema, updateApplicationSchema])

export type CreateOrUpdateApplicationInput = v.InferInput<typeof createOrUpdateApplicationSchema>
