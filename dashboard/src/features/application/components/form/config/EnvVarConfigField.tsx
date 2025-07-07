import { Field, FieldArray, getValue, getValues, insert, remove } from '@modular-forms/solid'
import { type Component, createReaction, For, onMount } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useEnvVarConfigForm } from '../../../provider/envVarConfigFormProvider'

const EnvVarConfigField: Component = () => {
  const { formStore } = useEnvVarConfigForm()

  // keyとvalueが空となるenv varを削除し、最後に空のenv varを追加する
  const stripEnvVars = () => {
    const envVars = getValues(formStore).variables ?? []

    for (let i = envVars.length - 1; i >= 0; i--) {
      if (envVars[i]?.key === '' && envVars[i]?.value === '') {
        remove(formStore, 'variables', { at: i })
      }
    }

    // add empty env var
    insert(formStore, 'variables', {
      value: { key: '', value: '', system: false },
    })

    // 次にvariablesが変更された時に1度だけ再度stripする
    track(() => getValues(formStore, 'variables'))
  }

  const track = createReaction(() => {
    stripEnvVars()
  })

  onMount(() => {
    stripEnvVars()
  })

  return (
    <div class="grid w-full grid-cols-2 gap-col-6 gap-row-2 text-bold text-text-black">
      <div>Key</div>
      <div>Value</div>
      <FieldArray of={formStore} name="variables">
        {(fieldArray) => (
          <For each={fieldArray.items}>
            {(_, index) => {
              const isSystem = () =>
                getValue(formStore, `variables.${index()}.system`, {
                  shouldActive: false,
                })

              return (
                <>
                  <Field of={formStore} name={`variables.${index()}.key`}>
                    {(field, fieldProps) => (
                      <TextField
                        tooltip={{
                          props: {
                            content: 'システム環境変数は変更できません',
                          },
                          disabled: !isSystem(),
                        }}
                        {...fieldProps}
                        value={field.value ?? ''}
                        error={field.error}
                        disabled={isSystem()}
                      />
                    )}
                  </Field>
                  <Field of={formStore} name={`variables.${index()}.value`}>
                    {(field, fieldProps) => (
                      <TextField
                        tooltip={{
                          props: {
                            content: 'システム環境変数は変更できません',
                          },
                          disabled: !isSystem(),
                        }}
                        {...fieldProps}
                        value={field.value ?? ''}
                        error={field.error}
                        disabled={isSystem()}
                        copyable
                      />
                    )}
                  </Field>
                </>
              )
            }}
          </For>
        )}
      </FieldArray>
    </div>
  )
}

export default EnvVarConfigField
