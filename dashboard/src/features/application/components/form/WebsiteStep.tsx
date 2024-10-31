import { FieldArray, getValue, getValues, insert } from '@modular-forms/solid'
import { type Component, For, Show } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { styled } from '/@/components/styled-components'
import { List } from '/@/components/templates/List'
import { systemInfo } from '/@/libs/api'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import { getInitialValueOfCreateWebsiteForm } from '../../schema/websiteSchema'
import WebsiteFieldGroup from './website/WebsiteFieldGroup'

const FieldRow = styled('div', 'w-full px-6 py-5 [&:not(:first-child)]:border-ui-border [&:not(:first-child)]:border-t')

const WebsiteStep: Component<{
  backToGeneralStep: () => void
}> = (props) => {
  const { formStore } = useApplicationForm()

  const defaultDomain = () => systemInfo()?.domains.at(0)

  const addFormStore = () => {
    const _defaultDomain = defaultDomain()
    if (!_defaultDomain) {
      throw new Error('Default domain is not found')
    }
    insert(formStore, 'form.websites', {
      value: getInitialValueOfCreateWebsiteForm(_defaultDomain),
    })
  }

  const isRuntimeApp = () => getValue(formStore, 'form.config.deployConfig.type') === 'runtime'

  const showAddMoreButton = () => {
    const websites = getValues(formStore, 'form.websites')
    return websites && websites.length > 0
  }

  return (
    <Show when={systemInfo()}>
      <div class="flex w-full flex-col items-center gap-10">
        <div class="flex w-full flex-col items-center gap-6">
          <div class="flex w-full flex-col gap-0.25 overflow-hidden rounded-lg border border-ui-border bg-ui-primary">
            <FieldArray of={formStore} name="form.websites">
              {(fieldArray) => (
                <For
                  each={fieldArray.items}
                  fallback={
                    <List.PlaceHolder>
                      <span class="i-material-symbols:link-off text-20/20" />
                      URLが設定されていません
                      <Button
                        variants="primary"
                        size="medium"
                        rightIcon={<span class="i-material-symbols:add text-2xl/6" />}
                        onClick={addFormStore}
                        type="button"
                      >
                        Add URL
                      </Button>
                    </List.PlaceHolder>
                  }
                >
                  {(_, index) => (
                    <FieldRow>
                      <WebsiteFieldGroup index={index()} isRuntimeApp={isRuntimeApp()} />
                    </FieldRow>
                  )}
                </For>
              )}
            </FieldArray>
            <Show when={showAddMoreButton()}>
              <FieldRow>
                <div class="flex justify-center">
                  <Button
                    onclick={() => {
                      addFormStore()
                    }}
                    variants="border"
                    size="small"
                    leftIcon={<span class="i-material-symbols:add text-xl/5" />}
                    type="button"
                  >
                    Add More
                  </Button>
                </div>
              </FieldRow>
            </Show>
          </div>
        </div>
        <div class="flex gap-5">
          <Button
            size="medium"
            variants="ghost"
            leftIcon={<span class="i-material-symbols:arrow-back text-2xl/6" />}
            onClick={props.backToGeneralStep}
          >
            Back
          </Button>
          <Button
            type="submit"
            size="medium"
            variants="primary"
            disabled={formStore.invalid || formStore.submitting}
            loading={formStore.submitting}
          >
            Create Application
          </Button>
        </div>
      </div>
    </Show>
  )
}

export default WebsiteStep
