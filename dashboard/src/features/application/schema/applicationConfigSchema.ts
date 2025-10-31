import { match, P } from 'ts-pattern'
import * as v from 'valibot'
import { type ApplicationConfig, AutoShutdownConfig_StartupBehavior } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { stringBooleanSchema } from '/@/libs/schemaUtil'

const unwrapStartupBehaviorMap: Record<`${AutoShutdownConfig_StartupBehavior}`, AutoShutdownConfig_StartupBehavior> = {
  [AutoShutdownConfig_StartupBehavior.UNDEFINED]: AutoShutdownConfig_StartupBehavior.UNDEFINED,
  [AutoShutdownConfig_StartupBehavior.LOADING_PAGE]: AutoShutdownConfig_StartupBehavior.LOADING_PAGE,
  [AutoShutdownConfig_StartupBehavior.BLOCKING]: AutoShutdownConfig_StartupBehavior.BLOCKING,
} as const

const autoShutdownSchema = v.optional(
  v.object({
    enabled: v.boolean(),
    startup: v.pipe(
      v.optional(
        v.union([
          v.literal(`${AutoShutdownConfig_StartupBehavior.LOADING_PAGE}`),
          v.literal(`${AutoShutdownConfig_StartupBehavior.BLOCKING}`),
          v.literal(`${AutoShutdownConfig_StartupBehavior.UNDEFINED}`),
        ]),
        `${AutoShutdownConfig_StartupBehavior.UNDEFINED}`,
      ),
      v.transform((input) => unwrapStartupBehaviorMap[input]),
    ),
  }),
)

const runtimeConfigSchema = v.object({
  useMariadb: v.optional(v.boolean(), false),
  useMongodb: v.optional(v.boolean(), false),
  entrypoint: v.string(),
  command: v.string(),
  autoShutdown: autoShutdownSchema,
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
  v.transform((input): ApplicationConfig => {
    return match([input.deployConfig, input.buildConfig])
      .returnType<ApplicationConfig>()
      .with([{ type: 'runtime' }, { type: 'buildpack' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'runtimeBuildpack',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigRuntimeBuildpack',
              ...buildConfig.value.buildpack,
              runtimeConfig: {
                $typeName: 'neoshowcase.protobuf.RuntimeConfig',
                ...deployConfig.value.runtime,
                autoShutdown: deployConfig.value.runtime.autoShutdown
                  ? {
                      $typeName: 'neoshowcase.protobuf.AutoShutdownConfig',
                      enabled: deployConfig.value.runtime.autoShutdown.enabled,
                      startup: deployConfig.value.runtime.autoShutdown.startup,
                    }
                  : undefined,
              },
            },
          },
        }
      })
      .with([{ type: 'runtime' }, { type: 'cmd' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'runtimeCmd',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigRuntimeCmd',
              ...buildConfig.value.cmd,
              runtimeConfig: {
                $typeName: 'neoshowcase.protobuf.RuntimeConfig',
                ...deployConfig.value.runtime,
                autoShutdown: deployConfig.value.runtime.autoShutdown
                  ? {
                      $typeName: 'neoshowcase.protobuf.AutoShutdownConfig',
                      enabled: deployConfig.value.runtime.autoShutdown.enabled,
                      startup: deployConfig.value.runtime.autoShutdown.startup,
                    }
                  : undefined,
              },
            },
          },
        }
      })
      .with([{ type: 'runtime' }, { type: 'dockerfile' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'runtimeDockerfile',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigRuntimeDockerfile',
              ...buildConfig.value.dockerfile,
              runtimeConfig: {
                $typeName: 'neoshowcase.protobuf.RuntimeConfig',
                ...deployConfig.value.runtime,
                autoShutdown: deployConfig.value.runtime.autoShutdown
                  ? {
                      $typeName: 'neoshowcase.protobuf.AutoShutdownConfig',
                      enabled: deployConfig.value.runtime.autoShutdown.enabled,
                      startup: deployConfig.value.runtime.autoShutdown.startup,
                    }
                  : undefined,
              },
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'buildpack' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'staticBuildpack',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigStaticBuildpack',
              ...buildConfig.value.buildpack,
              staticConfig: {
                $typeName: 'neoshowcase.protobuf.StaticConfig',
                ...deployConfig.value.static,
              },
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'cmd' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'staticCmd',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigStaticCmd',
              ...buildConfig.value.cmd,
              staticConfig: {
                $typeName: 'neoshowcase.protobuf.StaticConfig',
                ...deployConfig.value.static,
              },
            },
          },
        }
      })
      .with([{ type: 'static' }, { type: 'dockerfile' }], ([deployConfig, buildConfig]) => {
        return {
          $typeName: 'neoshowcase.protobuf.ApplicationConfig',
          buildConfig: {
            case: 'staticDockerfile',
            value: {
              $typeName: 'neoshowcase.protobuf.BuildConfigStaticDockerfile',
              ...buildConfig.value.dockerfile,
              staticConfig: {
                $typeName: 'neoshowcase.protobuf.StaticConfig',
                ...deployConfig.value.static,
              },
            },
          },
        }
      })
      .with(P.union([undefined, P._], [P._, undefined]), () => ({
        $typeName: 'neoshowcase.protobuf.ApplicationConfig',
        buildConfig: {
          case: undefined,
        },
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
          runtime: {
            ...buildConfig.value.runtimeConfig,
            entrypoint: buildConfig.value.runtimeConfig?.entrypoint ?? '',
            command: buildConfig.value.runtimeConfig?.command ?? '',
            autoShutdown: {
              enabled: buildConfig.value.runtimeConfig?.autoShutdown?.enabled ?? false,
              startup: buildConfig.value.runtimeConfig?.autoShutdown?.startup
                ? `${buildConfig.value.runtimeConfig.autoShutdown.startup}`
                : undefined,
            },
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
    .with(
      {
        case: P.union('functionBuildpack', 'functionDockerfile', 'functionCmd'),
      },
      (buildConfig) => {
        throw new Error(`Currently function build config type ${buildConfig.case} is not supported`)
      },
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
    .with(
      {
        case: P.union('functionBuildpack', 'functionCmd', 'functionDockerfile'),
      },
      (buildConfig) => {
        throw new Error(`Currently function build config type ${buildConfig.case} is not supported`)
      },
    )
    .with({ case: undefined }, () => undefined)
    .exhaustive()

  return {
    deployConfig,
    buildConfig,
  }
}
