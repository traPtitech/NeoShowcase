import type { PartialMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, FieldArray, Form, type SubmitHandler, getValues, insert, reset, setValues } from '@modular-forms/solid'
import { type Component, For, Show, createEffect, onMount, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { Application, UpdateApplicationRequest_UpdateWebsites } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import FormBox from '/@/components/layouts/FormBox'
import { List } from '/@/components/templates/List'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { colorVars } from '/@/theme'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationInput,
  handleSubmitUpdateApplicationForm,
  updateApplicationFormInitialValues,
} from '../../schema/applicationSchema'
import { createWebsiteInitialValues } from '../../schema/websiteSchema'
import WebsiteFieldGroup from './website/WebsiteFieldGroup'

const Container = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
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
    width: '100%',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
})

type Props = {
  app: Application
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const WebsiteConfigForm: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, updateApplicationFormInitialValues(props.app))
  })

  // reset forms when props.app changed
  createEffect(() => {
    reset(
      untrack(() => formStore),
      {
        initialValues: updateApplicationFormInitialValues(props.app),
      },
    )
  })

  const defaultDomain = () => systemInfo()?.domains.at(0)

  const addFormStore = () => {
    const _defaultDomain = defaultDomain()
    if (!_defaultDomain) {
      throw new Error('Default domain is not found')
    }
    insert(formStore, 'form.websites', {
      value: createWebsiteInitialValues(_defaultDomain),
    })
  }

  const isRuntimeApp = () => {
    const configCase = props.app.config?.buildConfig.case
    return configCase === 'runtimeBuildpack' || configCase === 'runtimeDockerfile' || configCase === 'runtimeCmd'
  }

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = (values) =>
    handleSubmitUpdateApplicationForm(values, async (output) => {
      try {
        console.log(output)
        // websiteがすべて削除されている場合、modularformsでは空配列ではなくundefinedになってしまう
        // undefinedを渡した場合、APIとしては 無更新 として扱われるため、空配列を渡す
        if (output.websites === undefined) {
          console.log('output.websites is undefined')
          output.websites = [] as PartialMessage<UpdateApplicationRequest_UpdateWebsites>
          console.log(output.websites)
        }

        await client.updateApplication(output)
        toast.success('ウェブサイト設定を更新しました')
        props.refetchApp()
        // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
        setTimeout(props.refetchApp, 1000)
      } catch (e) {
        handleAPIError(e, 'ウェブサイト設定の更新に失敗しました')
      }
    })

  const showAddMoreButton = () => {
    const websites = getValues(formStore, 'form.websites')
    return websites && websites.length > 0
  }

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Field of={formStore} name="form.id">
        {() => null}
      </Field>
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
                  <WebsiteFieldGroup index={index()} isRuntimeApp={isRuntimeApp()} readonly={!props.hasPermission} />
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
        <FormBox.Actions>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={formStore.invalid || !formStore.dirty || formStore.submitting || !props.hasPermission}
            loading={formStore.submitting}
            tooltip={{
              props: {
                content: !props.hasPermission
                  ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                  : undefined,
              },
            }}
          >
            Save
          </Button>
        </FormBox.Actions>
      </Container>
    </Form>
  )
}

export default WebsiteConfigForm
