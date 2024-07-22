import { safeParse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { createOrUpdateRepositorySchema } from './repositorySchema'

const validator = (input: unknown) => safeParse(createOrUpdateRepositorySchema, input)

describe('Create Repository Schema', () => {
  const base = {
    type: 'create',
    name: 'test repository',
    url: 'https://example.com/test/test.git',
    auth: {
      method: 'none',
      value: {
        none: {},
      },
    },
  }

  test('ok: valid input (auth none)', () => {
    expect(
      validator({
        ...base,
        auth: {
          method: 'none',
          value: {
            none: {},
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (auth basic)', () => {
    expect(
      validator({
        ...base,
        auth: {
          method: 'basic',
          value: {
            basic: {
              username: 'test name',
              password: 'test password',
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (auth ssh)', () => {
    expect(
      validator({
        ...base,
        auth: {
          method: 'ssh',
          value: {
            ssh: {
              keyId: 'test key id',
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty name', () => {
    expect(
      validator({
        ...base,
        name: '',
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Repository Name',
        path: [
          expect.objectContaining({
            key: 'name',
          }),
        ],
      }),
    ])
  })

  test('ng: empty url', () => {
    expect(
      validator({
        ...base,
        url: '',
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Repository URL',
        path: [
          expect.objectContaining({
            key: 'url',
          }),
        ],
      }),
    ])
  })

  test("ng: auth method is basic, but the URL starts with 'http://'", () => {
    expect(
      validator({
        ...base,
        url: 'http://example.com/test/test.git',
        auth: {
          method: 'basic',
          value: {
            basic: {
              username: 'test name',
              password: 'test password',
            },
          },
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
        path: [
          expect.objectContaining({
            key: 'url',
          }),
        ],
      }),
    ])
  })
})

describe('Update Repository Schema', () => {
  const base = {
    id: 'testRepositoryId',
    type: 'update',
    name: 'test repository',
    url: 'https://example.com/test/test.git',
    auth: {
      method: 'none',
      value: {
        none: {},
      },
    },
    ownerIds: ['owner1'],
  }

  test('ok: valid input', () => {
    expect(validator(base)).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty id', () => {
    expect(
      validator({
        ...base,
        id: undefined,
      }).issues,
    ).toEqual([
      expect.objectContaining({
        path: [
          expect.objectContaining({
            key: 'id',
          }),
        ],
      }),
    ])
  })
})
