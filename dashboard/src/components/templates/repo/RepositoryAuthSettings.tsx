import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, FormStore, ValidateField, getValue, required, setValue } from '@modular-forms/solid'
import { Match, Show, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { Suspense } from 'solid-js'
import {
  CreateRepositoryAuth,
  CreateRepositoryRequest,
  UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextField } from '/@/components/UI/TextField'
import { client, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { TooltipInfoIcon } from '../../UI/TooltipInfoIcon'
import { FormItem } from '../FormItem'
import { RadioGroup, RadioOption } from '../RadioGroups'

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
    width: '40px',
    height: '40px',
    padding: '8px',
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

const AuthMethods: RadioOption<Exclude<CreateRepositoryAuth['auth']['case'], undefined>>[] = [
  { label: '認証を使用しない', value: 'none' },
  { label: 'BASIC認証', value: 'basic' },
  { label: 'SSH公開鍵認証', value: 'ssh' },
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
    if (tmpKey.latest !== undefined) {
      setValue(props.formStore, 'auth.ssh.keyId', tmpKey().keyId)
    }
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemInfo()?.publicKey ?? '')

  const AuthMethod = () => (
    <Field of={props.formStore} name="case">
      {(field, fieldProps) => (
        <RadioGroup
          label="認証方法"
          {...fieldProps}
          options={AuthMethods}
          value={field.value}
          error={field.error}
          readOnly={!props.hasPermission}
        />
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
          <TextField
            label="Repository URL"
            required
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
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
              <TextField
                label="UserName"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <Field of={props.formStore} name="auth.basic.password" validate={required('Enter Password')}>
            {(field, fieldProps) => (
              <TextField
                label="Password"
                required
                type={showPassword() ? 'text' : 'password'}
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
                rightIcon={
                  <VisibilityButton onClick={() => setShowPassword((s) => !s)} type="button">
                    <Show when={showPassword()} fallback={<MaterialSymbols>visibility_off</MaterialSymbols>}>
                      <MaterialSymbols>visibility</MaterialSymbols>
                    </Show>
                  </VisibilityButton>
                }
              />
            )}
          </Field>
        </Match>
        <Match when={authMethod() === 'ssh'}>
          <Field of={props.formStore} name="auth.ssh.keyId">
            {() => (
              <FormItem title="デプロイキーの登録">
                <Suspense>
                  <SshKeyContainer>
                    以下のSSH公開鍵{useTmpKey() ? '(このリポジトリ専用)' : '(NeoShowcase全体共通)'}
                    を、リポジトリのデプロイキーとして登録してください。
                    <br />
                    公開リポジトリの場合は、この操作は不要です。
                    <TextField value={publicKey()} copyable={true} readonly />
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
                          専用公開鍵を生成する
                        </Button>
                        <TooltipInfoIcon
                          props={{
                            content: (
                              <>
                                <div>このリポジトリ専用のSSH用鍵ペアを生成します。</div>
                                <div>
                                  NeoShowcase全体で共通の公開鍵が、リポジトリに登録できない場合に生成してください。
                                </div>
                                <div>GitHubプライベートリポジトリの場合は必ず生成が必要です。</div>
                              </>
                            ),
                          }}
                          style="left"
                        />
                      </RefreshButtonContainer>
                    </Show>
                  </SshKeyContainer>
                </Suspense>
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
