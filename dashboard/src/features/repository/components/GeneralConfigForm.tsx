import { createForm, Field, Form } from '@formisch/solid/v2'
import { type Component, untrack } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { syncFormWithServer } from '/@/libs/serverForm'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateRepositorySchema,
  getInitialValueOfUpdateRepoForm,
  toUpdateRepositoryRequest,
} from '../schema/repositorySchema'

type Props = {
  repo: Repository
  hasPermission: boolean
  refetch: () => void
}

const GeneralConfigForm: Component<Props> = (props) => {
  // このパネルが編集するフィールドだけを持たせ、他のパネルの設定を検証・送信しない
  const toInput = (repo: Repository) => getInitialValueOfUpdateRepoForm(repo, ['name'] as const)
  const form = createForm({
    schema: createOrUpdateRepositorySchema,
    initialInput: untrack(() => toInput(props.repo)),
  })
  const sync = syncFormWithServer(form, { key: () => props.repo.id, input: () => toInput(props.repo), root: ['form'] })

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        const submitted = sync.snapshot()
        try {
          await client.updateRepository(toUpdateRepositoryRequest(output))
          sync.markSaved(submitted)
          toast.success('リポジトリ名を更新しました')
          props.refetch()
        } catch (e) {
          handleAPIError(e, 'リポジトリ名の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <Field of={form} path={['form', 'name']}>
            {(field) => (
              <TextField
                label="Repository Name"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
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
            disabled={!form.isDirty || form.isSubmitting || !props.hasPermission}
            loading={form.isSubmitting}
            tooltip={{
              content: !props.hasPermission ? '設定を変更するにはリポジトリのオーナーになる必要があります' : undefined,
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
