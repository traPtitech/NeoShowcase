import { createForm, Field, Form, reset } from '@formisch/solid/v2'
import { type Component, createEffect, untrack } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateRepositorySchema,
  getInitialValueOfUpdateRepoForm,
  toUpdateRepositoryRequest,
} from '../schema/repositorySchema'
import AuthMethodField from './AuthMethodField'

type Props = {
  repo: Repository
  hasPermission: boolean
  refetch: () => void
}

const AuthConfigForm: Component<Props> = (props) => {
  const form = createForm({
    schema: createOrUpdateRepositorySchema,
    initialInput: untrack(() => getInitialValueOfUpdateRepoForm(props.repo)),
  })

  // サーバー側のデータが更新されたらフォームを巻き戻す
  createEffect(
    () => props.repo,
    (repo) => reset(form, { initialInput: getInitialValueOfUpdateRepoForm(repo) }),
    { defer: true },
  )

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          await client.updateRepository(toUpdateRepositoryRequest(output))
          toast.success('リポジトリの設定を更新しました')
          props.refetch()
        } catch (e) {
          handleAPIError(e, 'リポジトリの設定の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <Field of={form} path={['form', 'url']}>
            {(field) => (
              <TextField
                label="Repository URL"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
              />
            )}
          </Field>
          <AuthMethodField form={form} />
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

export default AuthConfigForm
