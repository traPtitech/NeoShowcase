import { Repository, UpdateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextInput } from '/@/components/UI/TextInput'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { client, handleAPIError } from '/@/libs/api'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
    overflowY: 'auto',

    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
  },
})

const GeneralConfig: Component<{
  repo: Repository
  refetchRepo: () => void
}> = (props) => {
  let formRef: HTMLFormElement

  const [updateReq, setUpdateReq] = createStore<PlainMessage<UpdateRepositoryRequest>>({
    id: props.repo.id,
    name: props.repo.name,
  })
  const discardChanges = () => {
    setUpdateReq({
      name: props.repo.name,
    })
  }
  const nameChanged = () => props.repo.name !== updateReq.name
  const saveChanges = async () => {
    try {
      // validate form
      if (!formRef.reportValidity()) {
        return
      }
      await client.updateRepository(updateReq)
      toast.success('Project名を更新しました')
      props.refetchRepo()
    } catch (e) {
      handleAPIError(e, 'Project名の更新に失敗しました')
    }
  }

  return (
    <FormBox.Container ref={formRef}>
      <FormBox.Forms>
        <FormItem title="Project Name" required>
          <TextInput
            required
            value={updateReq.name}
            onInput={(e) => {
              setUpdateReq('name', (e.target as HTMLInputElement).value)
            }}
          />
        </FormItem>
      </FormBox.Forms>
      <FormBox.Actions>
        <Show when={nameChanged()}>
          <Button color="borderError" size="small" onClick={discardChanges} type="button">
            Discard Changes
          </Button>
        </Show>
        <Button color="primary" size="small" onClick={saveChanges} type="button" disabled={!nameChanged()}>
          Save
        </Button>
      </FormBox.Actions>
    </FormBox.Container>
  )
}

export default () => {
  const { repo, refetchRepo } = useRepositoryData()
  const loaded = () => !!repo()

  return (
    <Container>
      General
      <Show when={loaded()}>
        <GeneralConfig repo={repo()} refetchRepo={refetchRepo} />
      </Show>
    </Container>
  )
}
