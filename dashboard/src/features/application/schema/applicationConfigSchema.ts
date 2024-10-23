import type { PartialMessage } from '@bufbuild/protobuf'
import { P, match } from 'ts-pattern'
import * as v from 'valibot'
import type { ApplicationConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { stringBooleanSchema } from '/@/libs/schemaUtil'

const optionalBooleanSchema = (defaultValue = false) =>
  v.pipe(
    v.optional(v.boolean()),
    v.transform((i) => i ?? defaultValue),
  )

const runtimeConfigSchema = v.object({
  useMariadb: optionalBooleanSchema(),
  useMongodb: optionalBooleanSchema(),
  entrypoint: v.string(),
  command: v.string(),
  autoShutdown: v.boolean(),
})
const staticConfigSchema = v.object({
  artifactPath: v.pipe(v.string(), v.nonEmpty('Enter Artifact Path')),
  spa: v.optional(stringBooleanSchema, 'false'),
})

const deployConfigSchema = v.pipe(
  v.optional(
    v.variant(
      'type',
      [
        v.object({
          type: v.literal('runtime'),
          value: v.object({
            runtime: runtimeConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('static'),
          value: v.object({
            static: staticConfigSchema,
          }),
        }),
      ],
      'Select Deploy Type',
    ),
  ),
  // アプリ作成時には最初undefinedになっているが、submit時にはundefinedで無い必要がある
  v.check((input) => !!input, 'Select Deploy Type'),
)

const buildpackConfigSchema = v.object({
  context: v.string(),
})
const cmdConfigSchema = v.object({
  baseImage: v.string(),
  buildCmd: v.string(),
})
const dockerfileConfigSchema = v.object({
  dockerfileName: v.pipe(v.string(), v.nonEmpty('Enter Dockerfile Name')),
  context: v.string(),
})

const buildConfigSchema = v.pipe(
  v.optional(
    v.variant(
      'type',
      [
        v.object({
          type: v.literal('buildpack'),
          value: v.object({
            buildpack: buildpackConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('cmd'),
          value: v.object({
            cmd: cmdConfigSchema,
          }),
        }),
        v.object({
          type: v.literal('dockerfile'),
          value: v.object({
            dockerfile: dockerfileConfigSchema,
          }),
        }),
      ],
      'Select Build Type',
    ),
  ),
  // アプリ作成時には最初undefinedになっているが、submit時にはundefinedで無い必要がある
  v.check((input) => !!input, 'Select Build Type'),
)

export const applicationConfigSchema = v.pipe(
  v.object({
    deployConfig: deployConfigSchema,
    buildConfig: buildConfigSchema,
  }),
  v.transform((input): PartialMessage<ApplicationConfig> => {
    return match([input.deployConfig, input.buildConfig])
      .returnType<PartialMessage<ApplicationConfig>>()
      .with([{ type: 'runtime' }, { type: 'buildpack' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeBuildpack',
            value: {
              ...buildConfig.value.buildpack,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      })
      .with([{ type: 'runtime' }, { type: 'cmd' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeCmd',
            value: {
              ...buildConfig.value.cmd,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      })
      .with([{ type: 'runtime' }, { type: 'dockerfile' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'runtimeDockerfile',
            value: {
              ...buildConfig.value.dockerfile,
              runtimeConfig: deployConfig.value.runtime,
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'buildpack' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticBuildpack',
            value: {
              ...buildConfig.value.buildpack,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'cmd' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticCmd',
            value: {
              ...buildConfig.value.cmd,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'dockerfile' }], ([deployConfig, buildConfig]) => {
        return {
          buildConfig: {
            case: 'staticDockerfile',
            value: {
              ...buildConfig.value.dockerfile,
              staticConfig: deployConfig.value.static,
            },
          },
        }
      })
      .with(P.union([undefined, P._], [P._, undefined]), () => ({
        buildConfig: undefined,
      }))
      .exhaustive()
  }),
)

export type ApplicationConfigInput = v.InferInput<typeof applicationConfigSchema>

/** protobuf message -> valobot schema input */
export const configMessageToSchema = (config: ApplicationConfig): ApplicationConfigInput => {
  const deployConfig = match(config.buildConfig)
    .returnType<ApplicationConfigInput['deployConfig']>()
    .with(
      {
        case: P.union('runtimeBuildpack', 'runtimeDockerfile', 'runtimeCmd'),
      },
      (buildConfig) => ({
        type: 'runtime',
        value: {
          runtime: buildConfig.value.runtimeConfig ?? {
            command: '',
            entrypoint: '',
            useMariadb: false,
            useMongodb: false,
            autoShutdown: false,
          },
        },
      }),
    )
    .with(
      {
        case: P.union('staticBuildpack', 'staticDockerfile', 'staticCmd'),
      },
      (buildConfig) => ({
        type: 'static',
        value: {
          static: buildConfig.value.staticConfig
            ? {
                spa: buildConfig.value.staticConfig.spa ? 'true' : 'false',
                artifactPath: buildConfig.value.staticConfig.artifactPath,
              }
            : {
                spa: 'false',
                artifactPath: '',
              },
        },
      }),
    )
    .with({ case: undefined }, () => undefined)
    .exhaustive()

  const buildConfig = match(config.buildConfig)
    .returnType<ApplicationConfigInput['buildConfig']>()
    .with(
      {
        case: P.union('runtimeBuildpack', 'staticBuildpack'),
      },
      (buildConfig) => ({
        type: 'buildpack',
        value: {
          buildpack: buildConfig.value,
        },
      }),
    )
    .with(
      {
        case: P.union('runtimeCmd', 'staticCmd'),
      },
      (buildConfig) => ({
        type: 'cmd',
        value: {
          cmd: buildConfig.value,
        },
      }),
    )
    .with(
      {
        case: P.union('runtimeDockerfile', 'staticDockerfile'),
      },
      (buildConfig) => ({
        type: 'dockerfile',
        value: {
          dockerfile: buildConfig.value,
        },
      }),
    )
    .with({ case: undefined }, () => undefined)
    .exhaustive()

  return {
    deployConfig,
    buildConfig,
  }
}
