import { safeParse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { envVarSchema } from './envVarSchema'

const validator = (input: unknown) => safeParse(envVarSchema, input)

describe('Envvar Schema', () => {
  test('ok: valid input', () => {
    expect(
      validator({
        variables: [
          { key: 'key1', value: 'value1' },
          { key: 'key2', value: 'value2' },
          { key: '', value: '' },
        ],
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty key', () => {
    expect(
      validator({
        variables: [
          { key: '', value: 'value1' },
          { key: 'key2', value: 'value2' },
        ],
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Please enter a key',
        path: [
          expect.objectContaining({
            key: 'variables',
          }),
          expect.objectContaining({
            key: 0,
          }),
          expect.objectContaining({
            key: 'key',
          }),
        ],
      }),
    ])
  })

  test('ng: duplicate key', () => {
    expect(
      validator({
        variables: [
          { key: 'key1', value: 'value1' },
          { key: 'key1', value: 'value2' },
        ],
      }).issues,
    ).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          message: '同じキーの環境変数が存在します',
          path: [
            expect.objectContaining({
              key: 'variables',
            }),
            expect.objectContaining({
              key: 0,
            }),
            expect.objectContaining({
              key: 'key',
            }),
          ],
        }),
        expect.objectContaining({
          message: '同じキーの環境変数が存在します',
          path: [
            expect.objectContaining({
              key: 'variables',
            }),
            expect.objectContaining({
              key: 1,
            }),
            expect.objectContaining({
              key: 'key',
            }),
          ],
        }),
      ]),
    )
  })
})
