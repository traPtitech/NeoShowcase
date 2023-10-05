import { CreateRepositoryAuth } from '/@/api/neoshowcase/protobuf/gateway_pb'
import RefreshIcon from '/@/assets/icons/20/replay.svg'
import { client, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Match, Show, Switch, createEffect, createSignal } from 'solid-js'
import { createResource } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { Button } from '../UI/Button'
import { TextInput } from '../UI/TextInput'
import { CopyableInput } from './CopyableInput'
import { FormItem } from './FormItem'
import { RadioButtons, RadioItem } from './RadioButtons'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',

    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const SshKeyContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const RefreshButtonContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    color: colorVars.semantic.accent.error,
    ...textVars.caption.regular,
  },
})

const AuthMethods: RadioItem<CreateRepositoryAuth['auth']['case']>[] = [
  { title: 'SSH', value: 'ssh' },
  { title: 'HTTPS', value: 'basic' },
  { title: '認証を使用しない', value: 'none' },
]

interface Props {
  authConfig: PlainMessage<CreateRepositoryAuth>
  setAuthConfig: SetStoreFunction<PlainMessage<CreateRepositoryAuth>>
}

export const RepositoryAuthSettings: Component<Props> = (props) => {
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
    <>
      <Container>
        認証方法
        <RadioButtons
          items={AuthMethods}
          selected={props.authConfig.auth.case}
          setSelected={(v) => props.setAuthConfig('auth', 'case', v)}
        />
      </Container>
      <Switch>
        <Match when={props.authConfig.auth.case === 'basic' && props.authConfig.auth.value}>
          {(v) => (
            <>
              <FormItem title="UserName" required>
                <TextInput
                  value={v().username}
                  onInput={(e) =>
                    props.setAuthConfig('auth', 'value', {
                      username: e.currentTarget.value,
                    })
                  }
                />
              </FormItem>
              <FormItem title="Password" required>
                <TextInput
                  type="password"
                  value={v().password}
                  onInput={(e) =>
                    props.setAuthConfig('auth', 'value', {
                      password: e.currentTarget.value,
                    })
                  }
                />
              </FormItem>
            </>
          )}
        </Match>
        <Match when={props.authConfig.auth.case === 'ssh'}>
          <Container>
            SSH公開鍵
            <SshKeyContainer>
              以下のSSH公開鍵{!useTmpKey() && '(システムデフォルト)'}
              をリポジトリに登録してください
              <CopyableInput value={publicKey()} />
              <Show when={!useTmpKey()}>
                <RefreshButtonContainer>
                  <Button color="textError" size="small" onClick={() => setUseTmpKey(true)} leftIcon={<RefreshIcon />}>
                    再生成する
                  </Button>
                  For Github.com
                </RefreshButtonContainer>
              </Show>
            </SshKeyContainer>
          </Container>
        </Match>
      </Switch>
    </>
  )
}
