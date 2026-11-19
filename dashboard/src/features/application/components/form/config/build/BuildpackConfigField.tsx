import { Field } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const BuildpackConfigField: Component<Props> = (props) => (
  <Field of={props.form} path={['form', 'config', 'buildConfig', 'value', 'buildpack', 'context']}>
    {(field) => (
      <TextField
        label="Context"
        info={{
          content: (
            <>
              <div>ビルド対象ディレクトリ</div>
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
)

export default BuildpackConfigField
