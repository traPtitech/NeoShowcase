import { styled } from '/@/components/styled-components'

const Container = styled('div', 'w-full rounded-lg border border-ui-border border-solid bg-ui-primary')

const Forms = styled('div', 'flex w-full flex-col gap-6 px-6 py-5')

const Actions = styled(
  'div',
  'flex w-full items-center justify-end gap-2 rounded-b-lg border-ui-border border-t bg-ui-secondary px-6 py-4',
)

const FormBox = {
  Container,
  Forms,
  Actions,
}

export default FormBox
