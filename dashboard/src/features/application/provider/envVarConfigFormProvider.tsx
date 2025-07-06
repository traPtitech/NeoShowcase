import { useFormContext } from '/@/libs/useFormContext'
import { envVarSchema } from '../schema/envVarSchema'

export const { FormProvider: EnvVarConfigFormProvider, useForm: useEnvVarConfigForm } = useFormContext(envVarSchema)
