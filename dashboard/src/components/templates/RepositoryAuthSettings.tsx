import {
  CreateRepositoryAuth,
  CreateRepositoryRequest,
  UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, FormStore, ValidateField, getValue, required, setValue } from '@modular-forms/solid'
import { Match, Show, Switch, createEffect, createSignal } from 'solid-js'
import { createResource } from 'solid-js'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { TextInput } from '../UI/TextInput'
import { CopyableInput } from './CopyableInput'
import { FormItem } from './FormItem'
import { RadioButtons, RadioItem } from './RadioButtons'

const SshKeyContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const RefreshButtonContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    color: colorVars.semantic.accent.error,
    ...textVars.caption.regular,
  },
})
const VisibilityButton = styled('button', {
  base: {
    width: '24px',
    height: '24px',
    padding: '0',
    background: 'none',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',

    color: colorVars.semantic.text.black,
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:active': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
    },
  },
})

const AuthMethods: RadioItem<CreateRepositoryAuth['auth']['case']>[] = [
  { title: 'SSH', value: 'ssh' },
  { title: 'HTTPS', value: 'basic' },
  { title: '認証を使用しない', value: 'none' },
]

type AuthMethods = {
  [K in Exclude<PlainMessage<CreateRepositoryAuth>['auth']['case'], undefined>]: Extract<
    PlainMessage<CreateRepositoryAuth>['auth'],
    { case: K }
  >['value']
}

export type AuthForm = {
  url: PlainMessage<UpdateRepositoryRequest | CreateRepositoryRequest>['url']
  case: PlainMessage<CreateRepositoryAuth>['auth']['case']
  auth: AuthMethods
}

export const formToAuth = <T extends AuthForm>(form: T): PlainMessage<CreateRepositoryAuth>['auth'] => {
  const authMethod = form.case
  switch (authMethod) {
    case 'none':
      return {
        case: 'none',
        value: '',
      }
    case 'basic':
      return {
        case: 'basic',
        value: {
          username: form.auth.basic.username,
          password: form.auth.basic.password,
        },
      }
    case 'ssh':
      return {
        case: 'ssh',
        value: {
          keyId: form.auth.ssh.keyId,
        },
      }
  }
  throw new Error('unreachable')
}

interface Props {
  formStore: FormStore<AuthForm, undefined>
  hasPermission: boolean
}

export const RepositoryAuthSettings = (props: Props) => {
  const [showPassword, setShowPassword] = createSignal(false)
  const [useTmpKey, setUseTmpKey] = createSignal(false)
  const [tmpKey] = createResource(
    () => (useTmpKey() ? true : undefined),
    () => client.generateKeyPair({}),
  )
  createEffect(() => {
    if (!tmpKey()) return
    setValue(props.formStore, 'auth.ssh.keyId', tmpKey()?.keyId)
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemInfo()?.publicKey ?? '')

  const AuthMethod = () => (
    <Field of={props.formStore} name="case">
      {(field, fieldProps) => (
        <FormItem title="認証方法" error={field.error}>
          <RadioButtons
            items={AuthMethods}
            selected={field.value}
            setSelected={(v) => {
              setValue(props.formStore, 'case', v)
            }}
            {...fieldProps}
            readonly={!props.hasPermission}
          />
        </FormItem>
      )}
    </Field>
  )

  const validateUrl: ValidateField<AuthForm['url']> = (url) => {
    if (getValue(props.formStore, 'case') === 'basic' && !url?.startsWith('https')) {
      return 'Basic認証を使用する場合、URLはhttps://から始まる必要があります'
    }
    return ''
  }
  const Url = () => {
    return (
      <Field of={props.formStore} name="url" validate={[required('Enter Repository URL'), validateUrl]}>
        {(field, fieldProps) => (
          <FormItem title="Repository URL" required>
            <TextInput value={field.value} error={field.error} readonly={!props.hasPermission} {...fieldProps} />
          </FormItem>
        )}
      </Field>
    )
  }

  const AuthConfig = () => {
    const authMethod = () => getValue(props.formStore, 'case')
    return (
      <Switch>
        <Match when={authMethod() === 'basic'}>
          <Field of={props.formStore} name="auth.basic.username" validate={required('Enter UserName')}>
            {(field, fieldProps) => (
              <FormItem title="UserName" required>
                <TextInput value={field.value} error={field.error} readonly={!props.hasPermission} {...fieldProps} />
              </FormItem>
            )}
          </Field>
          <Field of={props.formStore} name="auth.basic.password" validate={required('Enter Password')}>
            {(field, fieldProps) => (
              <FormItem title="Password" required>
                <TextInput
                  value={field.value}
                  type={showPassword() ? 'text' : 'password'}
                  error={field.error}
                  readonly={!props.hasPermission}
                  {...fieldProps}
                  rightIcon={
                    <VisibilityButton onClick={() => setShowPassword((s) => !s)} type="button">
                      <Show when={showPassword()} fallback={<MaterialSymbols>visibility_off</MaterialSymbols>}>
                        <MaterialSymbols>visibility</MaterialSymbols>
                      </Show>
                    </VisibilityButton>
                  }
                />
              </FormItem>
            )}
          </Field>
        </Match>
        <Match when={authMethod() === 'ssh'}>
          <Field of={props.formStore} name="auth.ssh.keyId">
            {() => (
              <FormItem title="SSH公開鍵">
                <SshKeyContainer>
                  以下のSSH公開鍵{!useTmpKey() && '(システムデフォルト)'}
                  をリポジトリに登録してください
                  <CopyableInput value={publicKey()} />
                  <Show when={!useTmpKey()}>
                    <RefreshButtonContainer>
                      <Button
                        variants="textError"
                        size="small"
                        onClick={() => {
                          setUseTmpKey(true)
                        }}
                        leftIcon={<MaterialSymbols opticalSize={20}>replay</MaterialSymbols>}
                      >
                        再生成する
                      </Button>
                      For Github.com
                    </RefreshButtonContainer>
                  </Show>
                </SshKeyContainer>
              </FormItem>
            )}
          </Field>
        </Match>
      </Switch>
    )
  }

  return {
    AuthMethod,
    Url,
    AuthConfig,
  }
}
