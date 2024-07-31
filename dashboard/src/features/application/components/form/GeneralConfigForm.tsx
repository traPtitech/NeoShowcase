import { Field, Form, type SubmitHandler, reset, setValues } from '@modular-forms/solid'
import { type Component, Show, createEffect, onMount, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationInput,
  handleSubmitUpdateApplicationForm,
  updateApplicationFormInitialValues,
} from '../../schema/applicationSchema'
import BranchField from './general/BranchField'

type Props = {
  app: Application
  repo: Repository
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const GeneralConfigForm: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  const discardChanges = () => {
    reset(
      untrack(() => formStore),
      {
        initialValues: updateApplicationFormInitialValues(props.app),
      },
    )
  }

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, updateApplicationFormInitialValues(props.app))
  })

  // reset forms when props.app changed
  createEffect(() => {
    discardChanges()
  })

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = (values) =>
    handleSubmitUpdateApplicationForm(values, async (output) => {
      try {
        await client.updateApplication(output)
        toast.success('アプリケーション設定を更新しました')
        props.refetchApp()
        // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
        setTimeout(props.refetchApp, 1000)
      } catch (e) {
        handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
      }
    })

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Field of={formStore} name="form.id">
        {() => null}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <Field of={formStore} name="form.name">
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
          <Field of={formStore} name="form.repositoryId">
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
