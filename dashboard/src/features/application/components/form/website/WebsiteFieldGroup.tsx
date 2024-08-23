import { Collapsible } from '@kobalte/core/collapsible'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Field, getValue, remove, setValue } from '@modular-forms/solid'
import { type Component, Show, createEffect, createMemo } from 'solid-js'
import { AuthenticationType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { useApplicationForm } from '../../../provider/applicationFormProvider'
import UrlField from './UrlField'

const Container = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})

const collapsibleClass = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '8px',
  overflow: 'hidden',
  height: '0',

  selectors: {
    '[data-expanded] &': {
      height: 'var(--kb-collapsible-content-height)',
      overflow: 'visible',
    },
  },
})

const collapsibleTriggerClass = style({
  width: '100%',
  appearance: 'none',
  border: 'none',
  background: 'none',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '8px',
  cursor: 'pointer',
  color: colorVars.semantic.text.black,
  ...textVars.text.medium,
})

const ExpandIconContainer = styled('div', {
  base: {
    width: '24px',
    height: '24px',
    transition: 'transform 0.2s',
    selectors: {
      '[data-expanded] &': {
        transform: 'rotate(180deg)',
      },
    },
  },
})

const CollapsibleContentContainer = styled('div', {
  base: {
    paddingTop: '8px',
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})

const DeleteButtonContainer = styled('div', {
  base: {
    position: 'absolute',
    top: '-8px',
    right: '-8px',
  },
})

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
      <Container>
        <DeleteButtonContainer>
          <Button onClick={handleDelete} variants="textError" size="small" type="button">
            Delete
          </Button>
        </DeleteButtonContainer>
        <UrlField index={props.index} readonly={props.readonly} showHttpPort={props.isRuntimeApp} />
        {/* Field componentがmountされないとそのfieldがformに登録されないためforceMountする */}
        <Collapsible forceMount>
          <Collapsible.Trigger class={collapsibleTriggerClass}>
            詳細
            <ExpandIconContainer>
              <MaterialSymbols>expand_more</MaterialSymbols>
            </ExpandIconContainer>
          </Collapsible.Trigger>
          <Collapsible.Content class={collapsibleClass}>
            <CollapsibleContentContainer>
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
            </CollapsibleContentContainer>
          </Collapsible.Content>
        </Collapsible>
      </Container>
    </>
  )
}

export default WebsiteFieldGroup
