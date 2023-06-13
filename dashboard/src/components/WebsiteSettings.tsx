import { AuthenticationType, CreateWebsiteRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Checkbox } from '/@/components/Checkbox'
import { Radio, RadioItem } from '/@/components/Radio'
import { Button } from '/@/components/Button'
import { SetStoreFunction } from 'solid-js/store'
import { For } from 'solid-js'
import { storify } from '/@/libs/storify'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormButton, FormCheckBox, FormSettings, FormSettingsButton, SettingsContainer } from '/@/components/AppsNew'

interface WebsiteSettingProps {
  website: CreateWebsiteRequest
  setWebsite: <T extends keyof CreateWebsiteRequest>(valueName: T, value: CreateWebsiteRequest[T]) => void
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
          <Checkbox selected={props.website.h2c} setSelected={(selected) => props.setWebsite('h2c', selected)}>
            (advanced) アプリ通信にh2cを用いる
          </Checkbox>
        </FormCheckBox>
      </div>
      <div>
        <InputLabel>アプリのHTTP Port番号</InputLabel>
        <InputBar
          placeholder='80'
          type='number'
          value={props.website.httpPort}
          onChange={(e) => props.setWebsite('httpPort', +e.target.value)}
        />
      </div>
      <div>
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
  websiteConfigs: CreateWebsiteRequest[]
  setWebsiteConfigs: SetStoreFunction<CreateWebsiteRequest[]>
}

export const WebsiteSettings = (props: WebsiteSettingsProps) => {
  return (
    <SettingsContainer>
      <For each={props.websiteConfigs}>
        {(website, i) => (
          <WebsiteSetting
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
            props.setWebsiteConfigs([...props.websiteConfigs, storify(new CreateWebsiteRequest())])
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
