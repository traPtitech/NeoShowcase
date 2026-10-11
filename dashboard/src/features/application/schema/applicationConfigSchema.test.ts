import { parse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { AutoShutdownConfig_StartupBehavior } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationConfigSchema } from './applicationConfigSchema'

describe('applicationConfigSchema', () => {
  test('auto shutdown is placed on ApplicationConfig, not on RuntimeConfig', () => {
    const config = parse(applicationConfigSchema, {
      deployConfig: {
        type: 'runtime',
        value: {
          runtime: {
            entrypoint: '',
            command: '',
            autoShutdown: {
              enabled: true,
              startup: `${AutoShutdownConfig_StartupBehavior.BLOCKING}`,
            },
          },
        },
      },
      buildConfig: {
        type: 'buildpack',
        value: { buildpack: { context: '' } },
      },
    })

    expect(config.autoShutdown).toEqual(
      expect.objectContaining({ enabled: true, startup: AutoShutdownConfig_StartupBehavior.BLOCKING }),
    )
    expect(config.buildConfig.case).toBe('runtimeBuildpack')
    expect(config.buildConfig.value).toEqual(
      expect.objectContaining({ runtimeConfig: expect.not.objectContaining({ autoShutdown: expect.anything() }) }),
    )
  })
})
