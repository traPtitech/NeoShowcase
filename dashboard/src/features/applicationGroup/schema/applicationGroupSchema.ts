import * as v from 'valibot'

/** Form of the editable fields of an application group. */
export const applicationGroupSchema = v.object({
  name: v.pipe(
    v.string(),
    v.nonEmpty('Enter Group Name'),
    // The server counts characters as code points, while v.maxLength counts UTF-16 code units.
    v.check((name) => [...name].length <= 100, 'Group Name must be 100 characters or less'),
  ),
})

export type ApplicationGroupInput = v.InferInput<typeof applicationGroupSchema>
