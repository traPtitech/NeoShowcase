import { safeParse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { applicationGroupSchema } from './applicationGroupSchema'

const validator = (input: unknown) => safeParse(applicationGroupSchema, input)

describe('Application Group Schema', () => {
  test('ok: name with exactly 100 characters', () => {
    expect(validator({ name: 'あ'.repeat(100) })).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: name with 100 surrogate-pair characters (200 UTF-16 code units)', () => {
    expect(validator({ name: '𠮷'.repeat(100) })).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: name with 101 characters', () => {
    expect(validator({ name: 'あ'.repeat(101) })).toEqual(expect.objectContaining({ success: false }))
  })

  test('ng: empty name', () => {
    expect(validator({ name: '' })).toEqual(expect.objectContaining({ success: false }))
  })
})
