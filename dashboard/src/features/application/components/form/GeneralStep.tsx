import { Field, getValue, setValues } from '@modular-forms/solid'
import { type Component, onMount, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { Button } from '/@/components/UI/Button'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import { useEnvVarConfigForm } from '../../provider/envVarConfigFormProvider'
import BuildTypeField from './config/BuildTypeField'
import ConfigField from './config/ConfigField'
import EnvVarConfigField from './config/EnvVarConfigField'
import BranchField from './general/BranchField'
import NameField from './general/NameField'

const GeneralStep: Component<{
  repo: Repository
  backToRepoStep: () => void
  proceedToWebsiteStep: () => void
}> = (props) => {
  const { formStore } = useApplicationForm()
  const { formStore: envVarConfigFormStore } = useEnvVarConfigForm()

  onMount(() => {
    setValues(envVarConfigFormStore, 'variables', [])
  })

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
                    </>
                  ),
                },
              }}
            >
              <CheckBox.Option
                {...fieldProps}
                label="今すぐ起動する"
                checked={field.value ?? true}
                error={field.error}
              />
            </FormItem>
          )}
        </Field>

        <FormItem title="環境変数">
          <EnvVarConfigField />
        </FormItem>
      </div>
      <div class="flex gap-5">
        <Button
          size="medium"
          variants="border"
          onClick={props.backToRepoStep}
          leftIcon={<div class="i-material-symbols:arrow-back shrink-0 text-2xl/6" />}
        >
          Back
        </Button>
        <Button
          size="medium"
          variants="primary"
          type="button"
          onClick={props.proceedToWebsiteStep}
          rightIcon={<div class="i-material-symbols:arrow-forward shrink-0 text-2xl/6" />}
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
