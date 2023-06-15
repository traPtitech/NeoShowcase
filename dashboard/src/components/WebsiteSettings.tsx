import { AuthenticationType, CreateWebsiteRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Checkbox } from '/@/components/Checkbox'
import { Radio, RadioItem } from '/@/components/Radio'
import { Button } from '/@/components/Button'
import { SetStoreFunction } from 'solid-js/store'
import { For, Show } from 'solid-js'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormButton, FormCheckBox, FormSettings, FormSettingsButton, SettingsContainer } from '/@/components/AppsNew'
import { availableDomains } from '../libs/api'
import { styled } from '@macaron-css/solid'
import { vars } from '../theme'
import { PlainMessage } from '@bufbuild/protobuf'

const AvailableDomainContainer = styled('div', {
  base: {
    fontSize: '14px',
    color: vars.text.black2,
    padding: '8px',
  },
})
const AvailableDomainUl = styled('ul', {
  base: {
    margin: '8px 0',
  },
})

interface WebsiteSettingProps {
  runtime: boolean
  website: PlainMessage<CreateWebsiteRequest>
  setWebsite: <T extends keyof PlainMessage<CreateWebsiteRequest>>(
    valueName: T,
    value: PlainMessage<CreateWebsiteRequest>[T],
  ) => void
  deleteWebsite: () => void
}

const authenticationTypeItems: RadioItem<AuthenticationType>[] = [
  { value: AuthenticationType.OFF, title: 'OFF' },
  { value: AuthenticationType.SOFT, title: 'SOFT' },
  { value: AuthenticationType.HARD, title: 'HARD' },
]

export const WebsiteSetting = (props: WebsiteSettingProps) => {
  return (
    <FormSettings>
      <div>
        <InputLabel>ドメイン名</InputLabel>
        <InputBar
          placeholder='example.ns.trap.jp'
          value={props.website.fqdn}
          onInput={(e) => props.setWebsite('fqdn', e.target.value)}
        />
        <AvailableDomainContainer>
          使用可能なドメイン
          <AvailableDomainUl>
            <For each={availableDomains()?.domains}>
              {(domain) => (
                <li>
                  {domain.domain}
                  <Show when={domain.excludeDomains.length > 0}>&nbsp;({domain.excludeDomains.join(', ')}を除く)</Show>
                  ：{domain.authAvailable ? '部員認証の使用可能' : '部員認証の使用不可'}
                </li>
              )}
            </For>
          </AvailableDomainUl>
        </AvailableDomainContainer>
      </div>
      <div>
        <InputLabel>Path Prefix</InputLabel>
        <InputBar
          placeholder='/'
          value={props.website.pathPrefix}
          onInput={(e) => props.setWebsite('pathPrefix', e.target.value)}
        />
      </div>
      <div>
        <FormCheckBox>
          <Checkbox
            selected={props.website.stripPrefix}
            setSelected={(selected) => props.setWebsite('stripPrefix', selected)}
          >
            Strip Path Prefix
          </Checkbox>
          <Checkbox selected={props.website.https} setSelected={(selected) => props.setWebsite('https', selected)}>
            https
          </Checkbox>
          <Show when={props.runtime}>
            <Checkbox selected={props.website.h2c} setSelected={(selected) => props.setWebsite('h2c', selected)}>
              (advanced) アプリ通信にh2cを用いる
            </Checkbox>
          </Show>
        </FormCheckBox>
      </div>
      <Show when={props.runtime}>
        <div>
          <InputLabel>アプリのHTTP Port番号</InputLabel>
          <InputBar
            placeholder='80'
            type='number'
            value={props.website.httpPort}
            onChange={(e) => props.setWebsite('httpPort', +e.target.value)}
          />
        </div>
      </Show>
      <div>
        <InputLabel>部員認証</InputLabel>
        <Radio
          items={authenticationTypeItems}
          selected={props.website.authentication}
          setSelected={(selected) => props.setWebsite('authentication', selected)}
        />
      </div>
      <FormSettingsButton>
        <Button onclick={props.deleteWebsite} color='black1' size='large' width='auto' type='button'>
          Delete website setting
        </Button>
      </FormSettingsButton>
    </FormSettings>
  )
}

interface WebsiteSettingsProps {
  runtime: boolean
  websiteConfigs: PlainMessage<CreateWebsiteRequest>[]
  setWebsiteConfigs: SetStoreFunction<PlainMessage<CreateWebsiteRequest>[]>
}

export const WebsiteSettings = (props: WebsiteSettingsProps) => {
  return (
    <SettingsContainer>
      <For each={props.websiteConfigs}>
        {(website, i) => (
          <WebsiteSetting
            runtime={props.runtime}
            website={website}
            setWebsite={(valueName, value) => {
              props.setWebsiteConfigs(i(), valueName, value)
            }}
            deleteWebsite={() =>
              props.setWebsiteConfigs((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
            }
          />
        )}
      </For>

      <FormButton>
        <Button
          onclick={() => {
            props.setWebsiteConfigs([
              ...props.websiteConfigs,
              {
                fqdn: '',
                pathPrefix: '/',
                stripPrefix: false,
                https: false,
                h2c: false,
                httpPort: 0,
                authentication: AuthenticationType.OFF,
              },
            ])
          }}
          color='black1'
          size='large'
          width='auto'
          type='button'
        >
          Add website setting
        </Button>
      </FormButton>
    </SettingsContainer>
  )
}
