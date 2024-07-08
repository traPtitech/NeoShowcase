import { styled } from '@macaron-css/solid'
import {
  Field,
  Form,
  type FormStore,
  createFormStore,
  getValue,
  getValues,
  setValue,
  validate,
} from '@modular-forms/solid'
import { Title } from '@solidjs/meta'
import { A, useNavigate, useSearchParams } from '@solidjs/router'
import Fuse from 'fuse.js'
import {
  type Accessor,
  type Component,
  For,
  Match,
  type Setter,
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
  type Application,
  ApplicationConfig,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  type Repository,
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
import { AppGeneralConfig, type AppGeneralForm } from '/@/components/templates/app/AppGeneralConfig'
import {
  type BuildConfigForm,
  BuildConfigs,
  configToForm,
  formToConfig,
} from '/@/components/templates/app/BuildConfigs'
import { type WebsiteFormStatus, WebsiteSetting, newWebsite } from '/@/components/templates/app/WebsiteSettings'
import ReposFilter from '/@/components/templates/repo/ReposFilter'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { type RepositoryOrigin, originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'

const RepositoryStepContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    minHeight: '800px',
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
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
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
    rowGap: '2px',
    columnGap: '8px',
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
    background: colorVars.semantic.ui.primary,
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
      },
    },
  },
})

const RepositoryStep: Component<{
  setRepo: (repo: Repository) => void
}> = (props) => {
  const [repos] = createResource(() =>
    client.getRepositories({
      scope: GetRepositoriesRequest_Scope.CREATABLE,
    }),
  )
  const [apps] = createResource(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))

  const [query, setQuery] = createSignal('')
  const [origin, setOrigin] = createSignal<RepositoryOrigin[]>(['GitHub', 'Gitea', 'Others'])

  const filteredReposByOrigin = createMemo(() => {
    const p = origin()
    return repos()?.repositories.filter((r) => p.includes(repositoryURLToOrigin(r.url)))
  })
  const repoWithApps = createMemo(() => {
    const appsMap = apps()?.applications.reduce(
      (acc, app) => {
        if (!acc[app.repositoryId]) acc[app.repositoryId] = 0
        acc[app.repositoryId]++
        return acc
      },
      {} as { [id: Repository['id']]: number },
    )

    return (
      filteredReposByOrigin()?.map(
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
      <TextField
        placeholder="Search"
        value={query()}
        onInput={(e) => setQuery(e.currentTarget.value)}
        leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        rightIcon={<ReposFilter origin={origin()} setOrigin={setOrigin} />}
      />
      <List.Container>
        <A href="/repos/new">
          <RegisterRepositoryButton>
            <MaterialSymbols>add</MaterialSymbols>
            Register Repository
          </RegisterRepositoryButton>
        </A>
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
                  <RepositoryIcon>{originToIcon(repositoryURLToOrigin(repo.repo.url), 24)}</RepositoryIcon>
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
  createAppForm: FormStore<GeneralForm, undefined>
  backToRepoStep: () => void
  proceedToWebsiteStep: () => void
}> = (props) => {
  return (
    <Form of={props.createAppForm} onSubmit={props.proceedToWebsiteStep} style={{ width: '100%' }}>
      <FormsContainer>
        <FormContainer>
          <FormTitle>
            Create Application from
            {originToIcon(repositoryURLToOrigin(props.repo.url), 24)}
            {props.repo.name}
          </FormTitle>
          {/* 
            modular formsでは `FormStore<T extends AppGeneralForm, undefined>`のような
            genericsが使用できないためignoreしている
            */}
          {/* @ts-ignore */}
          <AppGeneralConfig repo={props.repo} formStore={props.createAppForm} editBranchId={false} hasPermission />
          {/* @ts-ignore */}
          <BuildConfigs formStore={props.createAppForm} disableEditDB={false} hasPermission />
          <Field of={props.createAppForm} name="startOnCreate" type="boolean">
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
            disabled={props.createAppForm.invalid || props.createAppForm.submitting}
            loading={props.createAppForm.submitting}
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
const AddMoreButtonContainer = styled('div', {
  base: {
    display: 'flex',
    justifyContent: 'center',
  },
})

const WebsiteStep: Component<{
  isRuntimeApp: boolean
  websiteForms: Accessor<FormStore<WebsiteFormStatus, undefined>[]>
  setWebsiteForms: Setter<FormStore<WebsiteFormStatus, undefined>[]>
  backToGeneralStep: () => void
  submit: () => Promise<void>
}> = (props) => {
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const addWebsiteForm = () => {
    const form = createFormStore<WebsiteFormStatus>({
      initialValues: {
        state: 'added',
        website: newWebsite(),
      },
    })
    props.setWebsiteForms((prev) => prev.concat([form]))
  }

  const handleSubmit = async () => {
    try {
      const isValid = (await Promise.all(props.websiteForms().map((form) => validate(form)))).every((v) => v)
      if (!isValid) return
      setIsSubmitting(true)
      await props.submit()
    } catch (err) {
      console.error(err);
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Show when={systemInfo()}>
      <FormsContainer>
        <DomainsContainer>
          <For
            each={props.websiteForms()}
            fallback={
              <List.PlaceHolder>
                <MaterialSymbols displaySize={80}>link_off</MaterialSymbols>
                URLが設定されていません
                <Button
                  variants="primary"
                  size="medium"
                  rightIcon={<MaterialSymbols>add</MaterialSymbols>}
                  onClick={addWebsiteForm}
                  type="button"
                >
                  Add URL
                </Button>
              </List.PlaceHolder>
            }
          >
            {(form, i) => (
              <WebsiteSetting
                isRuntimeApp={props.isRuntimeApp}
                formStore={form}
                deleteWebsite={() => props.setWebsiteForms((prev) => [...prev.slice(0, i()), ...prev.slice(i() + 1)])}
                hasPermission
              />
            )}
          </For>
          <Show when={props.websiteForms().length > 0}>
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
            leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}
            onClick={props.backToGeneralStep}
          >
            Back
          </Button>
          <Button
            size="medium"
            variants="primary"
            onClick={handleSubmit}
            disabled={isSubmitting()}
          // TODO: hostが空の状態でsubmitして一度requiredエラーが出たあとhostを入力してもエラーが消えない
          // disabled={props.websiteForms().some((form) => form.invalid)}
          >
            Create Application
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
type FormStep = (typeof formStep)[keyof typeof formStep]

export default () => {
  const [searchParams, setParam] = useSearchParams()
  const [currentStep, setCurrentStep] = createSignal<FormStep>(formStep.repository)

  const [repo, { mutate: mutateRepo }] = createResource(
    () => searchParams.repositoryID,
    (id) => client.getRepository({ repositoryId: id }),
  )

  // このページに遷移した時にURLパラメータにrepositoryIDがあれば
  // generalStepに遷移する
  onMount(() => {
    if (searchParams.repositoryID !== undefined) {
      setCurrentStep(formStep.general)
    }
  })

  const createAppForm = createFormStore<GeneralForm>({
    initialValues: {
      name: '',
      refName: '',
      repositoryId: repo()?.id,
      startOnCreate: false,
      ...configToForm(new ApplicationConfig()),
    },
  })
  const isRuntimeApp = () => {
    return (
      getValue(createAppForm, 'case') === 'runtimeBuildpack' ||
      getValue(createAppForm, 'case') === 'runtimeCmd' ||
      getValue(createAppForm, 'case') === 'runtimeDockerfile'
    )
  }
  // repo更新時にcreateAppFormのrepositoryIdを更新する
  createEffect(() => {
    setValue(createAppForm, 'repositoryId', repo()?.id)
  })

  const [websiteForms, setWebsiteForms] = createSignal<FormStore<WebsiteFormStatus, undefined>[]>([])

  // TODO: ブラウザバック時のrepositoryIDの設定

  // repositoryが指定されたらビルド設定に進む
  createEffect(() => {
    if (repo() !== undefined) {
      setParam({ repositoryID: repo()?.id })
      GoToGeneralStep()
    }
  })

  const backToRepoStep = () => {
    setCurrentStep(formStep.repository)
    // 選択していたリポジトリをリセットする
    setParam({ repositoryID: undefined })
  }
  const GoToGeneralStep = () => {
    setCurrentStep(formStep.general)
  }
  const GoToWebsiteStep = () => {
    setCurrentStep(formStep.website)
  }

  const createApp = async (): Promise<Application> => {
    const values = getValues(createAppForm, { shouldActive: false })
    const websitesToSave = websiteForms()
      .map((form) => getValues(form).website)
      .filter((w): w is Exclude<typeof w, undefined> => w !== undefined)

    const createdApp = await client.createApplication({
      name: values.name,
      refName: values.refName,
      repositoryId: values.repositoryId,
      config: {
        buildConfig: formToConfig({
          case: values.case,
          config: values.config as BuildConfigs,
        }),
      },
      websites: websitesToSave,
      startOnCreate: values.startOnCreate,
    })
    return createdApp
  }

  const navigate = useNavigate()
  const submit = async () => {
    try {
      const createdApp = await createApp()
      toast.success('アプリケーションを登録しました')
      navigate(`/apps/${createdApp.id}`)
    } catch (e) {
      handleAPIError(e, 'アプリケーションの登録に失敗しました')
    }
  }

  return (
    <WithNav.Container>
      <Title>Create Application - NeoShowcase</Title>
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
                    title: '1. Repository',
                    description: 'リポジトリの選択',
                    step: formStep.repository,
                  },
                  {
                    title: '2. Build Settings',
                    description: 'ビルド設定',
                    step: formStep.general,
                  },
                  {
                    title: '3. URLs',
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
                      createAppForm={createAppForm}
                      proceedToWebsiteStep={GoToWebsiteStep}
                    />
                  )}
                </Show>
              </Match>
              <Match when={currentStep() === formStep.website}>
                <WebsiteStep
                  isRuntimeApp={isRuntimeApp()}
                  backToGeneralStep={GoToGeneralStep}
                  websiteForms={websiteForms}
                  setWebsiteForms={setWebsiteForms}
                  submit={submit}
                />
              </Match>
            </Switch>
          </StepsContainer>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
