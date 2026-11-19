import { Field, getInput } from '@formisch/solid/v2'
import type { ComponentProps } from '@solidjs/web'
import { type Component, For, omit, type ParentComponent, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { FormItem } from '/@/components/templates/FormItem'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { TextField } from '/@/components/UI/TextField'
import { websiteWarnings } from '/@/libs/application'
import { inputToNumber, numberToInput } from '/@/libs/schemaUtil'
import { systemInfo } from '/@/libs/session'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'

const URLContainer = styled('div', 'flex items-start gap-2')

const URLItem: ParentComponent<ComponentProps<'div'> & { fixedWidth?: boolean }> = (props) => {
  const extraProps = omit(props, 'fixedWidth', 'class')
  return <div class={['flex h-12 items-center', props.fixedWidth && 'w-[calc(6ch+60px)] shrink-0']} {...extraProps} />
}

const schemeOptions: SelectOption<`${boolean}`>[] = [
  { value: 'false', label: 'http' },
  { value: 'true', label: 'https' },
]

type Props = {
  form: ApplicationFormStore
  index: number
  showHttpPort: boolean
  readonly?: boolean
}

const UrlField: Component<Props> = (props) => {
  const selectedDomain = () => getInput(props.form, { path: ['form', 'websites', props.index, 'domain'] })
  // 占有されているドメインはoptionに表示しない
  // すでに設定されているドメインはoptionに表示する
  const domainOptions = () =>
    systemInfo()
      .domains.filter((domain) => !domain.alreadyBound || selectedDomain() === domain.domain)
      .map((domain) => ({ value: domain.domain, label: domain.domain.replace(/\*/g, '') }))

  const warnings = () =>
    websiteWarnings(
      getInput(props.form, { path: ['form', 'websites', props.index, 'subdomain'] }),
      getInput(props.form, { path: ['form', 'websites', props.index, 'https'] }) === 'true',
    )

  return (
    <FormItem title="URL" required>
      <URLContainer>
        <URLItem fixedWidth>
          <Field of={props.form} path={['form', 'websites', props.index, 'https']}>
            {(field) => (
              <SingleSelect
                tooltip={{
                  content: (
                    <>
                      <div>スキーム</div>
                      <div>通常はhttpsが推奨です</div>
                    </>
                  ),
                }}
                {...field.props}
                options={schemeOptions}
                value={field.input}
                setValue={(v) => field.onInput(v)}
                readOnly={props.readonly}
              />
            )}
          </Field>
        </URLItem>
        <URLItem>://</URLItem>
        <Show when={selectedDomain()?.startsWith('*')}>
          <Field of={props.form} path={['form', 'websites', props.index, 'subdomain']}>
            {(field) => (
              <TextField
                placeholder="subdomain"
                tooltip={{ content: 'サブドメイン名' }}
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
                readOnly={props.readonly}
              />
            )}
          </Field>
        </Show>
        <Field of={props.form} path={['form', 'websites', props.index, 'domain']}>
          {(field) => (
            <SingleSelect
              tooltip={{ content: 'ドメイン名' }}
              {...field.props}
              options={domainOptions()}
              value={field.input}
              setValue={(v) => field.onInput(v)}
              readOnly={props.readonly}
            />
          )}
        </Field>
      </URLContainer>
      <URLContainer>
        <URLItem>/</URLItem>
        <Field of={props.form} path={['form', 'websites', props.index, 'pathPrefix']}>
          {(field) => (
            <TextField
              tooltip={{ content: '(Advanced) 指定Prefixが付いていたときのみアプリへルーティング' }}
              {...field.props}
              value={field.input ?? ''}
              error={field.errors?.[0]}
              readOnly={props.readonly}
            />
          )}
        </Field>
        <Show when={props.showHttpPort}>
          <URLItem> → </URLItem>
          <URLItem fixedWidth>
            <Field of={props.form} path={['form', 'websites', props.index, 'httpPort']}>
              {(field) => (
                <TextField
                  placeholder="80"
                  type="number"
                  min="0"
                  tooltip={{ content: 'アプリのHTTP Port番号' }}
                  {...field.props}
                  onInput={(e) => field.onInput(inputToNumber(e.currentTarget.value))}
                  value={numberToInput(field.input)}
                  error={field.errors?.[0]}
                  readOnly={props.readonly}
                />
              )}
            </Field>
          </URLItem>
          <URLItem>/TCP</URLItem>
        </Show>
      </URLContainer>
      <Show when={warnings().length > 0}>
        <div class="flex flex-col gap-1">
          <For each={warnings()}>{(item) => <div class="text-accent-error">{item}</div>}</For>
        </div>
      </Show>
    </FormItem>
  )
}

export default UrlField
