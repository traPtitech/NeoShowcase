import { Field, getValues, setValues } from '@modular-forms/solid'
import { type Component, Show, createEffect, createResource } from 'solid-js'
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

  const AutoShutdownField = () => (
    <Field
      of={formStore}
      name="form.config.deployConfig.value.runtime.autoShutdown"
      // @ts-expect-error: autoShutdown は deployConfig.type === "static" の時存在しないためtsの型の仕様上エラーが出る
      type="boolean"
    >
      {(field, fieldProps) => (
        <FormItem
          title="Auto Shutdown"
          tooltip={{
            props: { content: <div>アプリへのアクセスが一定期間ない場合、自動でアプリをシャットダウンします</div> },
          }}
        >
          <CheckBox.Option
            {...fieldProps}
            label="自動シャットダウン"
            checked={field.value ?? false}
            error={field.error}
          />
        </FormItem>
      )}
    </Field>
  )

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
      <AutoShutdownField />
      <Show when={buildType() === 'cmd'}>
        <EntryPointField />
      </Show>
      <FormItem title="高度な設定">
        <Show when={buildType() !== 'cmd'}>
          <EntryPointField />
        </Show>
        <CommandOverrideField />
      </FormItem>
    </>
  )
}

export default RuntimeConfigField
