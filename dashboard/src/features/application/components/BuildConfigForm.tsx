import { Field, Form, type SubmitHandler, getValues, reset } from '@modular-forms/solid'
import { type Component, Show, createEffect, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationForm } from '../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationInput,
  convertUpdateApplicationInput,
  updateApplicationFormInitialValues,
} from '../schema/applicationSchema'
import BuildTypeField from './BuildTypeField'

type Props = {
  app: Application
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const BuildConfigForm: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  // reset forms when props.app changed
  createEffect(() => {
    reset(
      untrack(() => formStore),
      {
        initialValues: updateApplicationFormInitialValues(props.app),
      },
    )
  })

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = async (values) => {
    try {
      await client.updateApplication(convertUpdateApplicationInput(values))
      toast.success('アプリケーション設定を更新しました')
      props.refetchApp()
      // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
      setTimeout(props.refetchApp, 1000)
    } catch (e) {
      handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
    }
  }

  const discardChanges = () => {
    reset(formStore)
  }

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      {JSON.stringify(getValues(formStore))}
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Field of={formStore} name="id">
        {() => null}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <BuildTypeField formStore={formStore} readonly={!props.hasPermission} />
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={formStore.dirty && !formStore.submitting}>
            <Button variants="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
          </Show>
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
      </FormBox.Container>
    </Form>
  )
}

export default BuildConfigForm
