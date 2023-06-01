import { createMemo, createResource } from 'solid-js'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client } from '/@/libs/api'

const [allUsersResource, { mutate: mutateAllUsers, refetch: refetchAllUsers }] = createResource(async () => {
  const allUsersRes = await client.getUsers({})
  return allUsersRes.users
})

export { allUsersResource, mutateAllUsers, refetchAllUsers }

// keyにID, valueにユーザー情報を持つMap
const allUsersMap = createMemo(() => {
  if (!allUsersResource()) return new Map<string, User>()
  return new Map(allUsersResource().map((user) => [user.id, user]))
})

export const userFromId = (id: string) => allUsersMap().get(id)
