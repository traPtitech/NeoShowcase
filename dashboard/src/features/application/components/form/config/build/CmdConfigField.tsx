import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'

type Props = {
  readonly?: boolean
}

const CmdConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  return (
    <>
      <Field of={formStore} name="form.config.buildConfig.value.cmd.baseImage">
        {(field, fieldProps) => (
          <TextField
            label="Base Image"
            info={{
              props: {
                content: (
                  <>
                    <div>ベースとなるDocker Image</div>
                    <div>「イメージ名:タグ名」の形式</div>
                    <div>ビルドが必要無い場合は空</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Field of={formStore} name="form.config.buildConfig.value.cmd.buildCmd">
        {(field, fieldProps) => (
          <TextField
            label="Build Command"
            info={{
              props: {
                content: (
                  <>
                    <div>イメージ上でビルド時に実行するコマンド</div>
                    <div>リポジトリルートで実行されます</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            multiline
            value={field.value ?? ''}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
    </>
  )
}

export default CmdConfigField
