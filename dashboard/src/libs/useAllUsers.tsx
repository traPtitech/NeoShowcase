import { createResource } from 'solid-js'
import { client } from '/@/libs/api'

/**
 * Create a all users resource:
 * ```typescript
 * const [users, { mutate, refetch }] = createAllUsersResource()
 * ```
 * 
 * @returns ```typescript
 * [Resource<User[]>, { mutate: Setter<User[]>, refetch: () => void }]
 * ```
 */
const useAllUsers = () => {
  return createResource(async () => {
    const allUsersRes = await client.getUsers({})
    return allUsersRes.users
  })
}

export default useAllUsers
