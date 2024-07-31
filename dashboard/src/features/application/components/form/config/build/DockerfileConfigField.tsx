import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'

type Props = {
  readonly?: boolean
}

const DockerfileConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  return (
    <>
      <Field of={formStore} name="form.config.buildConfig.value.dockerfile.context">
        {(field, fieldProps) => (
          <TextField
            label="Context"
            info={{
              props: {
                content: (
                  <>
                    <div>ビルドContext</div>
                    <div>(リポジトリルートからの相対パス)</div>
                  </>
                ),
              },
            }}
            value={field.value ?? ''}
            error={field.error}
            readOnly={props.readonly}
            {...fieldProps}
          />
        )}
      </Field>
      <Field of={formStore} name="form.config.buildConfig.value.dockerfile.dockerfileName">
        {(field, fieldProps) => (
          <TextField
            label="Dockerfile Name"
            required
            info={{
              props: {
                content: (
                  <>
                    <div>Dockerfileへのパス</div>
                    <div>(Contextからの相対パス)</div>
                  </>
                ),
              },
            }}
            value={field.value ?? ''}
            error={field.error}
            readOnly={props.readonly}
            {...fieldProps}
          />
        )}
      </Field>
    </>
  )
}

export default DockerfileConfigField
