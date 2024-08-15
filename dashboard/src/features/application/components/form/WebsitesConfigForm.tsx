import { styled } from '@macaron-css/solid'
import { createFormStore, getValues, setValue, valiForm } from '@modular-forms/solid'
import { type Component, For, createResource, createSignal } from 'solid-js'
import toast from 'solid-toast'
import { parse } from 'valibot'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import {
  type CreateWebsiteInput,
  createWebsiteInitialValues,
  createWebsiteSchema,
  websiteMessageToSchema,
} from '../../schema/websiteSchema'
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
  const [formStores, { mutate }] = createResource(
    () => props.app.websites,
    (websites) =>
      websites.map((w) =>
        createFormStore<CreateWebsiteInput, undefined>({
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
      const newForm = createFormStore<CreateWebsiteInput, undefined>({
        initialValues: createWebsiteInitialValues(defaultDomain),
        // validateOn: 'blur',
        validate: async (input) => {
          console.log(input)
          console.log(await valiForm(createWebsiteSchema)(input))
          return valiForm(createWebsiteSchema)(input)
        },
      })

      mutate((websites) => websites?.concat([newForm]))
    }
  }

  const isRuntimeApp = () => {
    const configCase = props.app.config?.buildConfig.case
    return configCase === 'runtimeBuildpack' || configCase === 'runtimeDockerfile' || configCase === 'runtimeCmd'
  }

  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const applyChanges = async () => {
    setIsSubmitting(true)

    try {
      /**
       * 送信するWebsite設定
       * - 変更を保存しないものの、initial value
       * - 変更して保存するもの ( = `readyToSave`)
       * - 追加するもの ( = `added`)
       * - 削除しないもの ( = not `readyToDelete`)
       */
      const websitesToSave =
        formStores()
          ?.map((form) => {
            const values = getValues(form)
            switch (values.state) {
              case 'noChange':
                return form.internal.initialValues
              case 'readyToChange':
              case 'added':
                return values
              case 'readyToDelete':
                return undefined
            }
          })
          .filter((w): w is Exclude<typeof w, undefined> => w !== undefined) ?? []

      const parsedWebsites = websitesToSave.map((w) => parse(createWebsiteSchema, w))

      await client.updateApplication({
        id: props.app.id,
        websites: {
          websites: parsedWebsites,
        },
      })
      toast.success('ウェブサイト設定を保存しました')
      void props.refetchApp()
    } catch (e) {
      // `readyToChange` を `noChange` に戻す
      for (const form of formStores() ?? []) {
        const values = getValues(form)
        if (values.state === 'readyToChange') {
          setValue(form, 'state', 'noChange')
        }
      }
      handleAPIError(e, 'Failed to save website settings')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <>
      <For each={formStores()}>{(formStore) => <pre>{JSON.stringify(getValues(formStore), null, 2)}</pre>}</For>
      <For each={formStores()}>
        {(formStore) => (
          <WebsiteFieldGroup
            formStore={formStore}
            isRuntimeApp={isRuntimeApp()}
            applyChanges={applyChanges}
            isSubmitting={isSubmitting()}
            readonly={!props.hasPermission}
          />
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
