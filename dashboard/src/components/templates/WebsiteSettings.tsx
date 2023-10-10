import {
  AuthenticationType,
  AvailableDomain,
  CreateWebsiteRequest,
  Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { getWebsiteURL } from '/@/libs/application'
import useModal from '/@/libs/useModal'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { For, Show, createEffect, createMemo, createSignal } from 'solid-js'
import { on } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { systemInfo } from '../../libs/api'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { TextInput } from '../UI/TextInput'
import { ToolTip } from '../UI/ToolTip'
import FormBox from '../layouts/FormBox'
import { CheckBox } from './CheckBox'
import { FormItem } from './FormItem'
import { RadioButtons } from './RadioButtons'
import { SelectItem, SingleSelect } from './Select'

const URLContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
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
  state: WebsiteSetting['state']
  website: PlainMessage<CreateWebsiteRequest>
  setWebsite: <T extends keyof PlainMessage<CreateWebsiteRequest>>(
    valueName: T,
    value: PlainMessage<CreateWebsiteRequest>[T],
  ) => void
  saveWebsite: () => void
  deleteWebsite: () => void
  discardChanges: () => void
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
  const [host, setHost] = createSignal('')
  const [domain, setDomain] = createSignal<PlainMessage<AvailableDomain>>()

  const { Modal, open, close } = useModal()

  const nonWildcardDomains = createMemo(() => systemInfo()?.domains.filter((d) => !d.domain.startsWith('*')) ?? [])
  const wildCardDomains = createMemo(() => systemInfo()?.domains.filter((d) => d.domain.startsWith('*')) ?? [])

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

    const matchDomain = wildCardDomains().find((d) => fqdn.endsWith(d.domain.replace(/\*/g, '')))
    if (matchDomain === undefined) {
      return {
        host: '',
        domain: systemInfo()?.domains[0],
      }
    }
    return {
      host: fqdn.slice(0, -matchDomain.domain.length + 1),
      domain: matchDomain,
    }
  }

  createEffect(() => {
    const { host, domain } = extractHost(props.website.fqdn)
    setHost(host)
    setDomain(domain)
  })

  createEffect(
    on([host, domain], () => {
      props.setWebsite('fqdn', `${host()}${domain()?.domain.replace(/\*/g, '')}`)
    }),
  )

  createEffect(() => {
    if (domain()?.authAvailable === false) {
      props.setWebsite('authentication', AuthenticationType.OFF)
    }
  })

  return (
    <>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="URL" required>
            <URLContainer>
              <HttpSelectContainer>
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
                    selected={props.website.https}
                    setSelected={(selected) => props.setWebsite('https', selected)}
                  />
                </ToolTip>
              </HttpSelectContainer>
              <span>://</span>
              <Show when={domain()?.domain.startsWith('*')}>
                <TextInput
                  placeholder="example.trap.show"
                  value={host()}
                  onInput={(e) => setHost(e.target.value)}
                  tooltip={{
                    props: {
                      content: 'ホスト名',
                    },
                  }}
                />
              </Show>
              <ToolTip
                props={{
                  content: 'ドメイン',
                }}
              >
                <SingleSelect
                  selected={domain()}
                  setSelected={(selected) => {
                    setDomain(selected)
                  }}
                  items={
                    systemInfo()?.domains.map((domain) => {
                      const domainName = domain.domain.replace(/\*/g, '')
                      return {
                        value: domain,
                        title: domainName,
                      }
                    }) ?? []
                  }
                />
              </ToolTip>
            </URLContainer>
            <URLContainer>
              <span>/</span>
              <TextInput
                value={props.website.pathPrefix.slice(1)}
                onInput={(e) => props.setWebsite('pathPrefix', `/${e.target.value}`)}
                tooltip={{
                  props: {
                    content: '(Advanced) 指定Prefixが付いていたときのみアプリへルーティング',
                  },
                }}
              />
              <Show when={props.isRuntimeApp}>
                <span> → </span>
                <HttpSelectContainer>
                  <TextInput
                    placeholder="80"
                    type="number"
                    min="0"
                    value={props.website.httpPort || ''}
                    onChange={(e) => props.setWebsite('httpPort', +e.target.value)}
                    tooltip={{
                      props: {
                        content: 'アプリのHTTP Port番号',
                      },
                    }}
                  />
                </HttpSelectContainer>
                <span>/TCP</span>
              </Show>
            </URLContainer>
          </FormItem>
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
                content: `${domain()?.domain}では部員認証が使用できません`,
              }}
              disabled={domain()?.authAvailable}
            >
              <RadioButtons
                items={authenticationTypeItems}
                selected={props.website.authentication}
                setSelected={(selected) => props.setWebsite('authentication', selected)}
                disabled={domain()?.authAvailable === false}
              />
            </ToolTip>
          </FormItem>
          <FormItem title="高度な設定">
            <CheckBox.Option
              title="Strip Path Prefix"
              checked={props.website.stripPrefix}
              setChecked={(selected) => props.setWebsite('stripPrefix', selected)}
              tooltip={{
                props: {
                  content: '(Advanced) 指定Prefixをアプリへのリクエスト時に削除',
                },
              }}
            />
            <Show when={props.isRuntimeApp}>
              <CheckBox.Option
                title="Use h2c"
                checked={props.website.h2c}
                setChecked={(selected) => props.setWebsite('h2c', selected)}
                tooltip={{
                  props: {
                    content: '(Advanced) アプリ通信に強制的にh2cを用いる',
                  },
                }}
              />
            </Show>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <DeleteButtonContainer>
            <Button onclick={open} color="textError" size="small" type="button">
              Delete
            </Button>
          </DeleteButtonContainer>
          <Show when={props.state === 'modified'}>
            <Button onclick={props.discardChanges} color="borderError" size="small" type="button">
              Discard Changes
            </Button>
          </Show>
          <Button
            onclick={props.saveWebsite}
            color="primary"
            size="small"
            type="button"
            disabled={props.state === 'noChange'}
          >
            <Show when={props.state === 'added'} fallback={'Save'}>
              Add
            </Show>
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Website</Modal.Header>
        <Modal.Body>
          <DeleteConfirm>{getWebsiteURL(props.website)}</DeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button onclick={close} color="text" size="medium" type="button">
            No, Cancel
          </Button>
          <Button
            onclick={props.deleteWebsite}
            color="primaryError"
            size="medium"
            type="button"
            disabled={props.state === 'noChange'}
          >
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export type WebsiteSetting =
  | {
      /**
       *  - `modified`: 既存の設定を変更した
       *  - `noChange`: 既存の設定を変更していない
       *  - `readyToDelete`: 次の保存時に削除する
       */
      state: 'modified' | 'noChange' | 'readyToDelete'
      website: PlainMessage<Website>
    }
  | {
      /**
       *  - `added`: 新規に設定を追加した
       *  - `readyToSave`: 次の保存時に変更を反映する
       */
      state: 'added' | 'readyToSave'
      website: PlainMessage<CreateWebsiteRequest>
    }

const newWebsite = (): PlainMessage<CreateWebsiteRequest> => ({
  fqdn: '',
  pathPrefix: '/',
  stripPrefix: false,
  https: true,
  h2c: false,
  httpPort: 0,
  authentication: AuthenticationType.OFF,
})

interface WebsiteSettingsProps {
  isRuntimeApp: boolean
  websiteConfigs: WebsiteSetting[]
  setWebsiteConfigs: SetStoreFunction<WebsiteSetting[]>
  applyChanges: () => void
  refetchApp: () => void
}

export const WebsiteSettings = (props: WebsiteSettingsProps) => {
  return (
    <Show when={systemInfo()}>
      <For each={props.websiteConfigs}>
        {(config, i) => (
          <WebsiteSetting
            isRuntimeApp={props.isRuntimeApp}
            state={config.state}
            website={config.website}
            setWebsite={(valueName, value) => {
              if (config.state === 'noChange') {
                props.setWebsiteConfigs(i(), 'state', 'modified')
              }
              props.setWebsiteConfigs(i(), 'website', valueName, value)
            }}
            saveWebsite={() => {
              props.setWebsiteConfigs(i(), 'state', 'readyToSave')
              props.applyChanges()
            }}
            deleteWebsite={() => {
              if (config.state === 'added') {
                props.setWebsiteConfigs((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
              } else {
                props.setWebsiteConfigs(i(), 'state', 'readyToDelete')
                props.applyChanges()
              }
            }}
            discardChanges={() => {
              props.setWebsiteConfigs(i(), 'state', 'noChange')
              props.refetchApp()
            }}
          />
        )}
      </For>
      <AddMoreButtonContainer>
        <Button
          onclick={() =>
            props.setWebsiteConfigs([
              ...props.websiteConfigs,
              {
                state: 'added',
                website: newWebsite(),
              },
            ])
          }
          color="border"
          size="small"
          leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
        >
          Add More
        </Button>
      </AddMoreButtonContainer>
    </Show>
  )
}
