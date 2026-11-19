import { revalidate } from '@solidjs/router'
import { createMemo, refresh } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client } from '/@/libs/api'
import { createSingleton } from '/@/libs/singleton'

/**
 * Session-wide data that does not depend on the current route.
 *
 * These are async memos: reading them before the first response suspends the
 * nearest `<Loading>` boundary. They are app-lifetime singletons (this dashboard is a client-only SPA).
 */
const session = createSingleton(() => {
  const me = createMemo(() => client.getMe({}))
  const systemInfo = createMemo(() => client.getSystemInfo({}))
  const availableMetrics = createMemo(() => client.getAvailableMetrics({}).then((r) => r.metricsNames), {
    lazy: true,
  })
  const users = createMemo(() => client.getUsers({}).then((r) => r.users), { lazy: true })
  const usersById = createMemo(() => new Map(users().map((u) => [u.id, u])), { lazy: true })
  return { me, systemInfo, availableMetrics, users, usersById }
})

export const { me, systemInfo, availableMetrics, users } = session

export const refreshUsers = () => refresh(session.users)

/** Resolves user IDs to users. Unknown IDs (e.g. deleted users) are skipped. */
export const usersFromIds = (ids: string[]): User[] =>
  ids.map((id) => session.usersById().get(id)).filter((u): u is User => u !== undefined)

/**
 * Re-asks every piece of data a boundary may be showing: the router's cached queries (`revalidate()`) and the
 * session memos that live for the lifetime of the module. Retrying goes through this rather than through a list
 * supplied by each boundary, so that adding a boundary cannot silently leave its Retry doing nothing.
 *
 * Memos created inside a boundary are not included; resetting the boundary recreates them.
 *
 * It settles rather than rejects, so that one source still being down cannot stop a caller from resetting a
 * boundary whose own data has recovered. Whatever is still failing throws again when the boundary reads it.
 */
export const retryAll = async () => {
  revalidate()
  await Promise.allSettled([
    refresh(session.me),
    refresh(session.systemInfo),
    refresh(session.availableMetrics),
    refresh(session.users),
  ])
}
