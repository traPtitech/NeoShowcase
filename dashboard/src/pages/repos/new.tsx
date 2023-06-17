import { JSX, JSXElement, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { CreateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/Button'
import { Header } from '/@/components/Header'
import { client } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { AuthConfig, RepositoryAuthSettings } from '/@/components/RepositoryAuthSettings'
import { InputBar, InputLabel } from '/@/components/Input'
import { NavContainer, NavTitleContainer } from '/@/components/Nav'

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
      <InputLabel>{props.label}</InputLabel>
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

export default () => {
  const navigate = useNavigate()

  // 認証情報
  // 認証方法の切り替え時に情報を保持するために、storeを使用して3種類の認証情報を保持する
  const [authConfig, setAuthConfig] = createStore<AuthConfig>({
    none: {
      case: 'none',
      value: {},
    },
    basic: {
      case: 'basic',
      value: {
        username: '',
        password: '',
      },
    },
    ssh: {
      case: 'ssh',
      value: {
        keyId: '',
      },
    },
    authMethod: 'none',
  })

  const [requestConfig, setRequestConfig] = createStore<PlainMessage<CreateRepositoryRequest>>({
    url: '',
    name: '',
    auth: undefined,
  })

  let formContainer: HTMLFormElement

  const createRepository: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    setRequestConfig('auth', { auth: authConfig[authConfig.authMethod] })
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

  // URLからリポジトリ名を自動入力
  createEffect(() => {
    const repositoryName = extractRepositoryNameFromURL(requestConfig.url)
    setRequestConfig('name', repositoryName)
  })

  return (
    <Container>
      <Header />
      <NavContainer>
        <NavTitleContainer>Create Repository</NavTitleContainer>
      </NavContainer>
      <ContentContainer>
        <InputFormContainer ref={formContainer}>
          <Form
            label='URL'
            // SSH URLはURLとしては不正なのでtypeを変更
            type={authConfig.authMethod === 'ssh' ? 'text' : 'url'}
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
          <RepositoryAuthSettings authConfig={authConfig} setAuthConfig={setAuthConfig} />
          <Button color='black1' size='large' width='auto' onclick={createRepository} type='submit'>
            + Create new Repository
          </Button>
        </InputFormContainer>
      </ContentContainer>
    </Container>
  )
}
