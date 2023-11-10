import {
  AuthenticationType,
  AvailableDomain,
  CreateWebsiteRequest,
  Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import useModal from '/@/libs/useModal'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, Form, FormStore, getValue, required, reset, setValue, toCustom } from '@modular-forms/solid'
import { For, Show, createEffect, createMemo, createReaction, onMount } from 'solid-js'
import { on } from 'solid-js'
import { systemInfo } from '../../libs/api'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { TextInput } from '../UI/TextInput'
import { ToolTip } from '../UI/ToolTip'
import FormBox from '../layouts/FormBox'
import { CheckBox } from './CheckBox'
import { FormItem } from './FormItem'
import { List } from './List'
import { RadioButtons } from './RadioButtons'
import { SelectItem, SingleSelect } from './Select'

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
const DeleteConfirm = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    overflowY: 'auto',
    borderRadius: '8px',
    background: colorVars.semantic.ui.secondary,
    color: colorVars.semantic.text.black,
    ...textVars.h3.regular,
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
}

const schemeOptions: SelectItem<boolean>[] = [
  { value: false, title: 'http' },
  { value: true, title: 'https' },
]

const authenticationTypeItems: SelectItem<AuthenticationType>[] = [
  { value: AuthenticationType.OFF, title: 'OFF' },
  { value: AuthenticationType.SOFT, title: 'SOFT' },
  { value: AuthenticationType.HARD, title: 'HARD' },
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
  })

  onMount(() => {
    // Reset when both host and domain are no longer undefined
    resetHostAndDomain(
      () =>
        getValue(props.formStore, 'website.host') !== undefined &&
        getValue(props.formStore, 'website.domain') !== undefined,
    )
  })

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

  createEffect(
    on(
      [() => getValue(props.formStore, 'website.host'), () => getValue(props.formStore, 'website.domain')],
      ([host, domain]) => {
        // set fqdn from host and domain on host or domain change
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
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="URL" required>
            <URLContainer>
              <HttpSelectContainer>
                <Field of={props.formStore} name={'website.https'} type="boolean">
                  {(field, fieldProps) => (
                    <ToolTip
                      props={{
                        content: (
                          <>
                            <div>スキーム</div>
                            <div>通常はhttpsが推奨です</div>
                          </>
                        ),
                      }}
                    >
                      <SingleSelect
                        items={schemeOptions}
                        selected={field.value}
                        setSelected={(selected) => {
                          if (selected !== undefined) {
                            setValue(props.formStore, 'website.https', selected)
                          }
                        }}
                        {...fieldProps}
                      />
                    </ToolTip>
                  )}
                </Field>
              </HttpSelectContainer>
              <URLItem>://</URLItem>
              <Show when={getValue(props.formStore, 'website.domain')?.startsWith('*')}>
                <Field of={props.formStore} name={'website.host'} validate={required('Please Enter Hostname')}>
                  {(field, fieldProps) => (
                    <TextInput
                      placeholder="example.trap.show"
                      value={field.value}
                      error={field.error}
                      tooltip={{
                        props: {
                          content: 'ホスト名',
                        },
                      }}
                      {...fieldProps}
                    />
                  )}
                </Field>
              </Show>
              <Field of={props.formStore} name={'website.domain'}>
                {(field, fieldProps) => (
                  <ToolTip
                    props={{
                      content: 'ドメイン',
                    }}
                  >
                    <SingleSelect
                      items={
                        systemInfo()?.domains.map((domain) => {
                          const domainName = domain.domain.replace(/\*/g, '')
                          return {
                            value: domain,
                            title: domainName,
                          }
                        }) ?? []
                      }
                      selected={systemInfo()?.domains.find((d) => d.domain === field.value)}
                      setSelected={(selected) => {
                        if (selected !== undefined) {
                          setValue(props.formStore, 'website.domain', selected.domain)
                        }
                      }}
                      {...fieldProps}
                    />
                  </ToolTip>
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
                  <TextInput
                    value={field.value?.slice(1)}
                    tooltip={{
                      props: {
                        content: '(Advanced) 指定Prefixが付いていたときのみアプリへルーティング',
                      },
                    }}
                    {...fieldProps}
                  />
                )}
              </Field>
              <Show when={props.isRuntimeApp}>
                <URLItem> → </URLItem>
                <HttpSelectContainer>
                  <Field of={props.formStore} name={'website.httpPort'} type="number">
                    {(field, fieldProps) => (
                      <TextInput
                        placeholder="80"
                        type="number"
                        min="0"
                        value={field.value}
                        tooltip={{
                          props: {
                            content: 'アプリのHTTP Port番号',
                          },
                        }}
                        {...fieldProps}
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
              <FormItem
                title="部員認証"
                tooltip={{
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
              >
                <ToolTip
                  props={{
                    content: `${getValue(props.formStore, 'website.domain')}では部員認証が使用できません`,
                  }}
                  disabled={getValue(props.formStore, 'website.authAvailable')}
                >
                  <RadioButtons
                    items={authenticationTypeItems}
                    selected={field.value}
                    setSelected={(selected) => {
                      if (selected !== undefined) {
                        setValue(props.formStore, 'website.authentication', selected)
                      }
                    }}
                    disabled={getValue(props.formStore, 'website.authAvailable') === false}
                    {...fieldProps}
                  />
                </ToolTip>
              </FormItem>
            )}
          </Field>
          <FormItem title="高度な設定">
            <Field of={props.formStore} name={'website.stripPrefix'} type="boolean">
              {(field, fieldProps) => (
                <CheckBox.Option
                  title="Strip Path Prefix"
                  checked={field.value ?? false}
                  setChecked={(selected) => setValue(props.formStore, 'website.stripPrefix', selected)}
                  tooltip={{
                    props: {
                      content: '(Advanced) 指定Prefixをアプリへのリクエスト時に削除',
                    },
                  }}
                  {...fieldProps}
                />
              )}
            </Field>
            <Show when={props.isRuntimeApp}>
              <Field of={props.formStore} name={'website.h2c'} type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    title="Use h2c"
                    checked={field.value ?? false}
                    setChecked={(selected) => setValue(props.formStore, 'website.h2c', selected)}
                    tooltip={{
                      props: {
                        content: '(Advanced) アプリ通信に強制的にh2cを用いる',
                      },
                    }}
                    {...fieldProps}
                  />
                )}
              </Field>
            </Show>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <DeleteButtonContainer>
            <Button onclick={open} variants="textError" size="small" type="button">
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
              disabled={props.formStore.invalid || !props.formStore.dirty || props.formStore.submitting}
            >
              {state() === 'added' ? 'Add' : 'Save'}
            </Button>
          </Show>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Website</Modal.Header>
        <Modal.Body>
          <DeleteConfirm>{websiteUrl()}</DeleteConfirm>
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
              <Button
                variants="primary"
                size="medium"
                rightIcon={<MaterialSymbols>add</MaterialSymbols>}
                onClick={props.addWebsite}
                type="button"
              >
                Add Website
              </Button>
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
          />
        )}
      </For>
      <Show when={props.formStores.length > 0}>
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
