import { render } from '@solidjs/web'
import { createMemo, Errored, flush, Loading } from 'solid-js'
import { expect, it, vi } from 'vitest'
import { ToolTip } from './ToolTip'

it('does not throw when content depends on a pending async value', async () => {
  let resolve!: (v: boolean) => void
  const allowed = new Promise<boolean>((r) => {
    resolve = r
  })
  const root = document.createElement('div')
  document.body.append(root)
  const dispose = render(() => {
    const hasPermission = createMemo(() => allowed)
    return (
      <Errored fallback={(e) => <p id="error">{String(e())}</p>}>
        <Loading fallback={<p id="loading" />}>
          <ToolTip content={hasPermission() ? undefined : 'owner only'}>
            <button type="button">Save</button>
          </ToolTip>
        </Loading>
      </Errored>
    )
  }, root)
  flush()
  expect(root.querySelector('#loading')).not.toBeNull()
  resolve(false)
  await vi.waitFor(() => {
    flush()
    expect(root.querySelector('button')?.textContent).toBe('Save')
  })
  expect(root.querySelector('#error')).toBeNull()
  dispose()
})
