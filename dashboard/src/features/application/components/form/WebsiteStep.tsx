import { styled } from '@macaron-css/solid'
import { type FormStore, createFormStore, validate } from '@modular-forms/solid'
import { type Accessor, type Component, For, type Setter, Show, createSignal } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { List } from '/@/components/templates/List'
import { type WebsiteFormStatus, WebsiteSetting, newWebsite } from '/@/components/templates/app/WebsiteSettings'
import { systemInfo } from '/@/libs/api'

const FormsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '40px',
  },
})
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
const ButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    gap: '20px',
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
      console.error(err)
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

export default WebsiteStep
