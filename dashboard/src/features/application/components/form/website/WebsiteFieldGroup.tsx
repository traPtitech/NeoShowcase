import { Collapsible } from '@kobalte/core/collapsible'
import { Field, getValue, remove, setValue } from '@modular-forms/solid'
import { type Component, Show, createEffect, createMemo } from 'solid-js'
import { AuthenticationType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { systemInfo } from '/@/libs/api'
import { clsx } from '/@/libs/clsx'
import { useApplicationForm } from '../../../provider/applicationFormProvider'
import UrlField from './UrlField'

const authenticationTypeOptions: RadioOption<`${AuthenticationType}`>[] = [
  { value: `${AuthenticationType.OFF}`, label: 'OFF' },
  { value: `${AuthenticationType.SOFT}`, label: 'SOFT' },
  { value: `${AuthenticationType.HARD}`, label: 'HARD' },
]

type Props = {
  index: number
  isRuntimeApp: boolean
  readonly?: boolean
}

const WebsiteFieldGroup: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  const availableDomains = systemInfo()?.domains ?? []
  const selectedDomain = createMemo(() => {
    const domainString = getValue(formStore, `form.websites.${props.index}.domain`)
    return availableDomains.find((d) => d.domain === domainString)
  })
  const authAvailable = (): boolean => selectedDomain()?.authAvailable ?? false

  createEffect(() => {
    if (!authAvailable()) {
      setValue(formStore, `form.websites.${props.index}.authentication`, `${AuthenticationType.OFF}`)
    }
  })

  const handleDelete = () => {
    remove(formStore, 'form.websites', { at: props.index })
    close()
  }

  return (
    <>
      <div class="relative flex w-full flex-col gap-6">
        <div class="-top-2 -right-2 absolute">
          <Button onClick={handleDelete} variants="textError" size="small" type="button">
            Delete
          </Button>
        </div>
        <UrlField index={props.index} readonly={props.readonly} showHttpPort={props.isRuntimeApp} />
        {/* Field componentがmountされないとそのfieldがformに登録されないためforceMountする */}
        <Collapsible forceMount>
          <Collapsible.Trigger class="flex w-full cursor-pointer appearance-none items-center gap-2 border-none bg-inherit text-medium text-text-black">
            詳細
            <div class="[[data-expanded]_&]:transform-rotate-180deg size-6 transition-transform duration-200">
              <span class="i-material-symbols:expand-more text-2xl/6" />
            </div>
          </Collapsible.Trigger>
          <Collapsible.Content
            class={clsx(
              'flex h-0 gap-2 overflow-hidden',
              '[data-[expanded]_&]:h-[var(--kb-collapsible-content-height)] [data-[expanded]_&]:overflow-visible',
            )}
          >
            <div class="flex w-full flex-col gap-2 pt-2">
              <Field of={formStore} name={`form.websites.${props.index}.authentication`}>
                {(field, fieldProps) => (
                  <RadioGroup<`${AuthenticationType}`>
                    label="部員認証"
                    info={{
                      style: 'left',
                      props: {
                        content: (
                          <>
                            <div>OFF: 誰でもアクセス可能</div>
                            <div>SOFT: 部員の場合X-Forwarded-Userをセット</div>
                            <div>HARD: 部員のみアクセス可能</div>
                          </>
                        ),
                      },
                    }}
                    {...fieldProps}
                    tooltip={{
                      props: {
                        content: `${selectedDomain()?.domain}では部員認証が使用できません`,
                      },
                      disabled: authAvailable(),
                    }}
                    options={authenticationTypeOptions}
                    value={field.value ?? `${AuthenticationType.OFF}`}
                    disabled={!authAvailable()}
                    readOnly={props.readonly}
                  />
                )}
              </Field>
              <FormItem title="高度な設定">
                <Field of={formStore} name={`form.websites.${props.index}.stripPrefix`} type="boolean">
                  {(field, fieldProps) => (
                    <CheckBox.Option
                      {...fieldProps}
                      label="Strip Path Prefix"
                      checked={field.value ?? false}
                      readOnly={props.readonly}
                    />
                  )}
                </Field>
                <Show when={props.isRuntimeApp}>
                  <Field of={formStore} name={`form.websites.${props.index}.h2c`} type="boolean">
                    {(field, fieldProps) => (
                      <CheckBox.Option
                        {...fieldProps}
                        label="Use h2c"
                        checked={field.value ?? false}
                        readOnly={props.readonly}
                      />
                    )}
                  </Field>
                </Show>
              </FormItem>
            </div>
          </Collapsible.Content>
        </Collapsible>
      </div>
    </>
  )
}

export default WebsiteFieldGroup
