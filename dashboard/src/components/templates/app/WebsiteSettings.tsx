import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, Form, FormStore, getValue, required, reset, setValue, toCustom } from '@modular-forms/solid'
import { For, Show, createEffect, createMemo, createReaction, onMount } from 'solid-js'
import { on } from 'solid-js'
import {
  AuthenticationType,
  AvailableDomain,
  CreateWebsiteRequest,
  Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import useModal from '/@/libs/useModal'
import { colorVars, textVars } from '/@/theme'
import { systemInfo } from '../../../libs/api'
import { MaterialSymbols } from '../../UI/MaterialSymbols'
import ModalDeleteConfirm from '../../UI/ModalDeleteConfirm'
import { TextField } from '../../UI/TextField'
import FormBox from '../../layouts/FormBox'
import { CheckBox } from '../CheckBox'
import { FormItem } from '../FormItem'
import { List } from '../List'
import { RadioGroup, RadioOption } from '../RadioGroups'
import { SelectOption, SingleSelect } from '../Select'

const URLContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-top',
    gap: '8px',
  },
})
const URLItem = styled('div', {
  base: {
    height: '48px',
    display: 'flex',
    alignItems: 'center',
  },
})
const HttpSelectContainer = styled('div', {
  base: {
    flexShrink: 0,
    width: 'calc(6ch + 60px)',
  },
})
const DeleteButtonContainer = styled('div', {
  base: {
    width: 'fit-content',
    marginRight: 'auto',
  },
})
const AddMoreButtonContainer = styled('div', {
  base: {
    display: 'flex',
    justifyContent: 'center',
  },
})

interface WebsiteSettingProps {
  isRuntimeApp: boolean
  formStore: FormStore<WebsiteSetting, undefined>
  saveWebsite?: () => void
  deleteWebsite: () => void
  hasPermission: boolean
}

const schemeOptions: SelectOption<`${boolean}`>[] = [
  { value: 'false', label: 'http' },
  { value: 'true', label: 'https' },
]

const authenticationTypeOptionsMap = {
  [`${AuthenticationType.OFF}`]: AuthenticationType.OFF,
  [`${AuthenticationType.SOFT}`]: AuthenticationType.SOFT,
  [`${AuthenticationType.HARD}`]: AuthenticationType.HARD,
}

const authenticationTypeOptions: RadioOption<`${AuthenticationType}`>[] = [
  { value: `${AuthenticationType.OFF}`, label: 'OFF' },
  { value: `${AuthenticationType.SOFT}`, label: 'SOFT' },
  { value: `${AuthenticationType.HARD}`, label: 'HARD' },
]

export const WebsiteSetting = (props: WebsiteSettingProps) => {
  const state = () => getValue(props.formStore, 'state')
  const discardChanges = () => reset(props.formStore)

  const { Modal, open, close } = useModal()

  const nonWildcardDomains = createMemo(() => systemInfo()?.domains.filter((d) => !d.domain.startsWith('*')) ?? [])
  const wildCardDomains = createMemo(() => systemInfo()?.domains.filter((d) => d.domain.startsWith('*')) ?? [])
  const websiteUrl = () => {
    const scheme = getValue(props.formStore, 'website.https') ? 'https' : 'http'
    const fqdn = getValue(props.formStore, 'website.fqdn')
    const pathPrefix = getValue(props.formStore, 'website.pathPrefix')
    return `${scheme}://${fqdn}${pathPrefix}`
  }

  const extractHost = (
    fqdn: string,
  ): {
    host: string
    domain: PlainMessage<AvailableDomain>
  } => {
    const matchNonWildcardDomain = nonWildcardDomains().find((d) => fqdn === d.domain)
    if (matchNonWildcardDomain !== undefined) {
      return {
        host: '',
        domain: matchNonWildcardDomain,
      }
    }

    const matchDomain = wildCardDomains().find((d) => fqdn?.endsWith(d.domain.replace(/\*/g, '')))
    if (matchDomain === undefined) {
      const fallbackDomain = systemInfo()?.domains[0]
      if (fallbackDomain === undefined) throw new Error('No domain available')
      return {
        host: '',
        domain: fallbackDomain,
      }
    }
    return {
      host: fqdn.slice(0, -matchDomain.domain.length + 1),
      domain: matchDomain,
    }
  }

  // set host and domain from fqdn on fqdn change
  createEffect(
    on(
      () => getValue(props.formStore, 'website.fqdn'),
      (fqdn) => {
        if (fqdn === undefined) return
        const { host, domain } = extractHost(fqdn)
        setValue(props.formStore, 'website.host', host)
        setValue(props.formStore, 'website.domain', domain.domain)
        setValue(props.formStore, 'website.authAvailable', domain.authAvailable)
        if (domain.authAvailable === false) {
          setValue(props.formStore, 'website.authentication', AuthenticationType.OFF)
        }
      },
    ),
  )

  const resetHostAndDomain = createReaction(() => {
    const fqdn = getValue(props.formStore, 'website.fqdn')
    if (fqdn === undefined) return
    const { host, domain } = extractHost(fqdn)
    reset(props.formStore, 'website.host', {
      initialValue: host,
    })
    reset(props.formStore, 'website.domain', {
      initialValue: domain.domain,
    })
    reset(props.formStore, 'website.authAvailable', {
      initialValue: domain.authAvailable,
    })
  })

  onMount(() => {
    // Reset host and domain on first fqdn change
    resetHostAndDomain(() => getValue(props.formStore, 'website.fqdn'))
  })

  // set fqdn from host and domain on host or domain change
  createEffect(
    on(
      [() => getValue(props.formStore, 'website.host'), () => getValue(props.formStore, 'website.domain')],
      ([host, domain]) => {
        if (host === undefined || domain === undefined) return
        const fqdn = `${host}${domain?.replace(/\*/g, '')}`
        setValue(props.formStore, 'website.fqdn', fqdn)
      },
    ),
  )

  return (
    <Form
      of={props.formStore}
      onSubmit={() => {
        if (props.saveWebsite) props.saveWebsite()
      }}
      style={{ width: '100%' }}
    >
      {/* 
          To make a field active, it must be included in the DOM
          see: https://modularforms.dev/solid/guides/add-fields-to-form#active-state
        */}
      <Field of={props.formStore} name={'state'}>
        {() => <></>}
      </Field>
      <Field of={props.formStore} name={'website.id'}>
        {() => <></>}
      </Field>
      <Field of={props.formStore} name={'website.fqdn'}>
        {() => <></>}
      </Field>
      <Field of={props.formStore} name={'website.authAvailable'} type="boolean">
        {() => <></>}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="URL" required>
            <URLContainer>
              <HttpSelectContainer>
                <Field of={props.formStore} name={'website.https'} type="boolean">
                  {(field, fieldProps) => (
                    <SingleSelect
                      tooltip={{
                        props: {
                          content: (
                            <>
                              <div>スキーム</div>
                              <div>通常はhttpsが推奨です</div>
                            </>
                          ),
                        },
                      }}
                      {...fieldProps}
                      options={schemeOptions}
                      value={field.value ? 'true' : 'false'}
                      setValue={(selected) => {
                        setValue(props.formStore, 'website.https', selected === 'true')
                      }}
                      readOnly={props.hasPermission}
                    />
                  )}
                </Field>
              </HttpSelectContainer>
              <URLItem>://</URLItem>
              <Field of={props.formStore} name={'website.host'} validate={required('Please Enter Hostname')}>
                {(field, fieldProps) => (
                  <Show when={getValue(props.formStore, 'website.domain')?.startsWith('*')}>
                    <TextField
                      placeholder="example.trap.show"
                      tooltip={{
                        props: {
                          content: 'ホスト名',
                        },
                      }}
                      {...fieldProps}
                      value={field.value ?? ''}
                      error={field.error}
                      readOnly={!props.hasPermission}
                    />
                  </Show>
                )}
              </Field>
              <Field of={props.formStore} name={'website.domain'}>
                {(field, fieldProps) => (
                  <SingleSelect
                    tooltip={{
                      props: {
                        content: 'ドメイン',
                      },
                    }}
                    {...fieldProps}
                    options={
                      systemInfo()?.domains.map((domain) => {
                        const domainName = domain.domain.replace(/\*/g, '')
                        return {
                          value: domain.domain,
                          label: domainName,
                        }
                      }) ?? []
                    }
                    value={field.value}
                    setValue={(domain) => {
                      setValue(props.formStore, 'website.domain', domain)
                    }}
                    readOnly={!props.hasPermission}
                  />
                )}
              </Field>
            </URLContainer>
            <URLContainer>
              <URLItem>/</URLItem>
              <Field
                of={props.formStore}
                name={'website.pathPrefix'}
                transform={toCustom((value) => `/${value}` as string, {
                  on: 'input',
                })}
              >
                {(field, fieldProps) => (
                  <TextField
                    tooltip={{
                      props: {
                        content: '(Advanced) 指定Prefixが付いていたときのみアプリへルーティング',
                      },
                    }}
                    {...fieldProps}
                    value={field.value?.slice(1) ?? ''}
                    error={field.error}
                    readOnly={!props.hasPermission}
                  />
                )}
              </Field>
              <Show when={props.isRuntimeApp}>
                <URLItem> → </URLItem>
                <HttpSelectContainer>
                  <Field of={props.formStore} name={'website.httpPort'} type="number">
                    {(field, fieldProps) => (
                      <TextField
                        placeholder="80"
                        type="number"
                        min="0"
                        tooltip={{
                          props: {
                            content: 'アプリのHTTP Port番号',
                          },
                        }}
                        {...fieldProps}
                        value={field.value?.toString() ?? ''}
                        error={field.error}
                        readOnly={!props.hasPermission}
                      />
                    )}
                  </Field>
                </HttpSelectContainer>
                <URLItem>/TCP</URLItem>
              </Show>
            </URLContainer>
          </FormItem>
          <Field of={props.formStore} name={'website.authentication'} type="number">
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
                tooltip={{
                  props: {
                    content: `${getValue(props.formStore, 'website.domain')}では部員認証が使用できません`,
                  },
                  disabled: getValue(props.formStore, 'website.authAvailable') && props.hasPermission,
                }}
                {...fieldProps}
                options={authenticationTypeOptions}
                value={`${field.value ?? AuthenticationType.OFF}`}
                setValue={(value) => {
                  setValue(props.formStore, 'website.authentication', authenticationTypeOptionsMap[value])
                }}
                disabled={!getValue(props.formStore, 'website.authAvailable')}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <FormItem title="高度な設定">
            <Field of={props.formStore} name={'website.stripPrefix'} type="boolean">
              {(field, fieldProps) => (
                <CheckBox.Option
                  {...fieldProps}
                  label="Strip Path Prefix"
                  checked={field.value ?? false}
                  readOnly={!props.hasPermission}
                />
              )}
            </Field>
            <Show when={props.isRuntimeApp}>
              <Field of={props.formStore} name={'website.h2c'} type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    {...fieldProps}
                    label="Use h2c"
                    checked={field.value ?? false}
                    readOnly={!props.hasPermission}
                  />
                )}
              </Field>
            </Show>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <DeleteButtonContainer>
            <Button
              onclick={open}
              variants="textError"
              size="small"
              type="button"
              disabled={!props.hasPermission}
              tooltip={{
                props: {
                  content: !props.hasPermission
                    ? '設定を削除するにはアプリケーションのオーナーになる必要があります'
                    : undefined,
                },
              }}
            >
              Delete
            </Button>
          </DeleteButtonContainer>
          <Show when={state() !== 'added' && props.formStore.dirty}>
            <Button onclick={discardChanges} variants="borderError" size="small" type="button">
              Discard Changes
            </Button>
          </Show>
          <Show when={props.saveWebsite !== undefined}>
            <Button
              variants="primary"
              size="small"
              type="submit"
              disabled={
                props.formStore.invalid || !props.formStore.dirty || props.formStore.submitting || !props.hasPermission
              }
              loading={props.formStore.submitting}
              tooltip={{
                props: {
                  content: !props.hasPermission
                    ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                    : undefined,
                },
              }}
            >
              {state() === 'added' ? 'Add' : 'Save'}
            </Button>
          </Show>
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
          <Button onclick={props.deleteWebsite} variants="primaryError" size="medium" type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </Form>
  )
}

type FQDN = {
  host: string
  domain: PlainMessage<AvailableDomain>['domain']
  authAvailable: PlainMessage<AvailableDomain>['authAvailable']
}

export type WebsiteSetting =
  | {
      /**
       *  - `noChange`: 既存の設定を変更していない
       *  - `readyToChange`: 次の保存時に変更を反映する
       *  - `readyToDelete`: 次の保存時に削除する
       */
      state: 'noChange' | 'readyToChange' | 'readyToDelete'
      website: PlainMessage<Website> & FQDN
    }
  | {
      /**
       *  - `added`: 新規に設定を追加した
       */
      state: 'added'
      website: PlainMessage<CreateWebsiteRequest> & FQDN
    }

export type WebsiteSettingForm = {
  websites: WebsiteSetting[]
}

export const newWebsite = (): PlainMessage<CreateWebsiteRequest> => ({
  fqdn: '',
  pathPrefix: '/',
  stripPrefix: false,
  https: true,
  h2c: false,
  httpPort: 0,
  authentication: AuthenticationType.OFF,
})

const PlaceHolder = styled('div', {
  base: {
    width: '100%',
    height: '400px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})

interface WebsiteSettingsProps {
  isRuntimeApp: boolean
  formStores: FormStore<WebsiteSetting, undefined>[]
  addWebsite: () => void
  applyChanges: () => void
  hasPermission: boolean
}

export const WebsiteSettings = (props: WebsiteSettingsProps) => {
  return (
    <Show when={systemInfo()}>
      <For
        each={props.formStores}
        fallback={
          <List.Container>
            <PlaceHolder>
              <MaterialSymbols displaySize={80}>link_off</MaterialSymbols>
              No Websites Configured
              <Show when={props.hasPermission}>
                <Button
                  variants="primary"
                  size="medium"
                  rightIcon={<MaterialSymbols>add</MaterialSymbols>}
                  onClick={props.addWebsite}
                  type="button"
                >
                  Add Website
                </Button>
              </Show>
            </PlaceHolder>
          </List.Container>
        }
      >
        {(form, index) => (
          <WebsiteSetting
            isRuntimeApp={props.isRuntimeApp}
            formStore={form}
            saveWebsite={() => {
              if (getValue(props.formStores[index()], 'state') === 'noChange') {
                setValue(props.formStores[index()], 'state', 'readyToChange')
              }
              props.applyChanges()
            }}
            deleteWebsite={() => {
              if (getValue(props.formStores[index()], 'state') === 'added') {
                // 新規追加した設定をformStoresから削除
                props.formStores.splice(index(), 1)
              } else {
                // すでに保存されている設定を削除
                setValue(props.formStores[index()], 'state', 'readyToDelete')
                props.applyChanges()
              }
            }}
            hasPermission={props.hasPermission}
          />
        )}
      </For>
      <Show when={props.formStores.length > 0 && props.hasPermission}>
        <AddMoreButtonContainer>
          <Button
            onclick={props.addWebsite}
            variants="border"
            size="small"
            leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
            type="button"
          >
            Add More
          </Button>
        </AddMoreButtonContainer>
      </Show>
    </Show>
  )
}
