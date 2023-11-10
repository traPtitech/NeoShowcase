import { CreateApplicationRequest, Repository, UpdateApplicationRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { useBranchesSuggestion } from '/@/libs/branchesSuggestion'
import { PlainMessage } from '@bufbuild/protobuf'
import { Field, FormStore, getValue, required, setValue } from '@modular-forms/solid'
import { Component, Show } from 'solid-js'
import { TextInput } from '../UI/TextInput'
import { FormItem } from './FormItem'
import { ComboBox } from './Select'

export type AppGeneralForm = Pick<
  PlainMessage<CreateApplicationRequest> | PlainMessage<UpdateApplicationRequest>,
  'name' | 'repositoryId' | 'refName'
>

interface GeneralConfigProps {
  repo: Repository
  formStore: FormStore<AppGeneralForm, undefined>
  editBranchId?: boolean
}

export const GeneralConfig: Component<GeneralConfigProps> = (props) => {
  const branchesSuggestion = useBranchesSuggestion(
    () => props.repo.id,
    () => getValue(props.formStore, 'refName') ?? '',
  )

  return (
    <>
      <Field of={props.formStore} name="name" validate={required('Please Enter Application Name')}>
        {(field, fieldProps) => (
          <FormItem title="Application Name" required>
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
      <Field of={props.formStore} name="repositoryId" validate={required('Please Enter Repository ID')}>
        {(field, fieldProps) => (
          <Show when={props.editBranchId}>
            <FormItem
              title="Repository ID"
              required
              tooltip={{
                props: {
                  content: 'リポジトリを移管する場合はIDを変更',
                },
              }}
            >
              <TextInput required value={field.value} error={field.error} {...fieldProps} />
            </FormItem>
          </Show>
        )}
      </Field>
      <Field of={props.formStore} name="refName" validate={required('Please Enter Branch Name')}>
        {(field, fieldProps) => (
          <FormItem
            title="Branch"
            required
            tooltip={{
              props: {
                content: (
                  <>
                    <div>Gitブランチ名またはRef</div>
                    <div>入力欄をクリックして候補を表示</div>
                  </>
                ),
              },
            }}
          >
            <ComboBox
              value={field.value}
              items={branchesSuggestion().map((branch) => ({
                title: branch,
                value: branch,
              }))}
              setSelected={(v) => {
                setValue(props.formStore, 'refName', v)
              }}
              error={field.error}
              {...fieldProps}
            />
          </FormItem>
        )}
      </Field>
    </>
  )
}
