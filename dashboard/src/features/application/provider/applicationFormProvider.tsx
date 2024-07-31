import { useFormContext } from '../../../libs/useFormContext'
import { createOrUpdateApplicationSchema } from '../schema/applicationSchema'

export const { FormProvider: ApplicationFormProvider, useForm: useApplicationForm } = useFormContext(
  createOrUpdateApplicationSchema,
)
