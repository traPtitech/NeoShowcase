import { For, JSX, JSXElement, Match, Show, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import { Empty } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import {
  CreateRepositoryAuth,
  CreateRepositoryAuthBasic,
  CreateRepositoryAuthSSH,
  CreateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/Button'
import { Header } from '/@/components/Header'
import { client } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'

// copy from /pages/apps AppsTitle component
const PageTitle = styled('div', {
  base: {
    marginTop: '48px',
    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

// copy from /pages/apps
// and delete unnecessary styles
const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
  },
})

interface FormProps {
  label: string
  type?: JSX.InputHTMLAttributes<HTMLInputElement>['type']
  placeholder?: JSX.InputHTMLAttributes<HTMLInputElement>['placeholder']
  value: JSX.InputHTMLAttributes<HTMLInputElement>['value']
  onInput: JSX.InputHTMLAttributes<HTMLInputElement>['onInput']
}

const Form = (props: FormProps): JSXElement => {
  return (
    <div>
      <label>{props.label}</label>
      <input
        type={props.type ?? 'text'}
        placeholder={props.placeholder ?? ''}
        value={props.value}
        onInput={props.onInput}
      />
    </div>
  )
}

interface ToggleButtonsProps<T> {
  items: {
    value: T
    label: string
  }[]
  selected: T
  onChange: (value: T) => void
}

const ToggleButtons = <T extends unknown>(props: ToggleButtonsProps<T>): JSXElement => {
  return (
    <div>
      <For each={props.items}>
        {(item) => (
          <Button color='black1' size='large' onclick={() => props.onChange(item.value)}>
            {item.label}
          </Button>
        )}
      </For>
    </div>
  )
}

const SystemPublicKey = (): JSXElement => {
  const [systemPublicKey] = createResource(() => client.getSystemPublicKey({}))

  return (
    <div>
      <div>公開鍵を入力せずにSSH認証でリポジトリを登録する場合、以下のSSH公開鍵が認証に使用されます。</div>
      <Switch>
        <Match when={systemPublicKey.loading}>
          <div>Loading...</div>
        </Match>
        <Match when={systemPublicKey()}>
          <code>{systemPublicKey().publicKey}</code>
        </Match>
      </Switch>
    </div>
  )
}

export default () => {
  // 認証方法 ("none" | "ssh" | "basic")
  type AuthMethod = CreateRepositoryAuth['auth']['case']
  const [authMethod, setAuthMethod] = createSignal<AuthMethod>('none')

  const [sshAuthConfig, setSshAuthConfig] = createStore(new CreateRepositoryAuthSSH())
  const [basicAuthConfig, setBasicAuthConfig] = createStore(new CreateRepositoryAuthBasic())

  const [requestConfig, setRequestConfig] = createStore(
    new CreateRepositoryRequest({
      auth: new CreateRepositoryAuth(),
    }),
  )

  const createRepository = async () => {
    // 認証方法に応じて認証情報を設定
    switch (authMethod()) {
      case 'none':
        setRequestConfig('auth', 'auth', { value: new Empty(), case: 'none' })
        break
      case 'ssh':
        setRequestConfig('auth', 'auth', { value: sshAuthConfig, case: 'ssh' })
        break
      case 'basic':
        setRequestConfig('auth', 'auth', { value: basicAuthConfig, case: 'basic' })
        break
    }

    const res = await client.createRepository(requestConfig)
    // TODO: navigate to repository page when success / show error message when failed
  }

  // URLからリポジトリ名を自動入力
  createEffect(() => {
    const segments = requestConfig.url.split('/')
    const lastSegment = segments.pop() || segments.pop() // 末尾のスラッシュを除去
    const repositoryName = lastSegment?.replace(/\.git$/, '') ?? ''
    setRequestConfig('name', repositoryName)
  })

  return (
    <Container>
      <Header />
      <PageTitle>Create Repository</PageTitle>
      <ContentContainer>
        <Form
          label='URL'
          type='url'
          placeholder='https://example.com/my-app.git'
          value={requestConfig.url}
          onInput={(e) =>
            setRequestConfig({
              url: e.currentTarget.value,
            })
          }
        />
        <Form
          label='リポジトリ名'
          placeholder='my-app'
          value={requestConfig.name}
          onInput={(e) =>
            setRequestConfig({
              name: e.currentTarget.value,
            })
          }
        />
        <ToggleButtons<AuthMethod>
          items={[
            { label: '認証を使用しない', value: 'none' },
            { label: 'Basic認証を使用', value: 'basic' },
            { label: 'SSH認証を使用', value: 'ssh' },
          ]}
          selected={authMethod()}
          onChange={setAuthMethod}
        />
        <Switch>
          <Match when={authMethod() === 'basic'}>
            <Form
              label='ユーザー名'
              value={basicAuthConfig.username}
              onInput={(e) => setBasicAuthConfig('username', e.currentTarget.value)}
            />
            <Form
              label='パスワード'
              type='password'
              value={basicAuthConfig.password}
              onInput={(e) => setBasicAuthConfig('password', e.currentTarget.value)}
            />
          </Match>
          <Match when={authMethod() === 'ssh'}>
            <Form
              label='SSH公開鍵'
              placeholder='ssh-ed25519 ******'
              value={sshAuthConfig.sshKey}
              onInput={(e) => setSshAuthConfig('sshKey', e.currentTarget.value)}
            />
            <Show when={sshAuthConfig.sshKey.length === 0}>
              <SystemPublicKey />
            </Show>
          </Match>
        </Switch>
        <Button color='black1' size='large' onclick={createRepository}>
          + Create new Repository
        </Button>
      </ContentContainer>
    </Container>
  )
}
