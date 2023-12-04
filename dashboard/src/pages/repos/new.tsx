import { styled } from '@macaron-css/solid'
import { SubmitHandler, createForm, getValue, required, setValue } from '@modular-forms/solid'
import { Title } from '@solidjs/meta'
import { useNavigate } from '@solidjs/router'
import { createEffect } from 'solid-js'
import toast from 'solid-toast'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextField } from '/@/components/UI/TextField'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import { AuthForm, RepositoryAuthSettings, formToAuth } from '/@/components/templates/repo/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { colorVars } from '/@/theme'

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
      // 新規アプリ作成ページに遷移
      navigate(`/apps/new?repositoryID=${res.id}`)
    } catch (e) {
      return handleAPIError(e, 'リポジトリの登録に失敗しました')
    }
  }

  // URLからリポジトリ名, 認証方法を自動入力
  createEffect(() => {
    const url = getValue(config, 'url')
    if (url === undefined) return
    const repositoryName = extractRepositoryNameFromURL(url)
    setValue(config, 'name', repositoryName)

    !url.startsWith('http') ? setValue(config, 'case', 'ssh') : setValue(config, 'case', 'none')
  })

  const AuthSetting = RepositoryAuthSettings({
    // @ts-ignore
    formStore: config,
    hasPermission: true,
  })

  return (
    <WithNav.Container>
      <Title>Register Repository - NeoShowcase</Title>
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
                loading={config.submitting}
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
