import { Field, Form, reset, type SubmitHandler, setValues } from '@modular-forms/solid'
import { type Component, createEffect, onMount, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { useRepositoryForm } from '/@/features/repository/provider/repositoryFormProvider'
import {
  type CreateOrUpdateRepositoryInput,
  getInitialValueOfUpdateRepoForm,
  handleSubmitUpdateRepositoryForm,
} from '/@/features/repository/schema/repositorySchema'
import { client, handleAPIError } from '/@/libs/api'

type Props = {
  repo: Repository
  refetchRepo: () => Promise<void>
  hasPermission: boolean
}

const GeneralConfigForm: Component<Props> = (props) => {
  const { formStore } = useRepositoryForm()

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, getInitialValueOfUpdateRepoForm(props.repo))
  })

  // reset forms when props.repo changed
  createEffect(() => {
    reset(
      untrack(() => formStore),
      {
        initialValues: getInitialValueOfUpdateRepoForm(props.repo),
      },
    )
  })

  const handleSubmit: SubmitHandler<CreateOrUpdateRepositoryInput> = (values) =>
    handleSubmitUpdateRepositoryForm(values, async (output) => {
      try {
        await client.updateRepository(output)
        toast.success('リポジトリ名を更新しました')
        await props.refetchRepo()
      } catch (e) {
        handleAPIError(e, 'リポジトリ名の更新に失敗しました')
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
                label="Repository Name"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
        </FormBox.Forms>
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
                  ? '設定を変更するにはリポジトリのオーナーになる必要があります'
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
