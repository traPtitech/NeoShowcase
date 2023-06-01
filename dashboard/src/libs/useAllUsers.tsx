import { createResource } from 'solid-js'
import { client } from '/@/libs/api'

const [allUsersResource, { mutate: mutateAllUsers, refetch: refetchAllUsers }] = createResource(async () => {
  const allUsersRes = await client.getUsers({})
  return allUsersRes.users
})

export { allUsersResource, mutateAllUsers, refetchAllUsers }
