import { styled } from '@macaron-css/solid'
import { createFormStore, getValue, getValues, valiForm } from '@modular-forms/solid'
import { type Component, For } from 'solid-js'
import { createMutable, createStore } from 'solid-js/store'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { systemInfo } from '/@/libs/api'
import { createWebsiteInitialValues, createWebsiteSchema, websiteMessageToSchema } from '../../schema/websiteSchema'
import WebsiteFieldGroup from './website/WebsiteFieldGroup'

const AddMoreButtonContainer = styled('div', {
  base: {
    display: 'flex',
    justifyContent: 'center',
  },
})

type Props = {
  app: Application
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const WebsiteConfigForm: Component<Props> = (props) => {
  const formStores = createMutable(
    props.app.websites.map((w) =>
      createFormStore({
        initialValues: websiteMessageToSchema(w),
        validate: async (input) => {
          console.log(input)
          console.log(await valiForm(createWebsiteSchema)(input))
          return valiForm(createWebsiteSchema)(input)
        },
      }),
    ),
  )

  const availableDomains = systemInfo()?.domains ?? []

  const addFormStore = () => {
    const defaultDomain = availableDomains.at(0)
    if (defaultDomain) {
      formStores.push(
        createFormStore({
          initialValues: createWebsiteInitialValues(defaultDomain),
          validateOn: 'blur',
          validate: async (input) => {
            console.log(input)
            console.log(await valiForm(createWebsiteSchema)(input))
            return valiForm(createWebsiteSchema)(input)
          },
        }),
      )
    }
  }

  const isRuntimeApp = () => {
    const configCase = props.app.config?.buildConfig.case
    return configCase === 'runtimeBuildpack' || configCase === 'runtimeDockerfile' || configCase === 'runtimeCmd'
  }

  const applyChanges = () => {
    console.log(formStores.map((f) => getValues(f)))
  }

  return (
    <>
      <For each={formStores}>
        {(formStore) => (
          <WebsiteFieldGroup formStore={formStore} isRuntimeApp={isRuntimeApp()} applyChanges={applyChanges} />
        )}
      </For>
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
    </>
  )
}

export default WebsiteConfigForm
