import { safeParse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { createOrUpdateRepositorySchema } from './repositorySchema'

const validator = (input: unknown) => safeParse(createOrUpdateRepositorySchema, input)

describe('Create Repository Schema', () => {
  test('ok: valid input (auth none)', () => {
    expect(
      validator({
        type: 'create',
        form: {
          name: 'test repository',
          url: 'https://example.com/test/test.git',
          auth: {
            method: 'none',
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (auth basic)', () => {
    expect(
      validator({
        type: 'create',
        form: {
          name: 'test repository',
          url: 'https://example.com/test/test.git',
          auth: {
            method: 'basic',
            value: {
              basic: {
                username: 'test name',
                password: 'test password',
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (auth ssh)', () => {
    expect(
      validator({
        type: 'create',
        form: {
          name: 'test repository',
          url: 'https://example.com/test/test.git',
          auth: {
            method: 'ssh',
            value: {
              ssh: {
                keyId: 'test key id',
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty name', () => {
    expect(
      validator({
        type: 'create',
        form: {
          name: '',
          url: 'https://example.com/test/test.git',
          auth: {
            method: 'none',
            value: {
              none: {},
            },
          },
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Repository Name',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
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
        type: 'create',
        form: {
          name: 'test repository',
          url: '',
          auth: {
            method: 'none',
            value: {
              none: {},
            },
          },
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Repository URL',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
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
        type: 'create',
        form: {
          name: 'test repository',
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
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
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
    expect(
      validator({
        type: 'update',
        form: base,
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: update name', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: base.id,
          name: base.name,
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: update auth config', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: base.id,
          url: base.url,
          auth: base.auth,
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: update ownerIds', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: base.id,
          ownerIds: base.ownerIds,
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty name', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: base.id,
          name: '',
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'name',
          }),
        ],
      }),
    ])
  })

  test("ng: auth method is basic, but the URL starts with 'http://'", () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: base.id,
          name: base.name,
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
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Basic認証を使用する場合、URLはhttps://から始まる必要があります',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'url',
          }),
        ],
      }),
    ])
  })
})
