import { FieldArray, getInput, insert, reset } from '@formisch/solid/v2'
import { type Component, For, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { List } from '/@/components/templates/List'
import { Button } from '/@/components/UI/Button'
import { systemInfo } from '/@/libs/session'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'
import { getInitialValueOfCreateWebsiteForm } from '../../../schema/websiteSchema'
import WebsiteFieldGroup from './WebsiteFieldGroup'

const FieldRow = styled(
  'div',
  'w-full px-6 py-5 [&:not(:first-child)]:border-t [&:not(:first-child)]:border-t-ui-border',
)

type Props = {
  form: ApplicationFormStore
  isRuntimeApp: boolean
  readonly?: boolean
}

/** ウェブサイト (URL) 設定の一覧。呼び出し側で枠を用意する。 */
const WebsiteListField: Component<Props> = (props) => {
  const addWebsite = () => {
    const defaultDomain = systemInfo().domains.at(0)
    if (!defaultDomain) throw new Error('Default domain is not found')
    const index = getInput(props.form, { path: ['form', 'websites'] })?.length ?? 0
    insert(props.form, { path: ['form', 'websites'], initialInput: getInitialValueOfCreateWebsiteForm(defaultDomain) })
    // Once the form has been submitted (create wizard: "Next"), Formisch revalidates everything on insert and the
    // untouched row would start out full of "required" errors. Reset only the new row's state.
    reset(props.form, { path: ['form', 'websites', index] })
  }

  return (
    <FieldArray of={props.form} path={['form', 'websites']}>
      {(fieldArray) => (
        <>
          <For
            each={fieldArray.items}
            fallback={
              <List.PlaceHolder>
                <div class="i-material-symbols:link-off text-20/20 shrink-0" />
                URLが設定されていません
                <Button
                  variants="primary"
                  size="medium"
                  rightIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
                  onClick={addWebsite}
                  type="button"
                >
                  Add URL
                </Button>
              </List.PlaceHolder>
            }
          >
            {(_, index) => (
              <FieldRow>
                <WebsiteFieldGroup
                  form={props.form}
                  index={index()}
                  isRuntimeApp={props.isRuntimeApp}
                  readonly={props.readonly}
                />
              </FieldRow>
            )}
          </For>
          <Show when={fieldArray.items.length > 0}>
            <FieldRow>
              <div class="flex w-full items-center justify-center">
                <Button
                  onClick={addWebsite}
                  variants="border"
                  size="small"
                  leftIcon={<div class="i-material-symbols:add shrink-0 text-xl/5" />}
                  type="button"
                >
                  Add More
                </Button>
              </div>
            </FieldRow>
          </Show>
        </>
      )}
    </FieldArray>
  )
}

export default WebsiteListField
