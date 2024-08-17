import { Form, type SubmitHandler, getValue, setValue, setValues } from '@modular-forms/solid'
import { useSearchParams } from '@solidjs/router'
import {
  type Component,
  Match,
  Show,
  Switch,
  createEffect,
  createResource,
  createSignal,
  onMount,
  untrack,
} from 'solid-js'
import toast from 'solid-toast'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationInput,
  createApplicationFormInitialValues,
  handleSubmitCreateApplicationForm,
} from '../../schema/applicationSchema'
import GeneralStep from './GeneralStep'
import RepositoryStep from './RepositoryStep'

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
    setValues(formStore, createApplicationFormInitialValues())
  })

  const [searchParams, setParam] = useSearchParams()
  const [currentStep, setCurrentStep] = createSignal(formStep.repository)

  const backToRepoStep = () => {
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

  const [repoBySearchParam, { mutate: mutateRepo }] = createResource(
    () => searchParams.repositoryID,
    (id) => client.getRepository({ repositoryId: id }),
  )
  // repoBySearchParam更新時にformのrepositoryIdを更新する
  createEffect(() => {
    setValue(
      untrack(() => formStore),
      'form.repositoryId',
      repoBySearchParam()?.id,
    )
  })

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = (values) =>
    handleSubmitCreateApplicationForm(values, async (output) => {
      try {
        await client.createApplication(output)
        toast.success('アプリケーション設定を更新しました')
      } catch (e) {
        handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
      }
    })

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      <Switch>
        <Match when={currentStep() === formStep.repository}>
          <RepositoryStep setRepo={(repo) => mutateRepo(repo)} />
        </Match>
        <Match when={currentStep() === formStep.general}>
          <Show when={repoBySearchParam()}>
            {(nonNullRepo) => (
              <GeneralStep
                repo={nonNullRepo()}
                backToRepoStep={backToRepoStep}
                proceedToWebsiteStep={goToWebsiteStep}
              />
            )}
          </Show>
        </Match>
        <Match when={currentStep() === formStep.website}>
          <WebsiteStep
            isRuntimeApp={isRuntimeApp()}
            backToGeneralStep={goToGeneralStep}
            websiteForms={websiteForms}
            setWebsiteForms={setWebsiteForms}
            submit={submit}
          />
        </Match>
      </Switch>
    </Form>
  )
}

export default CreateAppForm
