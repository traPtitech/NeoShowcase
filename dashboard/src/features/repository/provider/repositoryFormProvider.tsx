import { useFormContext } from '../../../libs/useFormContext'
import { createOrUpdateRepositorySchema } from '../schema/repositorySchema'

export const { FormProvider: RepositoryFormProvider, useForm: useRepositoryForm } =
  useFormContext(createOrUpdateRepositorySchema)
