import { InfoTooltip } from '/@/components/InfoTooltip'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Match, Show, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { CreateRepositoryAuth } from '../api/neoshowcase/protobuf/gateway_pb'
import { client, systemInfo } from '../libs/api'
import { vars } from '../theme'
import { Button } from './Button'
import { InputBar, InputLabel } from './Input'
import { Radio } from './Radio'

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

const Row = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
    alignItems: 'center',
  },
})

interface RepositoryAuthSettingsProps {
  authConfig: PlainMessage<CreateRepositoryAuth>
  setAuthConfig: SetStoreFunction<PlainMessage<CreateRepositoryAuth>>
}

export const RepositoryAuthSettings: Component<RepositoryAuthSettingsProps> = (props) => {
  const [useTmpKey, setUseTmpKey] = createSignal(false)
  const [tmpKey] = createResource(
    () => (useTmpKey() ? true : undefined),
    () => client.generateKeyPair({}),
  )
  createEffect(() => {
    if (!tmpKey()) return
    props.setAuthConfig('auth', 'value', { keyId: tmpKey()?.keyId })
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemInfo()?.publicKey)

  return (
    <div>
      <InputLabel>認証方法</InputLabel>
      <Radio
        items={[
          { title: '認証を使用しない', value: 'none' },
          { title: 'Basic認証を使用', value: 'basic' },
          { title: 'SSH認証を使用', value: 'ssh' },
        ]}
        selected={props.authConfig.auth.case}
        setSelected={(v) => props.setAuthConfig('auth', 'case', v)}
      />
      <Switch>
        <Match when={props.authConfig.auth.case === 'basic' && props.authConfig.auth.value}>
          {(v) => (
            <>
              <InputLabel>ユーザー名</InputLabel>
              <InputBar
                // SSH URLはURLとしては不正なのでtypeを変更
                value={v().username}
                onInput={(e) => props.setAuthConfig('auth', 'value', { username: e.currentTarget.value })}
              />
              <InputLabel>パスワード</InputLabel>
              <InputBar
                // SSH URLはURLとしては不正なのでtypeを変更
                type="password"
                value={v().password}
                onInput={(e) => props.setAuthConfig('auth', 'value', { password: e.currentTarget.value })}
              />
            </>
          )}
        </Match>
        <Match when={props.authConfig.auth.case === 'ssh'}>
          <SshDetails>以下のSSH公開鍵{!useTmpKey() && ' (システム共通) '}をリポジトリに登録してください。</SshDetails>
          <PublicKeyCode>{publicKey()}</PublicKeyCode>
          <Show when={!useTmpKey()}>
            <Row>
              <Button color="black1" size="large" width="auto" onclick={() => setUseTmpKey(true)} type="submit">
                新たな鍵ペアを生成 (for github.com)
              </Button>
              <InfoTooltip
                tooltip={[
                  '新しく鍵ペアを生成します',
                  'システム共通の鍵ペアでは動かない場合に利用します',
                  'github.com の場合は必ず生成してください',
                ]}
              />
            </Row>
          </Show>
        </Match>
      </Switch>
    </div>
  )
}
