import { create } from '@bufbuild/protobuf'
import { timestampFromDate } from '@bufbuild/protobuf/wkt'
import { describe, expect, test } from 'vitest'
import { ApplicationSchema } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { groupApplications } from './application'

const app = (id: string, groupId: string | undefined, t: number) =>
  create(ApplicationSchema, { id, groupId, updatedAt: timestampFromDate(new Date(t)) })

const shape = (apps: Parameters<typeof groupApplications>[0], sort: 'asc' | 'desc') =>
  groupApplications(apps, sort).map((e) => ({ groupId: e.groupId, ids: e.apps.map((a) => a.id) }))

describe('groupApplications', () => {
  const apps = [app('g1', 'G', 1000), app('s', undefined, 2000), app('g2', 'G', 3000)]

  test('desc: newest entry first, newest app first within an entry', () => {
    expect(shape(apps, 'desc')).toEqual([
      { groupId: 'G', ids: ['g2', 'g1'] },
      { groupId: undefined, ids: ['s'] },
    ])
  })

  test('asc: oldest entry first, oldest app first within an entry', () => {
    expect(shape(apps, 'asc')).toEqual([
      { groupId: undefined, ids: ['s'] },
      { groupId: 'G', ids: ['g1', 'g2'] },
    ])
  })

  test('standalone apps next to each other in the order share one entry', () => {
    const mixed = [
      app('s1', undefined, 5000),
      app('s2', undefined, 4000),
      app('g1', 'G', 3000),
      app('s3', undefined, 1000),
    ]
    expect(shape(mixed, 'desc')).toEqual([
      { groupId: undefined, ids: ['s1', 's2'] },
      { groupId: 'G', ids: ['g1'] },
      { groupId: undefined, ids: ['s3'] },
    ])
  })
})
