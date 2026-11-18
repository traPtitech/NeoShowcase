import { Field, getInput } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { TextField } from '/@/components/UI/TextField'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const StaticConfigField: Component<Props> = (props) => {
  const buildType = () => getInput(props.form, { path: ['form', 'config', 'buildConfig', 'type'] })

  return (
    <>
      <Field of={props.form} path={['form', 'config', 'deployConfig', 'value', 'static', 'artifactPath']}>
        {(field) => (
          <TextField
            label="Artifact Path"
            required
            info={{
              content: (
                <>
                  <div>静的ファイルが生成されるディレクトリ</div>
                  <div>({buildType() === 'cmd' ? 'リポジトリルート' : 'Context'}からの相対パス)</div>
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
      <Field of={props.form} path={['form', 'config', 'deployConfig', 'value', 'static', 'spa']}>
        {(field) => (
          <RadioGroup
            label="Is SPA (Single Page Application)"
            info={{
              content: (
                <>
                  <div>配信するファイルがSPAである</div>
                  <div>(いい感じのフォールバック設定が付きます)</div>
                </>
              ),
            }}
            {...field.props}
            options={[
              { value: 'true', label: 'Yes' },
              { value: 'false', label: 'No' },
            ]}
            value={field.input ?? 'false'}
            setValue={(v) => field.onInput(v)}
            readOnly={props.readonly}
          />
        )}
      </Field>
    </>
  )
}

export default StaticConfigField
