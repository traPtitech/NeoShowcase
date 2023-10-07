import { ApplicationEnvVars } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextInput } from '/@/components/UI/TextInput'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, For, Show, createEffect, createResource } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const EnvVarsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    rowGap: '8px',
    columnGap: '24px',

    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})

const EnvVarConfig: Component<{
  appId: string
  envVars: PlainMessage<ApplicationEnvVars>
  refetchEnvVars: () => void
}> = (props) => {
  let formRef: HTMLFormElement
  const [envVars, setEnvVars] = createStore<
    {
      key: string
      value: string
      system: boolean
    }[]
  >(props.envVars.variables)

  const updateEnvVars = () => {
    setEnvVars((prev) =>
      prev
        // strip empty env var
        .filter((envVar) => envVar.key !== '' || envVar.value !== '')
        // add empty env var
        .concat({
          key: '',
          value: '',
          system: false,
        }),
    )
  }
  createEffect(updateEnvVars)

  const discardChanges = () => {
    setEnvVars(props.envVars.variables)
    updateEnvVars()
  }

  const existSameKey = (key: string) => {
    const sameKey = envVars.map((envVar) => envVar.key).filter((k) => k === key)
    return sameKey.length > 1
  }

  const saveChanges = async () => {
    if (!formRef.reportValidity()) return

    const oldVars = new Map(
      props.envVars.variables.filter((envVar) => !envVar.system).map((envVar) => [envVar.key, envVar.value]),
    )
    const newVars = new Map(
      envVars.filter((envVar) => !envVar.system && envVar.key !== '').map((envVar) => [envVar.key, envVar.value]),
    )

    const addedKeys = [...newVars.keys()].filter((key) => !oldVars.has(key))
    const deletedKeys = [...oldVars.keys()].filter((key) => !newVars.has(key))
    const updatedKeys = [...oldVars.keys()].filter((key) =>
      newVars.has(key) ? oldVars.get(key) !== newVars.get(key) : false,
    )

    const addEnvVarRequests = Array.from([...addedKeys, ...updatedKeys]).map((key) => {
      return client.setEnvVar({
        applicationId: props.appId,
        key,
        value: newVars.get(key),
      })
    })
    const deleteEnvVarRequests = Array.from(deletedKeys).map((key) => {
      return client.deleteEnvVar({
        applicationId: props.appId,
        key,
      })
    })
    try {
      await Promise.all([...addEnvVarRequests, ...deleteEnvVarRequests])
      toast.success('環境変数を更新しました')
      props.refetchEnvVars()
    } catch (e) {
      handleAPIError(e, '環境変数の更新に失敗しました')
    }
  }

  return (
    <FormBox.Container ref={formRef}>
      <FormBox.Forms>
        <EnvVarsContainer>
          <div>Key</div>
          <div>Value</div>
          <For each={envVars}>
            {(envVar, index) => (
              <>
                <TextInput
                  value={envVar.key}
                  readOnly={envVar.system}
                  onInput={(e) => {
                    setEnvVars(index(), 'key', e.currentTarget.value)
                    updateEnvVars()

                    if (existSameKey(envVar.key)) {
                      e.currentTarget.setCustomValidity('Keyはユニークである必要があります')
                    } else {
                      e.currentTarget.setCustomValidity('')
                    }
                  }}
                  required={envVar.value !== ''}
                  disabled={envVar.system}
                />
                <TextInput
                  value={envVar.value}
                  readOnly={envVar.system}
                  onInput={(e) => {
                    setEnvVars(index(), 'value', e.currentTarget.value)
                    updateEnvVars()
                  }}
                  disabled={envVar.system}
                />
              </>
            )}
          </For>
        </EnvVarsContainer>
      </FormBox.Forms>
      <FormBox.Actions>
        <Show when={true}>
          <Button color="borderError" size="small" type="button" onClick={discardChanges}>
            Discard Changes
          </Button>
        </Show>
        <Button color="primary" size="small" type="button" onClick={saveChanges} disabled={false}>
          Save
        </Button>
      </FormBox.Actions>
    </FormBox.Container>
  )
}

export default () => {
  const { app } = useApplicationData()
  const [envVars, { refetch: refetchEnvVars }] = createResource(
    () => app()?.id,
    (id) => client.getEnvVars({ id }),
  )

  const loaded = () => !!envVars()
  return (
    <DataTable.Container>
      <DataTable.Title>Environment Variables</DataTable.Title>
      <Show when={loaded()}>
        <EnvVarConfig appId={app()?.id} envVars={structuredClone(envVars())} refetchEnvVars={refetchEnvVars} />
      </Show>
    </DataTable.Container>
  )
}
