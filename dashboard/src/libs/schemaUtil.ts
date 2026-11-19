import * as v from 'valibot'

// KobalteのRadioGroupでは値としてbooleanが使えずstringしか使えないため、
// RadioGroupでboolean入力を受け取りたい場合はこれを使用する
export const stringBooleanSchema = v.pipe(
  v.union([v.literal('true'), v.literal('false')]),
  v.transform((i) => i === 'true'),
)

// <input type="number"> の値は文字列で渡ってくるため、number として保持する form 向けに変換する。
// 空欄は NaN になり、`v.number()` のバリデーションで弾かれる。
export const inputToNumber = (value: string): number => (value === '' ? Number.NaN : Number(value))
export const numberToInput = (value: number | undefined): string =>
  value === undefined || Number.isNaN(value) ? '' : String(value)
