import { Field, getInput, remove, setInput } from '@formisch/solid/v2'
import { type Component, createEffect, createMemo, Show } from 'solid-js'
import { AuthenticationType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { Button } from '/@/components/UI/Button'
import { systemInfo } from '/@/libs/session'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'
import UrlField from './UrlField'

const authenticationTypeOptions: RadioOption<`${AuthenticationType}`>[] = [
  { value: `${AuthenticationType.OFF}`, label: 'OFF' },
  { value: `${AuthenticationType.SOFT}`, label: 'SOFT' },
  { value: `${AuthenticationType.HARD}`, label: 'HARD' },
]

type Props = {
  form: ApplicationFormStore
  index: number
  isRuntimeApp: boolean
  readonly?: boolean
}

const WebsiteFieldGroup: Component<Props> = (props) => {
  const selectedDomain = createMemo(() => {
    const domainString = getInput(props.form, { path: ['form', 'websites', props.index, 'domain'] })
    return systemInfo().domains.find((d) => d.domain === domainString)
  })
  const authAvailable = () => selectedDomain()?.authAvailable ?? false

  // 部員認証が使えないドメインが選ばれたら認証をOFFに戻す
  createEffect(
    () => selectedDomain() !== undefined && !authAvailable(),
    (authUnavailable) => {
      if (authUnavailable) {
        setInput(props.form, {
          path: ['form', 'websites', props.index, 'authentication'],
          input: `${AuthenticationType.OFF}`,
        })
      }
    },
  )

  return (
    <div class="relative flex w-full flex-col gap-6">
      <div class="absolute -top-2 -right-2">
        <Button
          onClick={() => remove(props.form, { path: ['form', 'websites'], at: props.index })}
          variants="textError"
          size="small"
          type="button"
        >
          Delete
        </Button>
      </div>
      <UrlField form={props.form} index={props.index} readonly={props.readonly} showHttpPort={props.isRuntimeApp} />
      <details class="group">
        <summary class="text-medium text-text-black flex w-full cursor-pointer list-none items-center gap-2 [&::-webkit-details-marker]:hidden">
          詳細
          <div class="i-material-symbols:expand-more size-6 shrink-0 text-2xl/6 transition-transform duration-200 group-open:rotate-180" />
        </summary>
        <div class="flex w-full flex-col gap-2 pt-2">
          <Field of={props.form} path={['form', 'websites', props.index, 'authentication']}>
            {(field) => (
              <RadioGroup<`${AuthenticationType}`>
                label="部員認証"
                info={{
                  style: 'left',
                  content: (
                    <>
                      <div>OFF: 誰でもアクセス可能</div>
                      <div>SOFT: 部員の場合X-Forwarded-Userをセット</div>
                      <div>HARD: 部員のみアクセス可能</div>
                    </>
                  ),
                }}
                {...field.props}
                tooltip={{
                  content: `${selectedDomain()?.domain}では部員認証が使用できません`,
                  disabled: authAvailable(),
                }}
                options={authenticationTypeOptions}
                value={field.input ?? `${AuthenticationType.OFF}`}
                setValue={(v) => field.onInput(v)}
                disabled={!authAvailable()}
                readOnly={props.readonly}
              />
            )}
          </Field>
          <FormItem title="高度な設定">
            <Field of={props.form} path={['form', 'websites', props.index, 'stripPrefix']}>
              {(field) => (
                <CheckBox.Option
                  {...field.props}
                  label="Strip Path Prefix"
                  checked={field.input ?? false}
                  setChecked={(v) => field.onInput(v)}
                  readOnly={props.readonly}
                />
              )}
            </Field>
            <Show when={props.isRuntimeApp}>
              <Field of={props.form} path={['form', 'websites', props.index, 'h2c']}>
                {(field) => (
                  <CheckBox.Option
                    {...field.props}
                    label="Use h2c"
                    checked={field.input ?? false}
                    setChecked={(v) => field.onInput(v)}
                    readOnly={props.readonly}
                  />
                )}
              </Field>
            </Show>
          </FormItem>
        </div>
      </details>
    </div>
  )
}

export default WebsiteFieldGroup
