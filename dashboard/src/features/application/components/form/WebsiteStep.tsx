import { getInput } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import type { ApplicationFormStore } from '../../schema/applicationSchema'
import WebsiteListField from './website/WebsiteListField'

const WebsiteStep: Component<{
  form: ApplicationFormStore
  backToGeneralStep: () => void
}> = (props) => (
  <div class="flex w-full flex-col items-center gap-10">
    <div class="flex w-full flex-col items-center gap-6">
      <div class="border-ui-border bg-ui-primary flex w-full flex-col gap-0.25 overflow-hidden rounded-lg border">
        <WebsiteListField
          form={props.form}
          isRuntimeApp={getInput(props.form, { path: ['form', 'config', 'deployConfig', 'type'] }) === 'runtime'}
        />
      </div>
    </div>
    <div class="flex gap-5">
      <Button
        size="medium"
        variants="ghost"
        type="button"
        leftIcon={<div class="i-material-symbols:arrow-back shrink-0 text-2xl/6" />}
        onClick={props.backToGeneralStep}
      >
        Back
      </Button>
      <Button
        type="submit"
        size="medium"
        variants="primary"
        disabled={props.form.isSubmitting}
        loading={props.form.isSubmitting}
      >
        Create Application
      </Button>
    </div>
  </div>
)

export default WebsiteStep
