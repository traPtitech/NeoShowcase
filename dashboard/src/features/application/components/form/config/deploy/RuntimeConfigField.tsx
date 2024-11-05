import { Field, getValues, setValues } from '@modular-forms/solid'
import { type Component, Show, createEffect, createResource } from 'solid-js'
import { AutoShutdownConfig_StartupBehavior } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { TextField } from '/@/components/UI/TextField'
import { ToolTip } from '/@/components/UI/ToolTip'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'

type Props = {
  readonly?: boolean
  disableEditDB?: boolean
}

const RuntimeConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const [useDB, { mutate: setUseDB }] = createResource(
    () =>
      getValues(formStore, {
        shouldActive: false,
      }).form?.config?.deployConfig,
    (config) => {
      if (!config?.type || config?.type === 'static') {
        return false
      }
      // @ts-expect-error: getValuesの結果のpropertyはすべてMaybeになるためnarrowingが正しく行われない
      return config?.value?.runtime?.useMariadb || config?.value?.runtime?.useMongodb
    },
  )

  const buildType = () => getValues(formStore).form?.config?.buildConfig?.type

  createEffect(() => {
    if (!useDB()) {
      setValues(formStore, {
        form: {
          config: {
            deployConfig: {
              value: {
                runtime: {
                  useMariadb: false,
                  useMongodb: false,
                },
              },
            },
          },
        },
      })
    }
  })

  // @ts-expect-error: autoShutdown は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
  const autoShutdown = () => getValues(formStore).form?.config?.deployConfig?.value?.runtime?.autoShutdown?.enabled

  const setStartupBehavior = (value: string | undefined) => {
    setValues(formStore, {
      form: {
        config: {
          deployConfig: {
            value: {
              // @ts-expect-error: deployConfig は form.type === "static" の時存在しないためtsの型の仕様上エラーが出る
              runtime: {
                autoShutdown: {
                  startup: value,
                },
              },
            },
          },
        },
      },
    })
  }

  const EntryPointField = () => (
    <Field of={formStore} name="form.config.deployConfig.value.runtime.entrypoint">
      {(field, fieldProps) => (
        <TextField
          label="Entrypoint"
          info={{
            props: {
              content: buildType() === 'cmd' ? 'アプリ起動コマンド' : 'コンテナのEntrypoint上書き',
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

  const CommandOverrideField = () => (
    <Field of={formStore} name="form.config.deployConfig.value.runtime.command">
      {(field, fieldProps) => (
        <TextField
          label="Command"
          info={{
            props: {
              content: 'コンテナのCommand上書き',
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

  return (
    <>
      <ToolTip>
        <FormItem
          title="Use Database"
          helpText="アプリ作成後は変更できません"
          tooltip={{
            props: {
              content: (
                <>
                  <div>データーベースを使用する場合はチェック</div>
                </>
              ),
            },
          }}
        >
          <RadioGroup
            tooltip={{
              props: {
                content: <>アプリ作成後は変更できません</>,
              },
              disabled: !props.disableEditDB,
            }}
            options={[
              { value: 'true', label: 'Yes' },
              { value: 'false', label: 'No' },
            ]}
            value={useDB() ? 'true' : 'false'}
            setValue={(v) => setUseDB(v === 'true')}
            disabled={props.disableEditDB}
          />
        </FormItem>
      </ToolTip>
      <Show when={useDB()}>
        <FormItem title="Database">
          <ToolTip
            props={{
              content: <>アプリ作成後は変更できません</>,
            }}
            disabled={!props.disableEditDB}
          >
            <CheckBox.Container>
              <Field
                of={formStore}
                name="form.config.deployConfig.value.runtime.useMariadb"
                // @ts-expect-error: useMariadb は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
                type="boolean"
              >
                {(field, fieldProps) => (
                  <CheckBox.Option
                    {...fieldProps}
                    label="MariaDB"
                    checked={field.value ?? false}
                    disabled={props.disableEditDB}
                  />
                )}
              </Field>
              <Field
                of={formStore}
                name="form.config.deployConfig.value.runtime.useMongodb"
                // @ts-expect-error: useMongodb は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
                type="boolean"
              >
                {(field, fieldProps) => (
                  <CheckBox.Option
                    {...fieldProps}
                    label="MongoDB"
                    checked={field.value ?? false}
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
        tooltip={{
          props: {
            content: '一定期間アクセスがない場合にアプリを自動でシャットダウンします',
          },
        }}
      >
        <CheckBox.Container>
          <Field
            of={formStore}
            name="form.config.deployConfig.value.runtime.autoShutdown.enabled"
            // @ts-expect-error: autoShutdown は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
            type="boolean"
          >
            {(field, fieldProps) => (
              <CheckBox.Option
                {...fieldProps}
                label="自動シャットダウン"
                checked={field.value ?? false}
                error={field.error}
              />
            )}
          </Field>
        </CheckBox.Container>
      </FormItem>
      <Show when={autoShutdown()}>
        <FormItem
          title="Startup Behavior"
          tooltip={{
            props: {
              content: '起動時の挙動',
            },
          }}
        >
          <Field
            of={formStore}
            name="form.config.deployConfig.value.runtime.autoShutdown.startup"
            // @ts-expect-error: autoShutdown は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
            type="string"
          >
            {(field, fieldProps) => (
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
                value={field.value}
                setValue={setStartupBehavior}
                required={true}
                {...fieldProps}
                error={field.error}
              />
            )}
          </Field>
        </FormItem>
      </Show>
    </>
  )
}

export default RuntimeConfigField
