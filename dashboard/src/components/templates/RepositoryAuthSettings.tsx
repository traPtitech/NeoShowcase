import { CreateRepositoryAuth, Repository_AuthMethod } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, systemInfo } from '/@/libs/api'
import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Match, Switch, createEffect, createSignal } from 'solid-js'
import { createResource } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { Button } from '../UI/Button'
import { TextInput } from '../UI/TextInput'
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
  const handleCopyPublicKey = () => writeToClipboard(publicKey())

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
        <Match when={props.authConfig.auth.case === 'basic'}>
          <FormItem title="UserName" required>
            <TextInput />
          </FormItem>
          <FormItem title="Password" required>
            <TextInput type="password" />
          </FormItem>
        </Match>
        <Match when={props.authConfig.auth.case === 'ssh'}>
          <Container>
            SSH公開鍵 以下のSSH公開鍵{!useTmpKey() && '(システムデフォルト)'}
            をリポジトリに登録してください。
            <TextInput value={publicKey()} />
            <button onClick={handleCopyPublicKey} type="button">
              Copy
            </button>
            <Button color="textError" size="small" onClick={() => setUseTmpKey(true)}>
              再生成する
            </Button>
          </Container>
        </Match>
      </Switch>
    </>
  )
}
