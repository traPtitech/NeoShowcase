import { createMemo, createResource, createRoot } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client } from '/@/libs/api'

const [users, { mutate: mutateUsers, refetch: refetchUsers }] = createResource(async () => {
  const getUsersRes = await client.getUsers({})
  return getUsersRes.users
})

export { users, mutateUsers, refetchUsers }

// keyにID, valueにユーザー情報を持つMap
// createRootを使わない場合``computations created outside a `createRoot` or `render` will never be disposed``の警告が出る
// see: https://www.solidjs.com/docs/latest#createroot
const usersMap = createRoot(() =>
  createMemo(() => {
    if (users.latest !== undefined) return new Map(users().map((user) => [user.id, user]))
    return new Map<string, User>()
  }),
)

export const userFromId = (id: string) => {
  const user = usersMap().get(id)
  if (user) return user
  throw new Error(`userFromId: user not found: ${id}`)
}
