import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'

type Props = {
  readonly?: boolean
}

const BuildpackConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  return (
    <Field of={formStore} name="form.config.buildConfig.value.buildpack.context">
      {(field, fieldProps) => (
        <TextField
          label="Context"
          info={{
            props: {
              content: (
                <>
                  <div>ビルド対象ディレクトリ</div>
                  <div>(リポジトリルートからの相対パス)</div>
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
  )
}

export default BuildpackConfigField
