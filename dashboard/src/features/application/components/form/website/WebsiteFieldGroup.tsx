import {
  Field,
  Form,
  type FormStore,
  type SubmitHandler,
  getErrors,
  getValue,
  reset,
  setValue,
  submit,
  validate,
} from '@modular-forms/solid'
import { type Component, Show, createMemo } from 'solid-js'
import { AuthenticationType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import FormBox from '/@/components/layouts/FormBox'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { systemInfo } from '/@/libs/api'
import useModal from '/@/libs/useModal'
import type { CreateWebsiteInput } from '../../../schema/websiteSchema'
import UrlField from './UrlField'

const authenticationTypeOptions: RadioOption<`${AuthenticationType}`>[] = [
  { value: `${AuthenticationType.OFF}`, label: 'OFF' },
  { value: `${AuthenticationType.SOFT}`, label: 'SOFT' },
  { value: `${AuthenticationType.HARD}`, label: 'HARD' },
]

type Props = {
  formStore: FormStore<CreateWebsiteInput>
  isRuntimeApp: boolean
  applyChanges: () => Promise<void> | void
  readonly?: boolean
}

const WebsiteFieldGroup: Component<Props> = (props) => {
  const { Modal, open, close } = useModal()

  const availableDomains = systemInfo()?.domains ?? []
  const selectedDomain = createMemo(() => {
    const domainString = getValue(props.formStore, 'domain')
    return availableDomains.find((d) => d.domain === domainString)
  })
  const authAvailable = (): boolean => selectedDomain()?.authAvailable ?? false

  const discardChanges = () => {
    reset(props.formStore)
  }

  const websiteUrl = () => {
    const scheme = getValue(props.formStore, 'https') ? 'https' : 'http'
    const subDomain = getValue(props.formStore, 'subdomain') ?? ''
    const domain = getValue(props.formStore, 'domain') ?? ''
    const fqdn = domain.startsWith('*') ? `${subDomain}${domain.slice(1)}` : domain

    const pathPrefix = getValue(props.formStore, 'pathPrefix')
    return `${scheme}://${fqdn}/${pathPrefix}`
  }

  const handleSave = async () => {
    // submit前のstateを保存しておく
    const originalState = getValue(props.formStore, 'state')
    if (!originalState) throw new Error("The field of state does not exist.")

    try {
      setValue(props.formStore, 'state', 'readyToChange')
      submit(props.formStore)
      
      const errors = getErrors(props.formStore)
      console.log(errors);
    } catch (e) {
      console.log(e)
      setValue(props.formStore, 'state', originalState)
    }
  }

  const handleDelete = () => {
    setValue(props.formStore, 'state', 'readyToDelete')
    submit(props.formStore)
    close()
  }

  const handleSubmit: SubmitHandler<CreateWebsiteInput> = () => {
    props.applyChanges()
  }

  const state = () => getValue(props.formStore, 'state')

  const disableSaveButton = () =>
    props.formStore.invalid ||
    props.formStore.submitting ||
    (state() !== 'added' && !props.formStore.dirty) ||
    props.readonly

  return (
    <Form of={props.formStore} onSubmit={handleSubmit}>
      <Field of={props.formStore} name={'state'}>
        {() => null}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <UrlField formStore={props.formStore} readonly={props.readonly} showHttpPort={props.isRuntimeApp} />
          <Field of={props.formStore} name={'authentication'}>
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
                value={`${field.value ?? AuthenticationType.OFF}`}
                disabled={!authAvailable()}
                readOnly={props.readonly}
              />
            )}
          </Field>
          <FormItem title="高度な設定">
            <Field of={props.formStore} name={'stripPrefix'} type="boolean">
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
              <Field of={props.formStore} name={'h2c'} type="boolean">
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
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            onclick={open}
            variants="textError"
            size="small"
            type="button"
            disabled={props.readonly}
            tooltip={{
              props: {
                content: props.readonly
                  ? '設定を削除するにはアプリケーションのオーナーになる必要があります'
                  : undefined,
              },
            }}
          >
            Delete
          </Button>
          <Show when={props.formStore.dirty && !props.formStore.submitting}>
            <Button variants="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
          </Show>
          <Button
            variants="primary"
            size="small"
            type="button"
            onClick={handleSave}
            disabled={disableSaveButton()}
            loading={props.formStore.submitting}
            tooltip={{
              props: {
                content: props.readonly
                  ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                  : undefined,
              },
            }}
          >
            {state() === 'added' ? 'Add' : 'Save'}
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Website</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            <MaterialSymbols>language</MaterialSymbols>
            {websiteUrl()}
          </ModalDeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button onclick={close} variants="text" size="medium" type="button">
            No, Cancel
          </Button>
          <Button onClick={handleDelete} variants="primaryError" size="medium" type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </Form>
  )
}

export default WebsiteFieldGroup
