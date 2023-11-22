import { CreateApplicationRequest, Repository, UpdateApplicationRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { useBranches } from '/@/libs/branchesSuggestion'
import { PlainMessage } from '@bufbuild/protobuf'
import { Field, FormStore, required, setValue } from '@modular-forms/solid'
import { Component, Show } from 'solid-js'
import { TextField } from '../UI/TextField'
import { ComboBox } from './Select'

export type AppGeneralForm = Pick<
  PlainMessage<CreateApplicationRequest> | PlainMessage<UpdateApplicationRequest>,
  'name' | 'repositoryId' | 'refName'
>

interface GeneralConfigProps {
  repo: Repository
  formStore: FormStore<AppGeneralForm, undefined>
  editBranchId?: boolean
  hasPermission: boolean
}

export const GeneralConfig: Component<GeneralConfigProps> = (props) => {
  const branches = useBranches(() => props.repo.id)

  return (
    <>
      <Field of={props.formStore} name="name" validate={required('Please Enter Application Name')}>
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
      <Field of={props.formStore} name="repositoryId" validate={required('Please Enter Repository ID')}>
        {(field, fieldProps) => (
          <Show when={props.editBranchId}>
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
          </Show>
        )}
      </Field>
      <Field of={props.formStore} name="refName" validate={required('Please Enter Branch Name')}>
        {(field, fieldProps) => (
          <ComboBox
            label="Branch"
            required
            info={{
              props: {
                content: (
                  <>
                    <div>Gitブランチ名またはRef</div>
                    <div>入力欄をクリックして候補を表示</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            options={branches().map((branch) => ({
              label: branch,
              value: branch,
            }))}
            value={field.value}
            error={field.error}
            setValue={(v) => {
              setValue(props.formStore, 'refName', v)
            }}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
    </>
  )
}
