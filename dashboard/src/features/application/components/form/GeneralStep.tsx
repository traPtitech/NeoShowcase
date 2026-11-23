import { Field, getInput } from '@formisch/solid/v2'
import { type Component, Loading, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { CheckBox } from '/@/components/templates/CheckBox'
import { FormItem } from '/@/components/templates/FormItem'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import GroupSelect from '/@/features/applicationGroup/components/GroupSelect'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import type { ApplicationFormStore } from '../../schema/applicationSchema'
import type { EnvVarFormStore } from '../../schema/envVarSchema'
import BuildTypeField from './config/BuildTypeField'
import ConfigField from './config/ConfigField'
import EnvVarConfigField from './config/EnvVarConfigField'
import BranchField from './general/BranchField'
import NameField from './general/NameField'

const GeneralStep: Component<{
  form: ApplicationFormStore
  envVarForm: EnvVarFormStore
  repo: Repository
  backToRepoStep: () => void
  proceedToWebsiteStep: () => void
}> = (props) => (
  <div class="flex w-full flex-col items-center gap-10">
    <div class="bg-ui-primary flex w-full flex-col gap-5 rounded-lg p-6">
      <h2 class="overflow-wrap-anywhere h2-medium text-text-black flex items-center gap-1">
        Create Application from
        {originToIcon(repositoryURLToOrigin(props.repo.url), 24)}
        {props.repo.name}
      </h2>
      <NameField form={props.form} />
      <Loading fallback={<Skeleton height={48} />}>
        <Field of={props.form} path={['form', 'groupId']}>
          {(field) => (
            <GroupSelect
              name={field.props.name}
              label="Group"
              info={{
                content: (
                  <>
                    <div>同じグループのアプリ同士は URL を入れ子にできます</div>
                    <div>オーナーはグループで管理します</div>
                  </>
                ),
              }}
              value={field.input}
              setValue={(groupId) => field.onInput(groupId)}
            />
          )}
        </Field>
      </Loading>
      <BranchField form={props.form} repo={props.repo} />
      <BuildTypeField form={props.form} />
      <Show
        when={
          getInput(props.form, { path: ['form', 'config', 'deployConfig', 'type'] }) &&
          getInput(props.form, { path: ['form', 'config', 'buildConfig', 'type'] })
        }
      >
        <ConfigField form={props.form} />
      </Show>
      <Field of={props.form} path={['form', 'startOnCreate']}>
        {(field) => (
          <FormItem title="Start Immediately" tooltip={{ content: <div>この設定で今すぐ起動するかどうか</div> }}>
            <CheckBox.Option
              {...field.props}
              label="今すぐ起動する"
              checked={field.input ?? true}
              setChecked={(v) => field.onInput(v)}
              error={field.errors?.[0]}
            />
          </FormItem>
        )}
      </Field>
      <FormItem title="環境変数">
        <EnvVarConfigField form={props.envVarForm} />
      </FormItem>
    </div>
    <div class="flex gap-5">
      <Button
        size="medium"
        variants="border"
        type="button"
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
        disabled={props.form.isSubmitting}
        loading={props.form.isSubmitting}
      >
        Next
      </Button>
    </div>
  </div>
)

export default GeneralStep
