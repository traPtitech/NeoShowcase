import { styled } from '@macaron-css/solid'
import { Field, Form, FormStore, SubmitHandler, createFormStore, getValues, validate } from '@modular-forms/solid'
import { A, useNavigate, useSearchParams } from '@solidjs/router'
import Fuse from 'fuse.js'
import {
  Component,
  For,
  Match,
  Show,
  Switch,
  createEffect,
  createMemo,
  createResource,
  createSignal,
  onMount,
} from 'solid-js'
import toast from 'solid-toast'
import {
  Application,
  ApplicationConfig,
  DeployType,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { Progress } from '/@/components/UI/StepProgress'
import { TextField } from '/@/components/UI/TextField'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { List } from '/@/components/templates/List'
import { Nav } from '/@/components/templates/Nav'
import { MultiSelect } from '/@/components/templates/Select'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { Provider, providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { colorVars, textVars } from '/@/theme'
import { AppGeneralConfig, AppGeneralForm } from '../../components/templates/app/AppGeneralConfig'
import { BuildConfigForm, BuildConfigs, configToForm, formToConfig } from '../../components/templates/app/BuildConfigs'
import { WebsiteSetting, newWebsite } from '../../components/templates/app/WebsiteSettings'

const RepositoryStepContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'hidden',
    padding: '24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
  },
})
const RepositoryListContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
  },
})
const RepositoryButton = styled('button', {
  base: {
    width: '100%',
    background: colorVars.semantic.ui.primary,
    border: 'none',
    cursor: 'pointer',

    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:not(:last-child)': {
        borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const RepositoryRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px',
    display: 'grid',
    gridTemplateColumns: '24px auto 1fr auto',
    gridTemplateRows: 'auto auto',
    gridTemplateAreas: `
      "icon name count button"
      ". url url button"`,
    gap: '8px',
    textAlign: 'left',
  },
})
const RepositoryIcon = styled('div', {
  base: {
    gridArea: 'icon',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    flexShrink: 0,
  },
})
const RepositoryName = styled('div', {
  base: {
    width: '100%',
    gridArea: 'name',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.black,
    ...textVars.h4.bold,
  },
})
const AppCount = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const RepositoryUrl = styled('div', {
  base: {
    gridArea: 'url',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const CreateAppText = styled('div', {
  base: {
    gridArea: 'button',
    display: 'flex',
    justifyContent: 'flex-end',
    alignItems: 'center',
    gap: '4px',
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const RegisterRepositoryButton = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    padding: '20px',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,

    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
    },
  },
})
const SortContainer = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '2fr 1fr',
    gap: '16px',
  },
})

const RepositoryStep: Component<{
  setRepo: (repo: Repository) => void
}> = (props) => {
  const [repos] = createResource(() =>
    client.getRepositories({
      scope: GetRepositoriesRequest_Scope.MINE,
    }),
  )
  const [apps] = createResource(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))

  const [query, setQuery] = createSignal('')
  const [provider, setProvider] = createSignal<Provider[]>(['GitHub', 'GitLab', 'Gitea'])

  const filteredReposByProvider = createMemo(() => {
    const p = provider()
    return repos()?.repositories.filter((r) => p.includes(repositoryURLToProvider(r.url)))
  })
  const repoWithApps = createMemo(() => {
    const appsMap = apps()?.applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = 0
      acc[app.repositoryId]++
      return acc
    }, {} as { [id: Repository['id']]: number })

    return (
      filteredReposByProvider()?.map(
        (
          repo,
        ): {
          repo: Repository
          appCount: number
        } => ({ repo, appCount: appsMap?.[repo.id] ?? 0 }),
      ) ?? []
    )
  })

  const fuse = createMemo(
    () =>
      new Fuse(repoWithApps(), {
        keys: ['repo.name', 'repo.htmlUrl'],
      }),
  )
  const filteredRepos = createMemo(() => {
    if (query() === '') return repoWithApps()
    return fuse()
      .search(query())
      .map((r) => r.item)
  })

  return (
    <RepositoryStepContainer>
      <SortContainer>
        <TextField
          placeholder="Search"
          value={query()}
          onInput={(e) => setQuery(e.currentTarget.value)}
          leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        />
        <MultiSelect
          options={[
            {
              label: 'GitHub',
              value: 'GitHub',
            },
            {
              label: 'GitLab',
              value: 'GitLab',
            },
            {
              label: 'Gitea',
              value: 'Gitea',
            },
          ]}
          placeholder="Provider"
          value={provider()}
          setValue={setProvider}
        />
      </SortContainer>
      <List.Container>
        <RepositoryListContainer>
          <For
            each={filteredRepos()}
            fallback={
              <List.Row>
                <List.RowContent>
                  <List.RowData>Repository Not Found</List.RowData>
                </List.RowContent>
              </List.Row>
            }
          >
            {(repo) => (
              <RepositoryButton
                onClick={() => {
                  props.setRepo(repo.repo)
                }}
                type="button"
              >
                <RepositoryRow>
                  <RepositoryIcon>{providerToIcon(repositoryURLToProvider(repo.repo.url), 24)}</RepositoryIcon>
                  <RepositoryName>{repo.repo.name}</RepositoryName>
                  <AppCount>{repo.appCount > 0 && `${repo.appCount} apps`}</AppCount>
                  <RepositoryUrl>{repo.repo.htmlUrl}</RepositoryUrl>
                  <CreateAppText>
                    Create App
                    <MaterialSymbols>arrow_forward</MaterialSymbols>
                  </CreateAppText>
                </RepositoryRow>
              </RepositoryButton>
            )}
          </For>
        </RepositoryListContainer>
        <A href="/repos/new?newApp=true">
          <RegisterRepositoryButton>
            <MaterialSymbols>add</MaterialSymbols>
            Register Repository
          </RegisterRepositoryButton>
        </A>
      </List.Container>
    </RepositoryStepContainer>
  )
}

const FormsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '40px',
  },
})
const FormContainer = styled('div', {
  base: {
    width: '100%',
    padding: '24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
  },
})
const FormTitle = styled('h2', {
  base: {
    display: 'flex',
    alignItems: 'center',
    gap: '4px',
    overflowWrap: 'anywhere',
    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
  },
})
const ButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    gap: '20px',
  },
})

type GeneralForm = AppGeneralForm & BuildConfigForm & { startOnCreate: boolean }

const GeneralStep: Component<{
  repo: Repository
  backToRepoStep: () => void
  setAppId: (appId: string) => void
}> = (props) => {
  const form = createFormStore<GeneralForm>({
    initialValues: {
      name: '',
      refName: '',
      repositoryId: props.repo.id,
      startOnCreate: false,
      ...configToForm(new ApplicationConfig()),
    },
  })

  const handleSubmit: SubmitHandler<GeneralForm> = async (values) => {
    try {
      const createdApp = await client.createApplication({
        name: values.name,
        refName: values.refName,
        repositoryId: values.repositoryId,
        config: {
          buildConfig: formToConfig({
            case: values.case,
            config: values.config,
          }),
        },
        startOnCreate: values.startOnCreate,
      })
      toast.success('アプリケーションを登録しました')
      props.setAppId(createdApp.id)
    } catch (e) {
      return handleAPIError(e, 'アプリケーションの登録に失敗しました')
    }
  }

  return (
    <Form of={form} onSubmit={handleSubmit} style={{ width: '100%' }}>
      <FormsContainer>
        <FormContainer>
          <FormTitle>
            create application from
            {providerToIcon(repositoryURLToProvider(props.repo.url), 24)}
            {props.repo.name}
          </FormTitle>
          {/* 
            modular formsでは `FormStore<T extends AppGeneralForm, undefined>`のような
            genericsが使用できないためignoreしている
            */}
          {/* @ts-ignore */}
          <AppGeneralConfig repo={props.repo} formStore={form} editBranchId={false} hasPermission />
          {/* @ts-ignore */}
          <BuildConfigs formStore={form} disableEditDB={false} hasPermission />
          <Field of={form} name="startOnCreate" type="boolean">
            {(field, fieldProps) => (
              <FormItem
                title="Start Immediately"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>この設定で今すぐ起動するかどうか</div>
                        <div>(環境変数はアプリ作成後設定可能になります)</div>
                      </>
                    ),
                  },
                }}
              >
                <CheckBox.Option
                  {...fieldProps}
                  label="今すぐ起動する"
                  checked={field.value ?? false}
                  error={field.error}
                />
              </FormItem>
            )}
          </Field>
        </FormContainer>
        <ButtonsContainer>
          <Button
            size="medium"
            variants="border"
            onClick={props.backToRepoStep}
            leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}
          >
            Back
          </Button>
          <Button
            size="medium"
            variants="primary"
            type="submit"
            rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
            disabled={form.invalid || form.submitting}
            loading={form.submitting}
          >
            Next
          </Button>
        </ButtonsContainer>
      </FormsContainer>
    </Form>
  )
}

const DomainsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '24px',
  },
})
const PlaceHolder = styled('div', {
  base: {
    width: '100%',
    height: '400px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})
const AddMoreButtonContainer = styled('div', {
  base: {
    display: 'flex',
    justifyContent: 'center',
  },
})

const WebsiteStep: Component<{
  app: Application
}> = (props) => {
  const [websiteForms, setWebsiteForms] = createSignal<FormStore<WebsiteSetting, undefined>[]>([])

  const navigate = useNavigate()
  const skipWebsiteConfig = () => {
    navigate(`/apps/${props.app.id}`)
  }

  const addWebsiteForm = () => {
    const form = createFormStore<WebsiteSetting>({
      initialValues: {
        state: 'added',
        website: newWebsite(),
      },
    })
    setWebsiteForms((prev) => prev.concat([form]))
  }

  const saveWebsiteConfig = async () => {
    try {
      const isValid = (await Promise.all(websiteForms().map((form) => validate(form)))).every((v) => v)
      if (!isValid) return
      const websitesToSave = websiteForms()
        .map((form) => getValues(form).website)
        .filter((w): w is Exclude<typeof w, undefined> => w !== undefined)
      await client.updateApplication({
        id: props.app.id,
        websites: {
          websites: websitesToSave,
        },
      })
      toast.success('ウェブサイト設定を保存しました')
      navigate(`/apps/${props.app.id}`)
    } catch (e) {
      handleAPIError(e, 'Failed to save website settings')
    }
  }

  return (
    <Show when={systemInfo()}>
      <FormsContainer>
        <DomainsContainer>
          <For
            each={websiteForms()}
            fallback={
              <PlaceHolder>
                <MaterialSymbols displaySize={80}>link_off</MaterialSymbols>
                No Websites Configured
                <Button
                  variants="primary"
                  size="medium"
                  rightIcon={<MaterialSymbols>add</MaterialSymbols>}
                  onClick={addWebsiteForm}
                  type="button"
                >
                  Add Website
                </Button>
              </PlaceHolder>
            }
          >
            {(form, i) => (
              <WebsiteSetting
                isRuntimeApp={props.app.deployType === DeployType.RUNTIME}
                formStore={form}
                deleteWebsite={() => setWebsiteForms((prev) => [...prev.slice(0, i()), ...prev.slice(i() + 1)])}
                hasPermission
              />
            )}
          </For>
          <Show when={websiteForms().length > 0}>
            <AddMoreButtonContainer>
              <Button
                onclick={addWebsiteForm}
                variants="border"
                size="small"
                leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
                type="button"
              >
                Add More
              </Button>
            </AddMoreButtonContainer>
          </Show>
        </DomainsContainer>
        <ButtonsContainer>
          <Button
            size="medium"
            variants="ghost"
            rightIcon={<MaterialSymbols>skip_next</MaterialSymbols>}
            onClick={skipWebsiteConfig}
          >
            Skip
          </Button>
          <Button size="medium" variants="primary" onClick={saveWebsiteConfig}>
            Save Website Config
          </Button>
        </ButtonsContainer>
      </FormsContainer>
    </Show>
  )
}

const StepsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '40px',
  },
  variants: {
    fit: {
      true: {
        maxHeight: '100%',
      },
    },
  },
})

const formStep = {
  repository: 0,
  general: 1,
  website: 2,
} as const
type FormStep = typeof formStep[keyof typeof formStep]

export default () => {
  const [searchParams, setParam] = useSearchParams()
  const [currentStep, setCurrentStep] = createSignal<FormStep>(formStep.repository)
  const [appId, setAppId] = createSignal<string | undefined>(undefined)

  const [repo, { mutate: mutateRepo }] = createResource(
    () => searchParams.repositoryID,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [app] = createResource(
    () => appId(),
    (id) => client.getApplication({ id }),
  )

  // repositoryIDが指定されている場合はビルド設定に進む
  onMount(() => {
    if (searchParams.repositoryID !== undefined) {
      setCurrentStep(formStep.general)
    }
  })

  // repositoryが指定されたらビルド設定に進む
  createEffect(() => {
    if (repo() !== undefined) {
      setParam({ repositoryID: repo()?.id })
      setCurrentStep(formStep.general)
    }
  })

  const backToRepoStep = () => {
    setCurrentStep(formStep.repository)
    setParam({ repositoryID: undefined })
    mutateRepo(undefined)
  }

  createEffect(() => {
    if (app() !== undefined) {
      setCurrentStep(formStep.website)
      setParam({ appID: app()?.id })
    }
  })

  return (
    <WithNav.Container>
      <WithNav.Navs>
        <Nav title="Create Application" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <StepsContainer fit={currentStep() === formStep.repository}>
            <Progress.Container>
              <For
                each={[
                  {
                    title: 'Repository',
                    description: 'リポジトリの選択',
                    step: formStep.repository,
                  },
                  {
                    title: 'Build Settings',
                    description: 'ビルド設定',
                    step: formStep.general,
                  },
                  {
                    title: 'Domains',
                    description: 'アクセスURLの設定',
                    step: formStep.website,
                  },
                ]}
              >
                {(step) => (
                  <Progress.Step
                    title={step.title}
                    description={step.description}
                    state={
                      currentStep() === step.step ? 'current' : currentStep() > step.step ? 'complete' : 'incomplete'
                    }
                  />
                )}
              </For>
            </Progress.Container>
            <Switch>
              <Match when={currentStep() === formStep.repository}>
                <RepositoryStep setRepo={(repo) => mutateRepo(repo)} />
              </Match>
              <Match when={currentStep() === formStep.general}>
                <Show when={repo()}>
                  {(nonNullRepo) => (
                    <GeneralStep
                      repo={nonNullRepo()}
                      backToRepoStep={backToRepoStep}
                      setAppId={(appId) => {
                        setAppId(appId)
                      }}
                    />
                  )}
                </Show>
              </Match>
              <Match when={currentStep() === formStep.website}>
                <Show when={app()}>{(nonNullApp) => <WebsiteStep app={nonNullApp()} />}</Show>
              </Match>
            </Switch>
          </StepsContainer>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
