import { Field, Form, type SubmitHandler, reset } from '@modular-forms/solid'
import { type Component, Show, createEffect, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationForm } from '../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationSchema,
  convertUpdateApplicationInput,
  updateApplicationFormInitialValues,
} from '../schema/applicationSchema'
import BranchField from './BranchField'

type Props = {
  app: Application
  repo: Repository
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const GeneralConfigForm: Component<Props> = (props) => {
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

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationSchema> = async (values) => {
    try {
      await client.updateRepository(convertUpdateApplicationInput(values))
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
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Field of={formStore} name="id">
        {() => null}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <Field of={formStore} name="name">
            {(field, fieldProps) => (
              <TextField
                label="Application Name"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <Field of={formStore} name="repositoryId">
            {(field, fieldProps) => (
              <TextField
                label="Repository ID"
                required
                info={{
                  props: {
                    content: 'リポジトリを移管する場合はIDを変更',
                  },
                }}
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <BranchField repo={props.repo} hasPermission={props.hasPermission} />
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

export default GeneralConfigForm
