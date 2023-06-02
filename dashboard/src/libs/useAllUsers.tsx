import { createMemo, createResource, createRoot } from 'solid-js'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client } from '/@/libs/api'

const [users, { mutate: mutateUsers, refetch: refetchUsers }] = createResource(async () => {
  const getUsersRes = await client.getUsers({})
  return getUsersRes.users
})

export { users, mutateUsers, refetchUsers }

// keyにID, valueにユーザー情報を持つMap
const usersMap = createRoot(() =>
  createMemo(() => {
    if (!users()) return new Map<string, User>()
    return new Map(users().map((user) => [user.id, user]))
  }),
)

export const userFromId = (id: string) => usersMap().get(id)
