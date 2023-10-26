import {
  Application,
  ApplicationConfig,
  CreateApplicationRequest,
  DeployType,
  Repository,
  RuntimeConfig,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { Progress } from '/@/components/UI/StepProgress'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { BuildConfigs } from '/@/components/templates/BuildConfigs'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { GeneralConfig } from '/@/components/templates/GeneralConfig'
import { List } from '/@/components/templates/List'
import { Nav } from '/@/components/templates/Nav'
import { WebsiteSetting, newWebsite } from '/@/components/templates/WebsiteSettings'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate, useSearchParams } from '@solidjs/router'
import { Component, For, JSX, Match, Show, Switch, createResource, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '40px',
  },
})
const FormContainer = styled('form', {
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

const GeneralStep: Component<{
  repo: Repository
  gotoNextStep: (appId: string) => void
}> = (props) => {
  let formContainer: HTMLFormElement

  const [buildConfig, setBuildConfig] = createStore<PlainMessage<ApplicationConfig>['buildConfig']>({
    case: 'runtimeBuildpack',
    value: {
      context: '',
      runtimeConfig: structuredClone(new RuntimeConfig()),
    },
  })
  const [request, setRequest] = createStore<PlainMessage<CreateApplicationRequest>>({
    name: '',
    refName: '',
    repositoryId: props.repo.id,
    config: { buildConfig },
    websites: [],
    portPublications: [],
    startOnCreate: false,
  })

  const createApplication: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    try {
      const createdApp = await client.createApplication(request)
      toast.success('アプリケーションを登録しました')
      props.gotoNextStep(createdApp.id)
    } catch (e) {
      return handleAPIError(e, 'アプリケーションの登録に失敗しました')
    }
  }

  return (
    <Container>
      <FormContainer ref={formContainer}>
        <GeneralConfig repo={props.repo} config={request} setConfig={setRequest} />
        <BuildConfigs buildConfig={buildConfig} setBuildConfig={setBuildConfig} disableEditDB={false} />
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
            checked={request.startOnCreate}
            setChecked={(checked) => setRequest('startOnCreate', checked)}
            title="今すぐ起動する"
          />
        </FormItem>
      </FormContainer>
      <Button
        size="medium"
        color="primary"
        rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
        onClick={createApplication}
      >
        Next
      </Button>
    </Container>
  )
}

const DomainsContainer = styled('form', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '24px',
  },
})
const ButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    gap: '20px',
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

    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})

const WebsiteStep: Component<{
  app: Application
}> = (props) => {
  const [websiteConfigs, setWebsiteConfigs] = createStore<WebsiteSetting[]>([])

  const navigate = useNavigate()
  const skipWebsiteConfig = () => {
    navigate(`/apps/${props.app.id}`)
  }

  const addWebsite = () =>
    setWebsiteConfigs([
      ...websiteConfigs,
      {
        state: 'added',
        website: newWebsite(),
      },
    ])

  const saveWebsiteConfig = async () => {
    try {
      const websitesToSave = websiteConfigs.map((website) => website.website)
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
      <Container>
        <DomainsContainer>
          <For
            each={websiteConfigs}
            fallback={
              <List.Container>
                <PlaceHolder>
                  <MaterialSymbols displaySize={80}>link_off</MaterialSymbols>
                  No Websites Configured
                  <Button
                    color="primary"
                    size="medium"
                    rightIcon={<MaterialSymbols>add</MaterialSymbols>}
                    onClick={addWebsite}
                  >
                    Add Website
                  </Button>
                </PlaceHolder>
              </List.Container>
            }
          >
            {(config, i) => (
              <WebsiteSetting
                isRuntimeApp={props.app.deployType === DeployType.RUNTIME}
                state={config.state}
                website={config.website}
                setWebsite={(valueName, value) => {
                  setWebsiteConfigs(i(), 'website', valueName, value)
                }}
                deleteWebsite={() =>
                  setWebsiteConfigs((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
                }
              />
            )}
          </For>
          <Show when={websiteConfigs.length > 0}>
            <Button
              onclick={addWebsite}
              color="border"
              size="small"
              type="button"
              leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
            >
              Add More
            </Button>
          </Show>
        </DomainsContainer>
        <ButtonsContainer>
          <Button
            size="medium"
            color="ghost"
            rightIcon={<MaterialSymbols>skip_next</MaterialSymbols>}
            onClick={skipWebsiteConfig}
          >
            Skip
          </Button>
          <Button size="medium" color="primary" onClick={saveWebsiteConfig}>
            Save Website Config
          </Button>
        </ButtonsContainer>
      </Container>
    </Show>
  )
}

const formStep = {
  general: 0,
  website: 1,
} as const
type FormStep = typeof formStep[keyof typeof formStep]

export default () => {
  const [searchParams] = useSearchParams()
  const [currentStep, setCurrentStep] = createSignal<FormStep>(formStep.general)
  const [appId, setAppId] = createSignal<string | undefined>(undefined)

  const [repo] = createResource(
    () => searchParams.repositoryID,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [app] = createResource(
    () => appId(),
    (id) => client.getApplication({ id }),
  )

  return (
    <WithNav.Container>
      <WithNav.Navs>
        <Nav title="Create Application" backToTitle="Back" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <Container>
            <Progress.Container>
              <Progress.Step
                title="Build Settings"
                description="ビルド設定"
                state={
                  currentStep() === formStep.general
                    ? 'current'
                    : currentStep() > formStep.general
                    ? 'complete'
                    : 'incomplete'
                }
              />
              <Progress.Step
                title="Domains"
                description="アクセスURLの設定"
                state={
                  currentStep() === formStep.website
                    ? 'current'
                    : currentStep() > formStep.website
                    ? 'complete'
                    : 'incomplete'
                }
              />
            </Progress.Container>
            <Switch>
              <Match when={currentStep() === formStep.general}>
                <Show when={repo()}>
                  {(nonNullRepo) => (
                    <GeneralStep
                      repo={nonNullRepo()}
                      gotoNextStep={(appId) => {
                        setAppId(appId)
                        setCurrentStep(formStep.website)
                      }}
                    />
                  )}
                </Show>
              </Match>
              <Match when={currentStep() === formStep.website}>
                <Show when={app()}>{(nonNullApp) => <WebsiteStep app={nonNullApp()} />}</Show>
              </Match>
            </Switch>
          </Container>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
