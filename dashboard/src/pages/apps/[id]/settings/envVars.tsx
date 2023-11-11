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
import { SubmitHandler, createForm, custom, getValue, getValues, insert, remove, reset } from '@modular-forms/solid'
import { Component, For, Show, createEffect, createReaction, createResource, on } from 'solid-js'
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
  const [envVarForm, EnvVar] = createForm<PlainMessage<ApplicationEnvVars>>({
    initialValues: {
      variables: props.envVars.variables,
    },
  })

  const discardChanges = () => {
    reset(envVarForm, {
      initialValues: {
        variables: props.envVars.variables,
      },
    })
    stripEnvVars()
  }

  // reset form when envVars updated
  createEffect(
    on(
      () => props.envVars,
      () => {
        discardChanges()
      },
    ),
  )

  const stripEnvVars = () => {
    const forms = getValues(envVarForm, 'variables') as PlainMessage<ApplicationEnvVars>['variables']
    // remove all empty env vars
    forms
      .map((envVar, index) => (envVar.key === '' && envVar.value === '' ? index : null))
      .filter((index): index is number => index !== null)
      .reverse()
      .forEach((index) => {
        remove(envVarForm, 'variables', { at: index })
      })
    // add empty env var
    insert(envVarForm, 'variables', {
      value: { key: '', value: '', system: false },
    })
    // 次にvariablesが変更された時に1度だけ再度stripする
    track(() => getValues(envVarForm, 'variables'))
  }
  const track = createReaction(() => {
    stripEnvVars()
  })

  const isUniqueKey = (key?: string) => {
    const sameKey = (getValues(envVarForm, 'variables') as PlainMessage<ApplicationEnvVars>['variables'])
      .map((envVar) => envVar.key)
      .filter((k) => k === key)
    return sameKey.length === 1
  }

  const handleSubmit: SubmitHandler<PlainMessage<ApplicationEnvVars>> = async (values) => {
    console.log(values)
    const oldVars = new Map(
      props.envVars.variables.filter((envVar) => !envVar.system).map((envVar) => [envVar.key, envVar.value]),
    )
    const newVars = new Map(
      values.variables
        .filter((envVar) => !envVar.system && envVar.key !== '')
        .map((envVar) => [envVar.key, envVar.value]),
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
    <EnvVar.Form onSubmit={handleSubmit}>
      <FormBox.Container>
        <FormBox.Forms>
          <EnvVarsContainer>
            <div>Key</div>
            <div>Value</div>
            <EnvVar.FieldArray name="variables">
              {(fieldArray) => (
                <For each={fieldArray.items}>
                  {(_, index) => (
                    <>
                      <EnvVar.Field name={`variables.${index()}.system`} type="boolean">
                        {/*
                            To make a field active, it must be included in the DOM
                            see: https://modularforms.dev/solid/guides/add-fields-to-form#active-state
                          */}
                        {() => <></>}
                      </EnvVar.Field>
                      <EnvVar.Field
                        name={`variables.${index()}.key`}
                        validate={[
                          custom(isUniqueKey, '同じキーの環境変数が存在します'),
                          (val) => (val === '' && index() !== fieldArray.items.length - 1 ? 'Please enter a key' : ''),
                        ]}
                      >
                        {(field, fieldProps) => (
                          <TextInput
                            value={field.value}
                            error={field.error}
                            {...fieldProps}
                            readOnly={getValue(envVarForm, `variables.${index()}.system`)}
                            disabled={getValue(envVarForm, `variables.${index()}.system`)}
                            tooltip={{
                              props: {
                                content: 'システム環境変数は変更できません',
                              },
                              disabled: !getValue(envVarForm, `variables.${index()}.system`),
                            }}
                          />
                        )}
                      </EnvVar.Field>
                      <EnvVar.Field name={`variables.${index()}.value`}>
                        {(field, fieldProps) => (
                          <TextInput
                            value={field.value}
                            {...fieldProps}
                            readOnly={getValue(envVarForm, `variables.${index()}.system`)}
                            disabled={getValue(envVarForm, `variables.${index()}.system`)}
                            tooltip={{
                              props: {
                                content: 'システム環境変数は変更できません',
                              },
                              disabled: !getValue(envVarForm, `variables.${index()}.system`),
                            }}
                          />
                        )}
                      </EnvVar.Field>
                    </>
                  )}
                </For>
              )}
            </EnvVar.FieldArray>
          </EnvVarsContainer>
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={envVarForm.dirty && !envVarForm.submitting}>
            <Button variants="borderError" size="small" type="button" onClick={discardChanges}>
              Discard Changes
            </Button>
          </Show>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={envVarForm.invalid || !envVarForm.dirty || envVarForm.submitting}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </EnvVar.Form>
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
