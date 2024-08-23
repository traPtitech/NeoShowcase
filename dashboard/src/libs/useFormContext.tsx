import { type FieldValues, type FormStore, type ValidationMode, createFormStore, valiForm } from '@modular-forms/solid'
import { type ParentComponent, createContext, useContext } from 'solid-js'
import type { BaseIssue, BaseSchema, InferInput } from 'valibot'

type FormContextValue<Schema extends FieldValues> = {
  formStore: FormStore<Schema, undefined>
}

export const useFormContext = <TInput extends FieldValues, TOutput, TIssue extends BaseIssue<unknown>>(
  schema: BaseSchema<TInput, TOutput, TIssue>,
  validationMode?: ValidationMode,
) => {
  const FormContext = createContext<FormContextValue<InferInput<typeof schema>>>()

  const FormProvider: ParentComponent = (props) => {
    const formStore = createFormStore<InferInput<typeof schema>>({
      validate: valiForm(schema),
      revalidateOn: validationMode,
    })

    return (
      <FormContext.Provider
        value={{
          formStore,
        }}
      >
        {props.children}
      </FormContext.Provider>
    )
  }

  const useForm = () => {
    const c = useContext(FormContext)
    if (!c) throw new Error('useForm must be used within a FormProvider')
    return c
  }

  return {
    FormProvider,
    useForm,
  }
}
