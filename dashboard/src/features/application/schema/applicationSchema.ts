import type { FormStore } from '@formisch/solid/v2'
import * as v from 'valibot'
import type {
  Application,
  CreateApplicationRequest,
  UpdateApplicationRequest,
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
  v.transform((input): CreateApplicationRequest => {
    return {
      $typeName: 'neoshowcase.protobuf.CreateApplicationRequest',
      name: input.name,
      repositoryId: input.repositoryId,
      refName: input.refName,
      config: input.config,
      websites: input.websites ?? [],
      portPublications: input.portPublications ?? [],
      startOnCreate: input.startOnCreate,
    }
  }),
)

type CreateApplicationOutput = v.InferOutput<typeof createApplicationSchema>

// --- update application

export const updateApplicationSchema = v.pipe(
  v.object({
    id: v.string(),
    name: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Application Name'))),
    refName: v.optional(v.pipe(v.string(), v.nonEmpty('Enter Branch Name'))),
    config: v.optional(applicationConfigSchema),
    websites: v.optional(v.array(createWebsiteSchema)),
    portPublications: v.optional(v.array(portPublicationSchema)),
    startOnCreate: v.optional(v.boolean()),
  }),
  v.transform(
    (input): UpdateApplicationRequest => ({
      $typeName: 'neoshowcase.protobuf.UpdateApplicationRequest',
      id: input.id,
      name: input.name,
      refName: input.refName,
      config: input.config,
      websites: input.websites
        ? {
            $typeName: 'neoshowcase.protobuf.UpdateApplicationRequest.UpdateWebsites',
            websites: input.websites,
          }
        : undefined,
      portPublications: input.portPublications
        ? {
            $typeName: 'neoshowcase.protobuf.UpdateApplicationRequest.UpdatePorts',
            portPublications: input.portPublications,
          }
        : undefined,
    }),
  ),
)

type UpdateApplicationOutput = v.InferOutput<typeof updateApplicationSchema>

export const getInitialValueOfCreateAppForm = (): CreateOrUpdateApplicationInput => ({
  type: 'create',
  form: {
    name: '',
    repositoryId: '',
    refName: '',
    config: {},
    websites: [],
    portPublications: [],
    startOnCreate: true,
  },
})

export const getInitialValueOfUpdateAppForm = (input: Application): CreateOrUpdateApplicationInput => ({
  type: 'update',
  form: {
    id: input.id,
    name: input.name,
    refName: input.refName,
    config: input.config ? configMessageToSchema(input.config) : undefined,
    websites: input.websites.map((w) => websiteMessageToSchema(w)),
    portPublications: input.portPublications.map((p) => portPublicationMessageToSchema(p)),
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
export type ApplicationFormStore = FormStore<typeof createOrUpdateApplicationSchema>

/** The create form is only ever submitted with `type: 'create'`; narrows the validated output. */
export const toCreateApplicationRequest = (output: CreateOrUpdateApplicationOutput): CreateApplicationOutput => {
  if (output.type !== 'create') throw new Error('The submitted application form must be of type "create"')
  return output.form
}

/** The update form is only ever submitted with `type: 'update'`; narrows the validated output. */
export const toUpdateApplicationRequest = (output: CreateOrUpdateApplicationOutput): UpdateApplicationOutput => {
  if (output.type !== 'update') throw new Error('The submitted application form must be of type "update"')
  return output.form
}
