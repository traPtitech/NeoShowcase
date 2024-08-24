import { styled } from '@macaron-css/solid'
import { FieldArray, getValue, getValues, insert } from '@modular-forms/solid'
import { type Component, For, Show } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { List } from '/@/components/templates/List'
import { systemInfo } from '/@/libs/api'
import { colorVars } from '/@/theme'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import { getInitialValueOfCreateWebsiteForm } from '../../schema/websiteSchema'
import WebsiteFieldGroup from './website/WebsiteFieldGroup'

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

const Container = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    background: colorVars.semantic.ui.primary,
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    display: 'flex',
    flexDirection: 'column',
    gap: '1px',
  },
})
const FieldRow = styled('div', {
  base: {
    width: '100%',
    padding: '20px 24px',

    selectors: {
      '&:not(:first-child)': {
        borderTop: `1px solid ${colorVars.semantic.ui.border}`,
      },
    },
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
      <FormsContainer>
        <DomainsContainer>
          <Container>
            <FieldArray of={formStore} name="form.websites">
              {(fieldArray) => (
                <For
                  each={fieldArray.items}
                  fallback={
                    <List.PlaceHolder>
                      <MaterialSymbols displaySize={80}>link_off</MaterialSymbols>
                      URLが設定されていません
                      <Button
                        variants="primary"
                        size="medium"
                        rightIcon={<MaterialSymbols>add</MaterialSymbols>}
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
                <AddMoreButtonContainer>
                  <Button
                    onclick={() => {
                      addFormStore()
                    }}
                    variants="border"
                    size="small"
                    leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
                    type="button"
                  >
                    Add More
                  </Button>
                </AddMoreButtonContainer>
              </FieldRow>
            </Show>
          </Container>
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
            type="submit"
            size="medium"
            variants="primary"
            disabled={formStore.invalid || formStore.submitting}
            loading={formStore.submitting}
          >
            Create Application
          </Button>
        </ButtonsContainer>
      </FormsContainer>
    </Show>
  )
}

export default WebsiteStep
