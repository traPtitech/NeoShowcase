import { Field } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const CmdConfigField: Component<Props> = (props) => (
  <>
    <Field of={props.form} path={['form', 'config', 'buildConfig', 'value', 'cmd', 'baseImage']}>
      {(field) => (
        <TextField
          label="Base Image"
          info={{
            content: (
              <>
                <div>ベースとなるDocker Image</div>
                <div>「イメージ名:タグ名」の形式</div>
                <div>ビルドが必要無い場合は空</div>
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
    <Field of={props.form} path={['form', 'config', 'buildConfig', 'value', 'cmd', 'buildCmd']}>
      {(field) => (
        <TextField
          label="Build Command"
          info={{
            content: (
              <>
                <div>イメージ上でビルド時に実行するコマンド</div>
                <div>リポジトリルートで実行されます</div>
              </>
            ),
          }}
          {...field.props}
          multiline
          value={field.input ?? ''}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
  </>
)

export default CmdConfigField
