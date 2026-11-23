import * as v from 'valibot'

/** Form of the editable fields of an application group. */
export const applicationGroupSchema = v.object({
  name: v.pipe(v.string(), v.nonEmpty('Enter Group Name')),
})

export type ApplicationGroupInput = v.InferInput<typeof applicationGroupSchema>
