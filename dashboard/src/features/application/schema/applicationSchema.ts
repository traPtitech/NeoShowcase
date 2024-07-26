import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import {
  type Application,
  AuthenticationType,
  type AvailableDomain,
  type CreateApplicationRequest,
  type CreateWebsiteRequest,
  type PortPublication,
  PortPublicationProtocol,
  type UpdateApplicationRequest,
  type UpdateApplicationRequest_UpdateOwners,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { applicationConfigSchema, configMessageToSchema } from './applicationConfigSchema'

const createWebsiteSchema = v.pipe(
  v.object({
    state: v.union([v.literal('noChange'), v.literal('readyToChange'), v.literal('readyToDelete'), v.literal('added')]),
    subdomain: v.string(),
    domain: v.string(),
    pathPrefix: v.string(),
    stripPrefix: v.boolean(),
    https: v.boolean(),
    h2c: v.boolean(),
    httpPort: v.pipe(v.number(), v.integer()),
    authentication: v.enum(AuthenticationType),
  }),
  v.transform((input): PartialMessage<CreateWebsiteRequest> => {
    // wildcard domainならsubdomainとdomainを結合
    const fqdn = input.domain.startsWith('*')
      ? `${input.subdomain}${input.domain.replace(/\*/g, '')}`
      : // non-wildcard domainならdomainをそのまま使う
        input.domain

    return {
      fqdn,
      authentication: input.authentication,
      h2c: input.h2c,
      httpPort: input.httpPort,
      https: input.https,
      pathPrefix: input.pathPrefix,
      stripPrefix: input.stripPrefix,
    }
  }),
)

const portPublicationSchema = v.pipe(
  v.object({
    internetPort: v.pipe(v.number(), v.integer()),
    applicationPort: v.pipe(v.number(), v.integer()),
    protocol: v.enum(PortPublicationProtocol),
  }),
  v.transform((input): PartialMessage<PortPublication> => input),
)

// --- create application

const createApplicationSchema = v.pipe(
  v.object({
    name: v.pipe(v.string(), v.nonEmpty('Enter Application Name')),
    repositoryId: v.string(),
    refName: v.pipe(v.string(), v.nonEmpty('Enter Branch Name')),
    config: applicationConfigSchema,
    websites: v.array(createWebsiteSchema),
    portPublications: v.array(portPublicationSchema),
    startOnCreate: v.boolean(),
  }),
  v.transform(
    (input): PartialMessage<CreateApplicationRequest> => ({
      name: input.name,
      repositoryId: input.repositoryId,
      refName: input.refName,
      config: input.config,
      portPublications: input.portPublications,
      websites: input.websites,
      startOnCreate: input.startOnCreate,
    }),
  ),
)

export const createApplicationFormInitialValues = (): CreateOrUpdateApplicationInput => ({
  type: 'create',
  form: {
    name: '',
    repositoryId: '',
    refName: '',
    config: {},
    websites: [],
    portPublications: [],
    startOnCreate: false,
  },
})

// --- update application

const ownersSchema = v.pipe(
  v.array(v.string()),
  v.transform(
    (input): PartialMessage<UpdateApplicationRequest_UpdateOwners> => ({
      ownerIds: input,
    }),
  ),
)

export const updateApplicationSchema = v.pipe(
  v.object({
    id: v.string(),
    name: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Application Name'))),
    repositoryId: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Repository ID'))),
    refName: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Branch Name'))),
    config: v.optional(applicationConfigSchema),
    websites: v.optional(v.array(createWebsiteSchema)),
    portPublications: v.optional(v.array(portPublicationSchema)),
    ownerIds: v.optional(ownersSchema),
  }),
  v.transform(
    (input): PartialMessage<UpdateApplicationRequest> => ({
      id: input.id,
      name: input.name,
      repositoryId: input.repositoryId,
      refName: input.refName,
      config: input.config,
      websites: {
        websites: input.websites,
      },
      portPublications: input.portPublications ? { portPublications: input.portPublications } : undefined,
      ownerIds: input.ownerIds,
    }),
  ),
)

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

export const updateApplicationFormInitialValues = (input: Application): CreateOrUpdateApplicationInput => ({
  type: 'update',
  form: {
    id: input.id,
    name: input.name,
    repositoryId: input.repositoryId,
    refName: input.refName,
    config: input.config ? configMessageToSchema(input.config) : undefined,
    websites: input.websites.map((w) => websiteMessageToSchema(w)),
    portPublications: input.portPublications,
    ownerIds: input.ownerIds,
  },
})

export const createOrUpdateApplicationSchema = v.variant('type', [
  v.object({
    type: v.literal('create'),
    form: createApplicationSchema,
  }),
  v.object({
    type: v.literal('update'),
    form: updateApplicationSchema,
  }),
])

export type CreateOrUpdateApplicationInput = v.InferInput<typeof createOrUpdateApplicationSchema>
export type CreateOrUpdateApplicationOutput = v.InferOutput<typeof createOrUpdateApplicationSchema>

export const handleSubmitCreateApplicationForm = (
  input: CreateOrUpdateApplicationInput,
  handler: (output: CreateOrUpdateApplicationOutput['form']) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateApplicationSchema, input)
  if (result.type !== 'create')
    throw new Error('The type of input passed to handleSubmitCreateApplicationForm must be "create"')
  return handler(result.form)
}

export const handleSubmitUpdateApplicationForm = (
  input: CreateOrUpdateApplicationInput,
  handler: (output: CreateOrUpdateApplicationOutput['form']) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateApplicationSchema, input)
  if (result.type !== 'update')
    throw new Error('The type of input passed to handleSubmitCreateApplicationForm must be "update"')
  return handler(result.form)
}
