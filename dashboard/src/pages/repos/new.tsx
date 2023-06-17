import { JSX, JSXElement, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { CreateRepositoryAuth, CreateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/Button'
import { Header } from '/@/components/Header'
import { client } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { RepositoryAuthSettings } from '/@/components/RepositoryAuthSettings'
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

  const [requestConfig, setRequestConfig] = createStore<PlainMessage<CreateRepositoryRequest>>({
    url: '',
    name: '',
    auth: undefined,
  })
  const [authConfig, setAuthConfig] = createStore<PlainMessage<CreateRepositoryAuth>>({
    auth: { case: 'none', value: {} },
  })

  let formContainer: HTMLFormElement

  const createRepository: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    try {
      const res = await client.createRepository({ ...requestConfig, auth: authConfig })
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
            type={authConfig.auth.case === 'ssh' ? 'text' : 'url'}
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
