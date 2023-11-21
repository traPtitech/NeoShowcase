import { Application, ApplicationConfig, DeployType, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { Progress } from '/@/components/UI/StepProgress'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { BuildConfigForm, BuildConfigs, configToForm, formToConfig } from '/@/components/templates/BuildConfigs'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { AppGeneralForm, GeneralConfig } from '/@/components/templates/GeneralConfig'
import { Nav } from '/@/components/templates/Nav'
import { WebsiteSetting, newWebsite } from '/@/components/templates/WebsiteSettings'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import {
  Field,
  Form,
  FormStore,
  SubmitHandler,
  createFormStore,
  getValues,
  setValue,
  validate,
} from '@modular-forms/solid'
import { useNavigate, useSearchParams } from '@solidjs/router'
import { For } from 'solid-js'
import { Component, Match, Show, Switch, createResource, createSignal } from 'solid-js'
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

type GeneralForm = AppGeneralForm & BuildConfigForm & { startOnCreate: boolean }

const GeneralStep: Component<{
  repo: Repository
  gotoNextStep: (appId: string) => void
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
      props.gotoNextStep(createdApp.id)
    } catch (e) {
      return handleAPIError(e, 'アプリケーションの登録に失敗しました')
    }
  }

  return (
    <Form of={form} onSubmit={handleSubmit} style={{ width: '100%' }}>
      <Container>
        <FormContainer>
          {/* 
            modular formsでは `FormStore<T extends AppGeneralForm, undefined>`のような
            genericsが使用できないためignoreしている
            */}
          {/* @ts-ignore */}
          <GeneralConfig repo={props.repo} formStore={form} editBranchId={false} hasPermission />
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
                  title="今すぐ起動する"
                  checked={field.value ?? false}
                  setChecked={(checked) => setValue(form, 'startOnCreate', checked)}
                  {...fieldProps}
                />
              </FormItem>
            )}
          </Field>
        </FormContainer>
        <Button
          size="medium"
          variants="primary"
          type="submit"
          rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
        >
          Next
        </Button>
      </Container>
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
const ButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    gap: '20px',
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
      <Container>
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
        <Nav title="Create Application" backTo={`/repos/${searchParams.repositoryID}`} backToTitle="Repository" />
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
