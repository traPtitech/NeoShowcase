import { createForm, Field, Form, setInput } from '@formisch/solid/v2'
import { useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateRepositorySchema,
  getInitialValueOfCreateRepoForm,
  toCreateRepositoryRequest,
} from '../schema/repositorySchema'
import AuthMethodField from './AuthMethodField'

const CreateForm: Component = () => {
  const navigate = useNavigate()
  const form = createForm({
    schema: createOrUpdateRepositorySchema,
    initialInput: getInitialValueOfCreateRepoForm(),
  })

  // URLからリポジトリ名, 認証方法を自動入力
  const autofillFromURL = (url: string) => {
    if (url === '') return

    setInput(form, { path: ['form', 'name'], input: extractRepositoryNameFromURL(url) })

    const isHTTPFormat = url.startsWith('http://') || url.startsWith('https://')
    if (!isHTTPFormat) {
      // Assume SSH or Git Protocol format
      setInput(form, { path: ['form', 'auth', 'method'], input: 'ssh' })
    }
  }

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          const res = await client.createRepository(toCreateRepositoryRequest(output))
          toast.success('リポジトリを登録しました')
          // 新規アプリ作成ページに遷移
          navigate(`/apps/new?repositoryID=${res.id}`)
        } catch (e) {
          handleAPIError(e, 'リポジトリの登録に失敗しました')
        }
      }}
    >
      <div class="flex w-full flex-col items-end gap-10">
        <div class="bg-ui-primary mx-auto flex w-full flex-col gap-6 rounded-lg px-6 py-5">
          <Field of={form} path={['form', 'url']}>
            {(field) => (
              <TextField
                label="Repository URL"
                required
                {...field.props}
                onInput={(e) => {
                  field.props.onInput(e)
                  autofillFromURL(e.currentTarget.value)
                }}
                value={field.input ?? ''}
                error={field.errors?.[0]}
              />
            )}
          </Field>
          <Field of={form} path={['form', 'name']}>
            {(field) => (
              <TextField
                label="Repository Name"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
              />
            )}
          </Field>
          <AuthMethodField form={form} />
        </div>
        <Button
          variants="primary"
          size="medium"
          rightIcon={<div class="i-material-symbols:arrow-forward shrink-0 text-2xl/6" />}
          type="submit"
          disabled={form.isSubmitting}
          loading={form.isSubmitting}
        >
          Register
        </Button>
      </div>
    </Form>
  )
}

export default CreateForm
