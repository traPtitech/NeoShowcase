import { Title } from '@solidjs/meta'
import ErrorView from '/@/components/layouts/ErrorView'

export default () => (
  <>
    <Title>Not Found - NeoShowcase</Title>
    <ErrorView error={new Error('Not Found')} />
  </>
)
