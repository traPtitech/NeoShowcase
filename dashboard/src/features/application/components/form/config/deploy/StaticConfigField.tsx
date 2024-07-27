import { Field, getValues } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'

type Props = {
  readonly?: boolean
}

const StaticConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const buildType = () => getValues(formStore).form?.config?.buildConfig?.type

  return (
    <>
      <Field of={formStore} name="form.config.deployConfig.value.static.artifactPath">
        {(field, fieldProps) => (
          <TextField
            label="Artifact Path"
            required
            info={{
              props: {
                content: (
                  <>
                    <div>静的ファイルが生成されるディレクトリ</div>
                    <div>({buildType() === 'cmd' ? 'リポジトリルート' : 'Context'}からの相対パス)</div>
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
      <Field of={formStore} name="form.config.deployConfig.value.static.spa">
        {(field, fieldProps) => (
          <RadioGroup
            label="Is SPA (Single Page Application)"
            info={{
              props: {
                content: (
                  <>
                    <div>配信するファイルがSPAである</div>
                    <div>(いい感じのフォールバック設定が付きます)</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            options={[
              { value: 'true', label: 'Yes' },
              { value: 'false', label: 'No' },
            ]}
            value={field.value ?? 'false'}
            readOnly={props.readonly}
          />
        )}
      </Field>
    </>
  )
}

export default StaticConfigField
