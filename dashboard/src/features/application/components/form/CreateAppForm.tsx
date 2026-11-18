import { createForm, Form, handleSubmit, setInput, validate } from '@formisch/solid/v2'
import { useNavigate, useSearchParams } from '@solidjs/router'
import { type Component, createEffect, createMemo, createSignal, For, Match, Show, Switch } from 'solid-js'
import { Progress } from '/@/components/UI/StepProgress'
import { client, getRepository, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateApplicationSchema,
  getInitialValueOfCreateAppForm,
  toCreateApplicationRequest,
} from '../../schema/applicationSchema'
import { envVarSchema, envVarsMessageToSchema } from '../../schema/envVarSchema'
import GeneralStep from './GeneralStep'
import RepositoryStep from './RepositoryStep'
import WebsiteStep from './WebsiteStep'

const steps = [
  { id: 'repository', title: '1. Repository', description: 'リポジトリの選択' },
  { id: 'general', title: '2. Build Settings', description: 'ビルド設定' },
  { id: 'website', title: '3. URLs', description: 'アクセスURLの設定' },
] as const

type StepId = (typeof steps)[number]['id']

const CreateAppForm: Component = () => {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()

  const form = createForm({
    schema: createOrUpdateApplicationSchema,
    initialInput: getInitialValueOfCreateAppForm(),
  })
  const envVarForm = createForm({ schema: envVarSchema, initialInput: envVarsMessageToSchema() })

  // 選択中のリポジトリは URL (?repositoryID=) が唯一の情報源。ステップもここから導出する
  const repositoryId = () => (typeof searchParams.repositoryID === 'string' ? searchParams.repositoryID : '')
  const repo = createMemo(() => (repositoryId() !== '' ? getRepository(repositoryId()) : undefined))
  const [urlStepReached, setUrlStepReached] = createSignal(false)
  const currentStep = (): StepId => (repositoryId() === '' ? 'repository' : urlStepReached() ? 'website' : 'general')

  createEffect(repositoryId, (id) => setInput(form, { path: ['form', 'repositoryId'], input: id }))

  const goToRepoStep = () => {
    setUrlStepReached(false)
    // 選択していたリポジトリをリセットする
    setSearchParams({ repositoryID: undefined })
  }

  // 全体を検証してエラーがなければ URL ステップへ進む
  // (環境変数は別のフォームなので、ここで一緒に検証しておく)
  const proceedToWebsiteStep = handleSubmit(form, async () => {
    const envVars = await validate(envVarForm, { shouldFocus: true })
    if (envVars.success) setUrlStepReached(true)
  })

  const handleCreate = async (request: ReturnType<typeof toCreateApplicationRequest>) => {
    const envVars = await validate(envVarForm)
    if (!envVars.success) {
      setUrlStepReached(false)
      return
    }
    try {
      const createdApp = await client.createApplication(request)

      try {
        await Promise.all(
          envVars.output.variables
            .filter((envVar) => !envVar.system && envVar.key !== '')
            .map((envVar) => client.setEnvVar({ applicationId: createdApp.id, ...envVar })),
        )
      } catch (e) {
        toast.error(
          <p>
            環境変数の設定に失敗しました (<a href={`/apps/${createdApp.id}/settings/envVars`}>再設定</a>)
            <br />
            {e instanceof Error ? e.message : '不明なエラー'}
          </p>,
        )
      }

      toast.success('アプリケーションを登録しました')
      navigate(`/apps/${createdApp.id}`)
    } catch (e) {
      handleAPIError(e, 'アプリケーションの登録に失敗しました')
    }
  }

  return (
    <div class={['flex w-full flex-col items-center gap-10', currentStep() === 'repository' && 'max-h-full']}>
      <Progress.Container>
        <For each={steps}>
          {(step, i) => (
            <Progress.Step
              title={step.title}
              description={step.description}
              state={
                currentStep() === step.id
                  ? 'current'
                  : steps.findIndex((s) => s.id === currentStep()) > i()
                    ? 'complete'
                    : 'incomplete'
              }
            />
          )}
        </For>
      </Progress.Container>
      <Form of={form} onSubmit={(output) => handleCreate(toCreateApplicationRequest(output))} class="w-full">
        <Switch>
          <Match when={currentStep() === 'repository'}>
            <RepositoryStep onSelect={(repo) => setSearchParams({ repositoryID: repo.id })} />
          </Match>
          <Match when={currentStep() === 'general'}>
            <Show when={repo()}>
              {(repo) => (
                <GeneralStep
                  form={form}
                  envVarForm={envVarForm}
                  repo={repo()}
                  backToRepoStep={goToRepoStep}
                  proceedToWebsiteStep={proceedToWebsiteStep}
                />
              )}
            </Show>
          </Match>
          <Match when={currentStep() === 'website'}>
            <WebsiteStep form={form} backToGeneralStep={() => setUrlStepReached(false)} />
          </Match>
        </Switch>
      </Form>
    </div>
  )
}

export default CreateAppForm
