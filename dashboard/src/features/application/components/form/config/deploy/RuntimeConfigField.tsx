import { Field, getInput, setInput } from '@formisch/solid/v2'
import { type Component, createSignal, Show, untrack } from 'solid-js'
import { AutoShutdownConfig_StartupBehavior } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { TextField } from '/@/components/UI/TextField'
import { ToolTip } from '/@/components/UI/ToolTip'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
  disableEditDB?: boolean
}

const runtimePath = ['form', 'config', 'deployConfig', 'value', 'runtime'] as const

const RuntimeConfigField: Component<Props> = (props) => {
  const buildType = () => getInput(props.form, { path: ['form', 'config', 'buildConfig', 'type'] })
  const autoShutdown = () => getInput(props.form, { path: [...runtimePath, 'autoShutdown', 'enabled'] })

  // 「Use Database」はチェックボックスの値から導出できない (Yes を選んだ直後は両方 false) ため、初期値だけを導出して別途保持する
  const [useDB, setUseDB] = createSignal(
    untrack(
      () =>
        !!getInput(props.form, { path: [...runtimePath, 'useMariadb'] }) ||
        !!getInput(props.form, { path: [...runtimePath, 'useMongodb'] }),
    ),
  )

  const updateUseDB = (value: boolean) => {
    setUseDB(value)
    if (!value) {
      setInput(props.form, { path: [...runtimePath, 'useMariadb'], input: false })
      setInput(props.form, { path: [...runtimePath, 'useMongodb'], input: false })
    }
  }

  const EntryPointField = () => (
    <Field of={props.form} path={[...runtimePath, 'entrypoint']}>
      {(field) => (
        <TextField
          label="Entrypoint"
          info={{ content: buildType() === 'cmd' ? 'アプリ起動コマンド' : 'コンテナのEntrypoint上書き' }}
          {...field.props}
          value={field.input ?? ''}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
  )

  const CommandOverrideField = () => (
    <Field of={props.form} path={[...runtimePath, 'command']}>
      {(field) => (
        <TextField
          label="Command"
          info={{ content: 'コンテナのCommand上書き' }}
          {...field.props}
          value={field.input ?? ''}
          error={field.errors?.[0]}
          readOnly={props.readonly}
        />
      )}
    </Field>
  )

  return (
    <>
      <FormItem
        title="Use Database"
        helpText="アプリ作成後は変更できません"
        tooltip={{ content: <div>データーベースを使用する場合はチェック</div> }}
      >
        <RadioGroup
          tooltip={{ content: 'アプリ作成後は変更できません', disabled: !props.disableEditDB }}
          options={[
            { value: 'true', label: 'Yes' },
            { value: 'false', label: 'No' },
          ]}
          value={useDB() ? 'true' : 'false'}
          setValue={(v) => updateUseDB(v === 'true')}
          disabled={props.disableEditDB}
        />
      </FormItem>
      <Show when={useDB()}>
        <FormItem title="Database">
          <ToolTip content="アプリ作成後は変更できません" disabled={!props.disableEditDB}>
            <CheckBox.Container>
              <Field of={props.form} path={[...runtimePath, 'useMariadb']}>
                {(field) => (
                  <CheckBox.Option
                    {...field.props}
                    label="MariaDB"
                    checked={field.input ?? false}
                    setChecked={(v) => field.onInput(v)}
                    disabled={props.disableEditDB}
                  />
                )}
              </Field>
              <Field of={props.form} path={[...runtimePath, 'useMongodb']}>
                {(field) => (
                  <CheckBox.Option
                    {...field.props}
                    label="MongoDB"
                    checked={field.input ?? false}
                    setChecked={(v) => field.onInput(v)}
                    disabled={props.disableEditDB}
                  />
                )}
              </Field>
            </CheckBox.Container>
          </ToolTip>
        </FormItem>
      </Show>
      <Show when={buildType() === 'cmd'}>
        <EntryPointField />
      </Show>
      <FormItem title="高度な設定">
        <Show when={buildType() !== 'cmd'}>
          <EntryPointField />
        </Show>
        <CommandOverrideField />
      </FormItem>
      <FormItem
        title="Auto Shutdown"
        tooltip={{ content: '一定期間アクセスがない場合にアプリを自動でシャットダウンします' }}
      >
        <CheckBox.Container>
          <Field of={props.form} path={[...runtimePath, 'autoShutdown', 'enabled']}>
            {(field) => (
              <CheckBox.Option
                {...field.props}
                label="自動シャットダウン"
                checked={field.input ?? false}
                setChecked={(v) => field.onInput(v)}
                error={field.errors?.[0]}
              />
            )}
          </Field>
        </CheckBox.Container>
      </FormItem>
      <Show when={autoShutdown()}>
        <FormItem title="Startup Behavior" tooltip={{ content: '起動時の挙動' }}>
          <Field of={props.form} path={[...runtimePath, 'autoShutdown', 'startup']}>
            {(field) => (
              <RadioGroup
                wrap={false}
                full
                options={[
                  {
                    value: `${AutoShutdownConfig_StartupBehavior.LOADING_PAGE}`,
                    label: 'Loading Page',
                    description: 'アプリ起動時にローディングページを表示します。Webアプリ向け',
                  },
                  {
                    value: `${AutoShutdownConfig_StartupBehavior.BLOCKING}`,
                    label: 'Blocking',
                    description: 'アクセスに対し、アプリが起動するまでリクエストを待機させます。APIサーバー向け',
                  },
                ]}
                required
                {...field.props}
                value={field.input}
                setValue={(v) => field.onInput(v)}
                error={field.errors?.[0]}
              />
            )}
          </Field>
        </FormItem>
      </Show>
    </>
  )
}

export default RuntimeConfigField
