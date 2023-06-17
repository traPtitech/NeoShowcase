import { Component, Match, Show, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { InputBar, InputLabel } from './Input'
import { Radio } from './Radio'
import { vars } from '../theme'
import { styled } from '@macaron-css/solid'
import { Button } from './Button'
import { CreateRepositoryAuth } from '../api/neoshowcase/protobuf/gateway_pb'
import { SetStoreFunction } from 'solid-js/store'
import { client } from '../libs/api'
import { PlainMessage } from '@bufbuild/protobuf'

const SshDetails = styled('div', {
  base: {
    color: vars.text.black2,
    marginBottom: '4px',
  },
})

const PublicKeyCode = styled('code', {
  base: {
    display: 'block',
    padding: '8px 12px',
    fontSize: '14px',
    background: vars.bg.white2,
    color: vars.text.black1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
  },
})

export type AuthMethod = Exclude<CreateRepositoryAuth['auth']['case'], undefined>
export type AuthConfig = {
  [K in AuthMethod]: Extract<PlainMessage<CreateRepositoryAuth>['auth'], { case: K }>
} & {
  authMethod: AuthMethod
}

interface RepositoryAuthSettingsProps {
  authConfig: AuthConfig
  setAuthConfig: SetStoreFunction<AuthConfig>
}

export const RepositoryAuthSettings: Component<RepositoryAuthSettingsProps> = (props) => {
  const [systemPublicKey] = createResource(() => client.getSystemPublicKey({}))
  const [useTmpKey, setUseTmpKey] = createSignal(false)
  const [tmpKey] = createResource(
    () => (useTmpKey() ? true : undefined),
    () => client.generateKeyPair({}),
  )
  createEffect(() => {
    if (!tmpKey()) return
    props.setAuthConfig('ssh', 'value', { keyId: tmpKey()?.keyId })
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemPublicKey()?.publicKey)

  return (
    <div>
      <InputLabel>認証方法</InputLabel>
      <Radio
        items={[
          { title: '認証を使用しない', value: 'none' },
          { title: 'Basic認証を使用', value: 'basic' },
          { title: 'SSH認証を使用', value: 'ssh' },
        ]}
        selected={props.authConfig.authMethod}
        setSelected={(v) => props.setAuthConfig('authMethod', v)}
      />
      <Switch>
        <Match when={props.authConfig.authMethod === 'basic'}>
          <InputLabel>ユーザー名</InputLabel>
          <InputBar
            // SSH URLはURLとしては不正なのでtypeを変更
            value={props.authConfig.basic.value.username}
            onInput={(e) => props.setAuthConfig('basic', 'value', { username: e.currentTarget.value })}
          />
          <InputLabel>パスワード</InputLabel>
          <InputBar
            // SSH URLはURLとしては不正なのでtypeを変更
            type='password'
            value={props.authConfig.basic.value.password}
            onInput={(e) => props.setAuthConfig('basic', 'value', { password: e.currentTarget.value })}
          />
        </Match>
        <Match when={props.authConfig.authMethod === 'ssh'}>
          <SshDetails>
            以下のSSH公開鍵{!useTmpKey() && ' (システムデフォルト) '}をリポジトリに登録してください。
          </SshDetails>
          <PublicKeyCode>{publicKey()}</PublicKeyCode>
          <Show when={!useTmpKey()}>
            <Button color='black1' size='large' width='auto' onclick={() => setUseTmpKey(true)} type='submit'>
              新たなSSH鍵を生成する (for github.com)
            </Button>
          </Show>
        </Match>
      </Switch>
    </div>
  )
}
