import * as v from 'valibot'
import type { ApplicationEnvVars } from '/@/api/neoshowcase/protobuf/gateway_pb'

export const envVarSchema = v.object({
  variables: v.pipe(
    v.array(
      v.object({
        key: v.pipe(v.string(), v.toUpperCase()),
        value: v.string(),
        system: v.optional(v.boolean()),
      }),
    ),
    // 最後のkey以外は必須
    // see: https://valibot.dev/api/rawCheck/
    v.rawCheck(({ dataset, addIssue }) => {
      if (dataset.typed) {
        dataset.value.forEach((kv, i) => {
          if (i < dataset.value.length - 1 && kv.key.length === 0) {
            addIssue({
              message: 'Please enter a key',
              path: [
                {
                  type: 'array',
                  origin: 'value',
                  input: dataset.value,
                  key: i,
                  value: kv,
                },
                {
                  type: 'object',
                  origin: 'value',
                  input: dataset.value[i],
                  key: 'key',
                  value: kv.key,
                },
              ],
            })
          }
        })
      }
    }),
    // keyの重複チェック
    v.rawCheck(({ dataset, addIssue }) => {
      if (dataset.typed) {
        dataset.value.forEach((kv, i) => {
          const isDuplicate = dataset.value.some((other, j) => i !== j && other.key === kv.key)

          if (isDuplicate) {
            addIssue({
              message: '同じキーの環境変数が存在します',
              path: [
                {
                  type: 'array',
                  origin: 'value',
                  input: dataset.value,
                  key: i,
                  value: kv,
                },
                {
                  type: 'object',
                  origin: 'value',
                  input: dataset.value[i],
                  key: 'key',
                  value: kv.key,
                },
              ],
            })
          }
        })
      }
    }),
  ),
})

export type EnvVarInput = v.InferInput<typeof envVarSchema>

export const parseCreateWebsiteInput = (input: unknown) => {
  const result = v.parse(envVarSchema, input)
  return result
}

export const envVarsMessageToSchema = (envVars: ApplicationEnvVars): EnvVarInput => {
  return {
    variables: envVars.variables.map((variable) => ({
      key: variable.key,
      value: variable.value,
      system: variable.system,
    })),
  }
}

export const handleSubmitEnvVarForm = (
  input: EnvVarInput,
  handler: (output: v.InferOutput<typeof envVarSchema>) => Promise<unknown>,
) => {
  const result = v.parse(envVarSchema, input)
  return handler(result)
}
