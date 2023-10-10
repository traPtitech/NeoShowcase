import { CreateRepositoryAuth, CreateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextInput } from '/@/components/UI/TextInput'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { FormItem } from '/@/components/templates/FormItem'
import { Nav } from '/@/components/templates/Nav'
import { RepositoryAuthSettings } from '/@/components/templates/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { colorVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { JSX, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const Container = styled('form', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-end',
    gap: '40px',
  },
})
const InputsContainer = styled('div', {
  base: {
    width: '100%',
    margin: '0 auto',
    padding: '20px 24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',

    borderRadius: '8px',
    background: colorVars.semantic.ui.primary,
  },
})

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

  let formRef: HTMLFormElement

  const createRepository: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formRef.reportValidity()) {
      return
    }

    try {
      const res = await client.createRepository({
        ...requestConfig,
        auth: authConfig,
      })
      toast.success('リポジトリを登録しました')
      // リポジトリページに遷移
      navigate(`/repos/${res.id}`)
    } catch (e) {
      return handleAPIError(e, 'リポジトリの登録に失敗しました')
    }
  }

  // URLからリポジトリ名を自動入力
  createEffect(() => {
    const repositoryName = extractRepositoryNameFromURL(requestConfig.url)
    setRequestConfig('name', repositoryName)
  })

  const Auth = RepositoryAuthSettings({
    url: requestConfig.url,
    setUrl: (v) => setRequestConfig('url', v),
    authConfig: authConfig,
    setAuthConfig: setAuthConfig,
  })

  return (
    <WithNav.Container>
      <WithNav.Navs>
        <Nav title="Register Repository" backToTitle="Back" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <Container ref={formRef}>
            <InputsContainer>
              <Auth.Url />
              <FormItem title="Repository Name" required>
                <TextInput
                  value={requestConfig.name}
                  onInput={(e) => setRequestConfig('name', e.currentTarget.value)}
                />
              </FormItem>
              <Auth.AuthMethod />
              <Auth.AuthConfig />
            </InputsContainer>
            <Button
              color="primary"
              size="medium"
              onClick={createRepository}
              rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
            >
              Register
            </Button>
          </Container>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
