import { createForm, Field, Form, isDirty } from '@formisch/solid/v2'
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
import AuthMethodField from './AuthMethodField'

type Props = {
  repo: Repository
  hasPermission: boolean
  refetch: () => void
}

const AuthConfigForm: Component<Props> = (props) => {
  // このパネルが編集するフィールドだけを持たせ、他のパネルの設定を検証・送信しない
  const toInput = (repo: Repository) => getInitialValueOfUpdateRepoForm(repo, ['url', 'auth'] as const)
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
        const request = toUpdateRepositoryRequest(output)
        // 認証設定に触れていなければ送らない (SSH の場合、送ると専用鍵が共通鍵に置き換わってしまう)
        if (!isDirty(form, { path: ['form', 'auth'] })) request.auth = undefined
        try {
          await client.updateRepository(request)
          sync.markSaved(submitted)
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
          <AuthMethodField form={form} repositoryID={props.repo.id} />
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
