import { Field, Form, getValue, type SubmitHandler, setValue, validate } from '@modular-forms/solid'
import { useNavigate, useSearchParams } from '@solidjs/router'
import {
  type Component,
  createEffect,
  createResource,
  createSignal,
  For,
  Match,
  onMount,
  Show,
  Switch,
  untrack,
} from 'solid-js'
import toast from 'solid-toast'
import { Progress } from '/@/components/UI/StepProgress'
import { client, handleAPIError } from '/@/libs/api'
import { clsx } from '/@/libs/clsx'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import { type CreateOrUpdateApplicationInput, handleSubmitCreateApplicationForm } from '../../schema/applicationSchema'
import GeneralStep from './GeneralStep'
import RepositoryStep from './RepositoryStep'
import WebsiteStep from './WebsiteStep'

enum formStep {
  repository = 0,
  general = 1,
  website = 2,
}

const CreateAppForm: Component = () => {
  const { formStore } = useApplicationForm()

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValue(formStore, 'type', 'create')
  })

  const [searchParams, setParam] = useSearchParams()
  const [currentStep, setCurrentStep] = createSignal(formStep.repository)

  const goToRepoStep = () => {
    setCurrentStep(formStep.repository)
    // 選択していたリポジトリをリセットする
    setParam({ repositoryID: undefined })
  }
  const goToGeneralStep = () => {
    setCurrentStep(formStep.general)
  }
  const goToWebsiteStep = () => {
    setCurrentStep(formStep.website)
  }

  const [repoBySearchParam] = createResource(
    () => searchParams.repositoryID ?? '',
    (id) => {
      return typeof id === 'string' && id !== '' ? client.getRepository({ repositoryId: id }) : undefined
    },
  )

  // repositoryIDがない場合はリポジトリ選択ステップに遷移
  createEffect(() => {
    if (!searchParams.repositoryID) {
      goToRepoStep()
    }
  })

  // repoBySearchParam更新時にformのrepositoryIdを更新する
  createEffect(() => {
    setValue(
      untrack(() => formStore),
      'form.repositoryId',
      repoBySearchParam()?.id,
    )
  })

  // リポジトリ選択ステップで中に、リポジトリが選択された場合は次のステップに遷移
  createEffect(() => {
    if (currentStep() === formStep.repository && getValue(formStore, 'form.repositoryId')) {
      goToGeneralStep()
    }
  })

  const handleGeneralToWebsiteStep = async () => {
    const isValid = await validate(formStore)
    // modularformsではsubmitフラグが立っていないとrevalidateされないため、手動でsubmitフラグを立てる
    // TODO: internalのAPIを使っているため、将来的には変更が必要
    formStore.internal.submitted.set(true)
    if (isValid) {
      goToWebsiteStep()
    }
  }

  const navigate = useNavigate()
  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = (values) =>
    handleSubmitCreateApplicationForm(values, async (output) => {
      try {
        const createdApp = await client.createApplication(output)
        toast.success('アプリケーションを登録しました')
        navigate(`/apps/${createdApp.id}`)
      } catch (e) {
        handleAPIError(e, 'アプリケーションの登録に失敗しました')
      }
    })

  return (
    <div
      class={clsx('flex w-full flex-col items-center gap-10', currentStep() === formStep.repository && 'max-h-full')}
    >
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
              state={currentStep() === step.step ? 'current' : currentStep() > step.step ? 'complete' : 'incomplete'}
            />
          )}
        </For>
      </Progress.Container>
      <Form of={formStore} onSubmit={handleSubmit} shouldActive={false} class="w-full">
        <Field of={formStore} name="type">
          {() => null}
        </Field>
        <Field of={formStore} name="form.repositoryId">
          {() => null}
        </Field>
        <Switch>
          <Match when={currentStep() === formStep.repository}>
            <RepositoryStep
              onSelect={(repo) => {
                setParam({
                  repositoryID: repo.id,
                })
                goToGeneralStep()
              }}
            />
          </Match>
          <Match when={currentStep() === formStep.website}>
            <WebsiteStep backToGeneralStep={goToGeneralStep} />
          </Match>
        </Switch>
        <div class={clsx(currentStep() !== formStep.general && 'hidden')}>
          <Show when={repoBySearchParam()}>
            {(nonNullRepo) => (
              <GeneralStep
                repo={nonNullRepo()}
                backToRepoStep={goToRepoStep}
                proceedToWebsiteStep={handleGeneralToWebsiteStep}
              />
            )}
          </Show>
        </div>
      </Form>
    </div>
  )
}

export default CreateAppForm
