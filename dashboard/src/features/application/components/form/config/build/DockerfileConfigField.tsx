import { Field } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const DockerfileConfigField: Component<Props> = (props) => (
  <>
    <Field of={props.form} path={['form', 'config', 'buildConfig', 'value', 'dockerfile', 'context']}>
      {(field) => (
        <TextField
          label="Context"
          info={{
            content: (
              <>
                <div>ビルドContext</div>
                <div>(リポジトリルートからの相対パス)</div>
              </>
            ),
          }}
          {...field.props}
          value={field.input ?? ''}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
    <Field of={props.form} path={['form', 'config', 'buildConfig', 'value', 'dockerfile', 'dockerfileName']}>
      {(field) => (
        <TextField
          label="Dockerfile Name"
          required
          info={{
            content: (
              <>
                <div>Dockerfileへのパス</div>
                <div>(Contextからの相対パス)</div>
              </>
            ),
          }}
          {...field.props}
          value={field.input ?? ''}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
  </>
)

export default DockerfileConfigField
