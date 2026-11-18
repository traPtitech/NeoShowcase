import { Field, FieldArray, getInput, insert, remove } from '@formisch/solid/v2'
import { type Component, createEffect, For } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { type EnvVarFormStore, emptyEnvVar } from '../../../schema/envVarSchema'

type Props = {
  form: EnvVarFormStore
}

const isBlank = (envVar: { key?: string; value?: string } | undefined) => !envVar?.key && !envVar?.value

const EnvVarConfigField: Component<Props> = (props) => {
  // 末尾に空の行をちょうど1つだけ残し、それ以外の空の行は削除する
  // 書き込みは正規化済みの状態では発生しないため、ループしない
  createEffect(
    () => getInput(props.form, { path: ['variables'] }),
    (variables) => {
      const rows = variables ?? []
      for (let i = rows.length - 2; i >= 0; i--) {
        if (isBlank(rows[i])) remove(props.form, { path: ['variables'], at: i })
      }
      if (rows.length === 0 || !isBlank(rows.at(-1))) {
        insert(props.form, { path: ['variables'], initialInput: emptyEnvVar })
      }
    },
    { defer: true },
  )

  return (
    <div class="gap-col-6 gap-row-2 text-bold text-text-black grid w-full grid-cols-2">
      <div>Key</div>
      <div>Value</div>
      <FieldArray of={props.form} path={['variables']}>
        {(fieldArray) => (
          <For each={fieldArray.items}>
            {(_, index) => {
              const isSystem = () => !!getInput(props.form, { path: ['variables', index(), 'system'] })
              const systemTooltip = () => ({ content: 'システム環境変数は変更できません', disabled: !isSystem() })

              return (
                <>
                  <Field of={props.form} path={['variables', index(), 'key']}>
                    {(field) => (
                      <TextField
                        tooltip={systemTooltip()}
                        {...field.props}
                        value={field.input ?? ''}
                        error={field.errors?.[0]}
                        disabled={isSystem()}
                      />
                    )}
                  </Field>
                  <Field of={props.form} path={['variables', index(), 'value']}>
                    {(field) => (
                      <TextField
                        tooltip={systemTooltip()}
                        {...field.props}
                        value={field.input ?? ''}
                        error={field.errors?.[0]}
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
