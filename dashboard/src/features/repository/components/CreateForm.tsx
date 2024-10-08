import { styled } from '@macaron-css/solid'
import { Field, Form, type SubmitHandler, getValue, reset, setValue, setValues } from '@modular-forms/solid'
import { useNavigate } from '@solidjs/router'
import { type Component, createEffect, onMount } from 'solid-js'
import toast from 'solid-toast'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextField } from '/@/components/UI/TextField'
import { useRepositoryForm } from '/@/features/repository/provider/repositoryFormProvider'
import {
  type CreateOrUpdateRepositoryInput,
  getInitialValueOfCreateRepoForm,
  handleSubmitCreateRepositoryForm,
} from '/@/features/repository/schema/repositorySchema'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'

import { colorVars } from '/@/theme'
import AuthMethodField from './AuthMethodField'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-end',
    gap: '40px',
  },
})
const InputsContainer = styled('div', {
  base: {
    width: '100%',
    margin: '0 auto',
    padding: '20px 24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',

    borderRadius: '8px',
    background: colorVars.semantic.ui.primary,
  },
})

const CreateForm: Component = () => {
  const navigate = useNavigate()
  const { formStore } = useRepositoryForm()

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, getInitialValueOfCreateRepoForm())
    reset(formStore, {
      initialValues: getInitialValueOfCreateRepoForm(),
    })
  })

  // URLからリポジトリ名, 認証方法を自動入力
  createEffect(() => {
    const url = getValue(formStore, 'form.url')
    if (url === undefined || url === '') return

    // リポジトリ名を自動入力
    const repositoryName = extractRepositoryNameFromURL(url)
    setValue(formStore, 'form.name', repositoryName)

    // 認証方法を自動入力
    const isHTTPFormat = url.startsWith('http://') || url.startsWith('https://')
    if (!isHTTPFormat) {
      // Assume SSH or Git Protocol format
      setValues(formStore, {
        form: {
          auth: {
            method: 'ssh',
          },
        },
      })
    }
  })

  const handleSubmit: SubmitHandler<CreateOrUpdateRepositoryInput> = (values) =>
    handleSubmitCreateRepositoryForm(values, async (output) => {
      try {
        const res = await client.createRepository(output)
        toast.success('リポジトリを登録しました')
        // 新規アプリ作成ページに遷移
        navigate(`/apps/new?repositoryID=${res.id}`)
      } catch (e) {
        return handleAPIError(e, 'リポジトリの登録に失敗しました')
      }
    })

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Container>
        <InputsContainer>
          <Field of={formStore} name="form.url">
            {(field, fieldProps) => (
              <TextField
                label="Repository URL"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
              />
            )}
          </Field>
          <Field of={formStore} name="form.name">
            {(field, fieldProps) => (
              <TextField
                label="Repository Name"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
              />
            )}
          </Field>
          <AuthMethodField />
        </InputsContainer>
        <Button
          variants="primary"
          size="medium"
          rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
          type="submit"
          disabled={formStore.invalid || formStore.submitting}
          loading={formStore.submitting}
        >
          Register
        </Button>
      </Container>
    </Form>
  )
}

export default CreateForm
