import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import type {
  Application,
  CreateApplicationRequest,
  UpdateApplicationRequest,
  UpdateApplicationRequest_UpdateOwners,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationConfigSchema, configMessageToSchema } from './applicationConfigSchema'
import { portPublicationMessageToSchema, portPublicationSchema } from './portPublicationSchema'
import { createWebsiteSchema, websiteMessageToSchema } from './websiteSchema'

// --- create application

const createApplicationSchema = v.pipe(
  v.object({
    name: v.pipe(v.string(), v.nonEmpty('Enter Application Name')),
    repositoryId: v.string(),
    refName: v.pipe(v.string(), v.nonEmpty('Enter Branch Name')),
    config: applicationConfigSchema,
    websites: v.optional(v.array(createWebsiteSchema)),
    portPublications: v.optional(v.array(portPublicationSchema)),
    startOnCreate: v.boolean(),
  }),
  v.transform((input): PartialMessage<CreateApplicationRequest> => input),
)

type CreateApplicationOutput = v.InferOutput<typeof createApplicationSchema>

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
    startOnCreate: v.optional(v.boolean()),
  }),
  v.transform(
    (input): PartialMessage<UpdateApplicationRequest> => ({
      id: input.id,
      name: input.name,
      repositoryId: input.repositoryId,
      refName: input.refName,
      config: input.config,
      websites: input.websites
        ? {
            websites: input.websites,
          }
        : undefined,
      portPublications: input.portPublications ? { portPublications: input.portPublications } : undefined,
      ownerIds: input.ownerIds,
    }),
  ),
)

type UpdateApplicationOutput = v.InferOutput<typeof updateApplicationSchema>

export const updateApplicationFormInitialValues = (input: Application): CreateOrUpdateApplicationInput => ({
  type: 'update',
  form: {
    id: input.id,
    name: input.name,
    repositoryId: input.repositoryId,
    refName: input.refName,
    config: input.config ? configMessageToSchema(input.config) : undefined,
    websites: input.websites.map((w) => websiteMessageToSchema(w)),
    portPublications: input.portPublications.map((p) => portPublicationMessageToSchema(p)),
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
  handler: (output: CreateApplicationOutput) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateApplicationSchema, input)
  if (result.type !== 'create')
    throw new Error('The type of input passed to handleSubmitCreateApplicationForm must be "create"')
  return handler(result.form)
}

export const handleSubmitUpdateApplicationForm = (
  input: CreateOrUpdateApplicationInput,
  handler: (output: UpdateApplicationOutput) => Promise<unknown>,
) => {
  const result = v.parse(createOrUpdateApplicationSchema, input)
  if (result.type !== 'update')
    throw new Error('The type of input passed to handleSubmitCreateApplicationForm must be "update"')
  return handler(result.form)
}
