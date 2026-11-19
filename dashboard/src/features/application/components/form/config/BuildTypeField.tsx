import { Field, getInput } from '@formisch/solid/v2'
import { type Component, Show } from 'solid-js'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const BuildTypeField: Component<Props> = (props) => (
  <>
    <Field of={props.form} path={['form', 'config', 'deployConfig', 'type']}>
      {(field) => (
        <RadioGroup
          label="Deploy Type"
          wrap={false}
          full
          required
          {...field.props}
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
          value={field.input}
          setValue={(v) => field.onInput(v)}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
    <Show when={getInput(props.form, { path: ['form', 'config', 'deployConfig', 'type'] }) !== undefined}>
      <Field of={props.form} path={['form', 'config', 'buildConfig', 'type']}>
        {(field) => (
          <RadioGroup
            label="Build Type"
            wrap={false}
            full
            required
            {...field.props}
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
            value={field.input}
            setValue={(v) => field.onInput(v)}
            error={field.errors?.[0]}
            readOnly={props.readonly}
          />
        )}
      </Field>
    </Show>
  </>
)

export default BuildTypeField
