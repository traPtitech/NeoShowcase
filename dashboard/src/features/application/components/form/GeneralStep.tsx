import { Field, getValue } from '@modular-forms/solid'
import { type Component, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { colorVars, textVars } from '/@/theme'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import BuildTypeField from './config/BuildTypeField'
import ConfigField from './config/ConfigField'
import BranchField from './general/BranchField'
import NameField from './general/NameField'

const GeneralStep: Component<{
  repo: Repository
  backToRepoStep: () => void
  proceedToWebsiteStep: () => void
}> = (props) => {
  const { formStore } = useApplicationForm()

  return (
    <div class="flex w-full flex-col items-center gap-10">
      <div class="flex w-full flex-col gap-5 rounded-lg bg-ui-primary p-6">
        <h2 class="overflow-wrap-anywhere h2-medium flex items-center gap-1 text-text-black">
          Create Application from
          {originToIcon(repositoryURLToOrigin(props.repo.url), 24)}
          {props.repo.name}
        </h2>
        <NameField />
        <BranchField repo={props.repo} />
        <BuildTypeField />
        <Show
          when={
            getValue(formStore, 'form.config.deployConfig.type') && getValue(formStore, 'form.config.buildConfig.type')
          }
        >
          <ConfigField />
        </Show>
        <Field of={formStore} name="form.startOnCreate" type="boolean">
          {(field, fieldProps) => (
            <FormItem
              title="Start Immediately"
              tooltip={{
                props: {
                  content: (
                    <>
                      <div>この設定で今すぐ起動するかどうか</div>
                      <div>(環境変数はアプリ作成後設定可能になります)</div>
                    </>
                  ),
                },
              }}
            >
              <CheckBox.Option
                {...fieldProps}
                label="今すぐ起動する"
                checked={field.value ?? false}
                error={field.error}
              />
            </FormItem>
          )}
        </Field>
      </div>
      <div class="flex gap-5">
        <Button
          size="medium"
          variants="border"
          onClick={props.backToRepoStep}
          leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}
        >
          Back
        </Button>
        <Button
          size="medium"
          variants="primary"
          type="button"
          onClick={props.proceedToWebsiteStep}
          rightIcon={<MaterialSymbols>arrow_forward</MaterialSymbols>}
          disabled={formStore.invalid || formStore.submitting}
          loading={formStore.submitting}
        >
          Next
        </Button>
      </div>
    </div>
  )
}
export default GeneralStep
