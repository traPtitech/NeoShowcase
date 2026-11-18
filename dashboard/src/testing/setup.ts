import { vi } from 'vitest'
import { mockTransport } from '/@/testing/mockBackend'

// Route every Connect call of the app to the in-memory backend.
// (mockBackend is imported statically so that the tests share its state: `vi.resetModules()` does not re-run
// mock factories, so the app keeps talking to this very instance.)
vi.mock('@connectrpc/connect-web', () => ({ createConnectTransport: () => mockTransport }))

// jsdom has no canvas 2d context, so chart.js cannot draw. Keep the real exports (scales, controllers, ...)
// but replace the chart itself with a no-op.
vi.mock('chart.js', async (importOriginal) => {
  const actual = await importOriginal<typeof import('chart.js')>()
  class FakeChart {
    static register() {}
    update() {}
    destroy() {}
  }
  return { ...actual, Chart: FakeChart }
})

// jsdom does not implement scrolling (the router scrolls to top on navigation)
window.scrollTo = () => {}
