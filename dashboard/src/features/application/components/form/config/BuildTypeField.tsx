import { Field, getValues } from '@modular-forms/solid'
import { type Component, Show } from 'solid-js'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

type Props = {
  readonly?: boolean
}

const BuildTypeField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const runType = () => getValues(formStore).form?.config?.deployConfig?.type

  return (
    <>
      <Field of={formStore} name="form.config.deployConfig.type">
        {(field, fieldProps) => (
          <RadioGroup
            label="Deploy Type"
            wrap={false}
            full
            required
            {...fieldProps}
            options={[
              {
                value: 'runtime',
                label: 'Runtime',
                description:
                  'コマンドを実行してアプリを起動します。サーバープロセスやバックグラウンド処理がある場合、こちらを選びます。',
              },
              {
                value: 'static',
                label: 'Static',
                description: '静的ファイルを配信します。ビルド（任意）を実行できます。',
              },
            ]}
            value={field.value}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Show when={runType() !== undefined}>
        <Field of={formStore} name="form.config.buildConfig.type">
          {(field, fieldProps) => (
            <RadioGroup
              label="Build Type"
              wrap={false}
              full
              required
              {...fieldProps}
              options={[
                {
                  value: 'buildpack',
                  label: 'Buildpack',
                  description: 'ビルド設定を、リポジトリ内ファイルから自動検出します。（オススメ）',
                },
                {
                  value: 'cmd',
                  label: 'Command',
                  description: 'ベースイメージとビルドコマンド（任意）を設定します。',
                },
                {
                  value: 'dockerfile',
                  label: 'Dockerfile',
                  description: 'リポジトリ内Dockerfileからビルドを行います。',
                },
              ]}
              value={field.value}
              error={field.error}
              readOnly={props.readonly}
            />
          )}
        </Field>
      </Show>
    </>
  )
}

export default BuildTypeField
