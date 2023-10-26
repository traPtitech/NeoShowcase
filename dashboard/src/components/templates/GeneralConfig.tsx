import { CreateApplicationRequest, Repository, UpdateApplicationRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { useBranchesSuggestion } from '/@/libs/branchesSuggestion'
import { PlainMessage } from '@bufbuild/protobuf'
import { Component, Show } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { TextInput } from '../UI/TextInput'
import { FormItem } from './FormItem'
import { ComboBox } from './Select'

interface GeneralConfigProps {
  repo: Repository
  config: PlainMessage<CreateApplicationRequest> | PlainMessage<UpdateApplicationRequest>
  setConfig: SetStoreFunction<GeneralConfigProps['config']>
  editBranchId?: boolean
}

export const GeneralConfig: Component<GeneralConfigProps> = (props) => {
  const branchesSuggestion = useBranchesSuggestion(
    () => props.repo.id,
    () => props.config.refName ?? '',
  )

  return (
    <>
      <FormItem title="Application Name" required>
        <TextInput
          required
          value={props.config.name}
          onInput={(e) => {
            props.setConfig('name', e.target.value)
          }}
        />
      </FormItem>
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
          <TextInput
            required
            value={props.config.repositoryId}
            onInput={(e) => {
              props.setConfig('repositoryId', e.target.value)
            }}
          />
        </FormItem>
      </Show>
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
          required
          value={props.config.refName}
          onInput={(e) => props.setConfig('refName', e.target.value)}
          items={branchesSuggestion().map((branch) => ({
            title: branch,
            value: branch,
          }))}
          setSelected={(branch) => {
            props.setConfig('refName', branch)
          }}
        />
      </FormItem>
    </>
  )
}
