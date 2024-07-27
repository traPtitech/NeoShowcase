import { Field, getValues } from '@modular-forms/solid'
import { type Component, Show, createSignal } from 'solid-js'
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
  const [useDB, setUseDB] = createSignal(false)

  const buildType = () => getValues(formStore).form?.config?.buildConfig?.type

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
        <Field of={formStore} name="form.config.deployConfig.value.runtime.useMariadb">
          {(field, fieldProps) => (
            <CheckBox.Option
              {...fieldProps}
              label="MariaDB"
              checked={field.value ?? false}
              disabled={props.disableEditDB}
            />
          )}
        </Field>
        <Field of={formStore} name="form.config.deployConfig.value.runtime.useMongodb">
          {(field, fieldProps) => (
            <CheckBox.Option
              {...fieldProps}
              label="MongoDB"
              checked={field.value ?? false}
              disabled={props.disableEditDB}
            />
          )}
        </Field>
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
    </>
  )
}

export default RuntimeConfigField
