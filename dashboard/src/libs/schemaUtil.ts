import * as v from 'valibot'

// KobalteのRadioGroupでは値としてbooleanが使えずstringしか使えないため、
// RadioGroupでboolean入力を受け取りたい場合はこれを使用する
export const stringBooleanSchema = v.pipe(
  v.union([v.literal('true'), v.literal('false')]),
  v.transform((i) => i === 'true'),
)
