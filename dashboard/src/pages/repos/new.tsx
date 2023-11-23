import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import { AuthForm, RepositoryAuthSettings, formToAuth } from '/@/components/templates/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { SubmitHandler, createForm, getValue, required, setValue } from '@modular-forms/solid'
import { useNavigate, useSearchParams } from '@solidjs/router'
import { createEffect } from 'solid-js'
import toast from 'solid-toast'
import { TextField } from '../../components/UI/TextField'

const Container = styled('div', {
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

type Config = AuthForm & {
  name: string
}

export default () => {
  const navigate = useNavigate()
  const [params] = useSearchParams()

  const [config, Form] = createForm<Config>({
    initialValues: {
      url: '',
      name: '',
      case: 'none',
      auth: {
        basic: {
          username: '',
          password: '',
        },
        ssh: {
          keyId: '',
        },
      },
    },
  })
  const handleSubmit: SubmitHandler<Config> = async (values) => {
    try {
      const res = await client.createRepository({
        name: values.name,
        url: values.url,
        auth: {
          auth: formToAuth(values),
        },
      })
      toast.success('リポジトリを登録しました')

      if (params.newApp === 'true') {
        // アプリ作成ページのリポジトリ登録ボタンから来た場合は新規アプリ作成ページに遷移
        navigate(`/apps/new?repositoryID=${res.id}`)
      } else {
        // リポジトリページに遷移
        navigate(`/repos/${res.id}`)
      }
    } catch (e) {
      return handleAPIError(e, 'リポジトリの登録に失敗しました')
    }
  }

  // URLからリポジトリ名を自動入力
  createEffect(() => {
    const repositoryName = extractRepositoryNameFromURL(getValue(config, 'url') ?? '')
    setValue(config, 'name', repositoryName)
  })

  const AuthSetting = RepositoryAuthSettings({
    // @ts-ignore
    formStore: config,
    hasPermission: true,
  })

  return (
    <WithNav.Container>
      <WithNav.Navs>
        <Nav title="Register Repository" backTo="/apps" backToTitle="Apps" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <Form.Form onSubmit={handleSubmit}>
            <Container>
              <InputsContainer>
                <AuthSetting.Url />
                <Form.Field name="name" validate={required('Enter Repository Name')}>
                  {(field, fieldProps) => (
                    <TextField
                      label="Repository Name"
                      required
                      {...fieldProps}
                      value={field.value ?? ''}
                      error={field.error}
                    />
                  )}
                </Form.Field>
                <AuthSetting.AuthMethod />
                <AuthSetting.AuthConfig />
              </InputsContainer>
              <Button
                variants="primary"
                size="medium"
                rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
                type="submit"
                disabled={config.invalid || config.submitting}
              >
                Register
              </Button>
            </Container>
          </Form.Form>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
