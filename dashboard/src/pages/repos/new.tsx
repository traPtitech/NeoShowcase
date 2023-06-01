import { JSX, JSXElement, Match, Show, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { Empty } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import {
  CreateRepositoryAuth,
  CreateRepositoryAuthBasic,
  CreateRepositoryAuthSSH,
  CreateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/Button'
import { Header } from '/@/components/Header'
import { Radio } from '/@/components/Radio'
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

// copy from /pages/apps/new
const InputFormContainer = styled('form', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})
const InputForm = styled('div', {
  base: {},
})
const InputFormText = styled('div', {
  base: {
    fontSize: '16px',
    alignItems: 'center',
    fontWeight: 700,
    color: vars.text.black1,

    marginBottom: '4px',
  },
})
const InputBar = styled('input', {
  base: {
    padding: '8px 12px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    fontSize: '14px',
    marginLeft: '4px',

    width: '320px',

    display: 'flex',
    flexDirection: 'column',

    '::placeholder': {
      color: vars.text.black3,
    },
  },
})

interface FormProps {
  label: string
  type?: JSX.InputHTMLAttributes<HTMLInputElement>['type']
  placeholder?: JSX.InputHTMLAttributes<HTMLInputElement>['placeholder']
  value: JSX.InputHTMLAttributes<HTMLInputElement>['value']
  onInput: JSX.InputHTMLAttributes<HTMLInputElement>['onInput']
  required?: JSX.InputHTMLAttributes<HTMLInputElement>['required']
}

const Form = (props: FormProps): JSXElement => {
  return (
    <InputForm>
      <InputFormText>{props.label}</InputFormText>
      <InputBar
        type={props.type ?? 'text'}
        placeholder={props.placeholder ?? ''}
        value={props.value}
        onInput={props.onInput}
        required={props.required}
      />
    </InputForm>
  )
}

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
    fontFamily: 'monospace',
    fontSize: '14px',
    background: vars.bg.white2,
    color: vars.text.black1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
  },
})

export default () => {
  const navigate = useNavigate()
  // 認証方法 ("none" | "ssh" | "basic")
  type AuthMethod = CreateRepositoryAuth['auth']['case']
  const [authMethod, setAuthMethod] = createSignal<AuthMethod>('none')

  // 認証情報
  // 認証方法の切り替え時に情報を保持するために、storeを使用して3種類の認証情報を保持する
  const [authConfig, setAuthConfig] = createStore<{
    [K in AuthMethod]: Extract<CreateRepositoryAuth['auth'], { case: K }>
  }>({
    none: {
      case: 'none',
      value: new Empty(),
    },
    basic: {
      case: 'basic',
      value: new CreateRepositoryAuthBasic(),
    },
    ssh: {
      case: 'ssh',
      value: new CreateRepositoryAuthSSH(),
    },
  })

  const [requestConfig, setRequestConfig] = createStore(
    new CreateRepositoryRequest({
      auth: new CreateRepositoryAuth(),
    }),
  )

  let formContainer: HTMLFormElement

  const createRepository: JSX.EventHandler<HTMLInputElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (formContainer.reportValidity()) {
      setRequestConfig('auth', 'auth', authConfig[authMethod()])
      try {
        const res = await client.createRepository(requestConfig)
        toast.success('リポジトリを登録しました')
        // リポジトリページに遷移
        navigate(`/repos/${res.id}`)
      } catch (e) {
        console.error(e)
        // gRPCエラー
        if (e instanceof ConnectError) {
          toast.error('リポジトリの登録に失敗しました\n' + e.message)
        }
      }
    }
  }

  // URLからリポジトリ名を自動入力
  createEffect(() => {
    const segments = requestConfig.url.split('/')
    const lastSegment = segments.pop() || segments.pop() // 末尾のスラッシュを除去
    const repositoryName = lastSegment?.replace(/\.git$/, '') ?? ''
    setRequestConfig('name', repositoryName)
  })

  const [systemPublicKey] = createResource(() => client.getSystemPublicKey({}))

  return (
    <Container>
      <Header />
      <PageTitle>Create Repository</PageTitle>
      <ContentContainer>
        <InputFormContainer ref={formContainer}>
          <Form
            label='URL'
            // SSH URLはURLとしては不正なのでtypeを変更
            type={authMethod() === 'ssh' ? 'text' : 'url'}
            placeholder='https://example.com/my-app.git'
            value={requestConfig.url}
            onInput={(e) => setRequestConfig('url', e.currentTarget.value)}
            required
          />
          <Form
            label='リポジトリ名'
            placeholder='my-app'
            value={requestConfig.name}
            onInput={(e) => setRequestConfig('name', e.currentTarget.value)}
            required
          />
          <InputForm>
            <InputFormText>認証方法</InputFormText>
            <Radio
              items={[
                { title: '認証を使用しない', value: 'none' },
                { title: 'Basic認証を使用', value: 'basic' },
                { title: 'SSH認証を使用', value: 'ssh' },
              ]}
              selected={authMethod()}
              setSelected={setAuthMethod}
            />
          </InputForm>
          <Switch>
            <Match when={authMethod() === 'basic'}>
              <Form
                label='ユーザー名'
                value={authConfig.basic.value.username}
                onInput={(e) => setAuthConfig('basic', 'value', 'username', e.currentTarget.value)}
              />
              <Form
                label='パスワード'
                type='password'
                value={authConfig.basic.value.password}
                onInput={(e) => setAuthConfig('basic', 'value', 'password', e.currentTarget.value)}
              />
            </Match>
            <Match when={authMethod() === 'ssh'}>
              <Form
                label='SSH秘密鍵'
                value={authConfig.ssh.value.sshKey}
                onInput={(e) => setAuthConfig('ssh', 'value', 'sshKey', e.currentTarget.value)}
              />
              <Show when={authConfig.ssh.value.sshKey.length === 0}>
                <div>
                  <SshDetails>
                    秘密鍵を入力せずにSSH認証でリポジトリを登録する場合、以下のSSH公開鍵が認証に使用されます。
                  </SshDetails>
                  <Switch>
                    <Match when={systemPublicKey.loading}>
                      <div>Loading...</div>
                    </Match>
                    <Match when={systemPublicKey()}>
                      <PublicKeyCode>{systemPublicKey().publicKey}</PublicKeyCode>
                    </Match>
                  </Switch>
                </div>
              </Show>
            </Match>
          </Switch>
          <Button color='black1' size='large' onclick={createRepository} type='submit'>
            + Create new Repository
          </Button>
        </InputFormContainer>
      </ContentContainer>
    </Container>
  )
}
