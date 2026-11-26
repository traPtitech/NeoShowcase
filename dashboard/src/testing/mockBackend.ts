import { clone, create } from '@bufbuild/protobuf'
import { type Timestamp, TimestampSchema } from '@bufbuild/protobuf/wkt'
import { Code, ConnectError, type ConnectRouter, createRouterTransport, type ServiceImpl } from '@connectrpc/connect'
import {
  AdditionalLinkSchema,
  APIService,
  type Application,
  Application_ContainerState,
  ApplicationConfigSchema,
  ApplicationEnvVarSchema,
  ApplicationEnvVarsSchema,
  ApplicationMetricSchema,
  ApplicationMetricsSchema,
  ApplicationOutputSchema,
  ApplicationOutputsSchema,
  ApplicationSchema,
  ArtifactContentSchema,
  AuthenticationType,
  AutoShutdownConfig_StartupBehavior,
  AutoShutdownConfigSchema,
  AvailableDomainSchema,
  AvailableMetricsSchema,
  AvailablePortSchema,
  type Build,
  BuildConfigRuntimeBuildpackSchema,
  BuildConfigStaticCmdSchema,
  BuildLogSchema,
  BuildSchema,
  BuildStatus,
  DeployType,
  GenerateKeyPairResponseSchema,
  GetApplicationsRequest_Scope,
  GetApplicationsResponseSchema,
  GetBuildsResponseSchema,
  GetRepositoriesResponseSchema,
  GetRepositoryCommitsResponseSchema,
  GetRepositoryRefsResponseSchema,
  GetUserKeysResponseSchema,
  GetUsersResponseSchema,
  GitRefSchema,
  PortPublicationProtocol,
  PortPublicationSchema,
  type Repository,
  Repository_AuthMethod,
  RepositorySchema,
  RuntimeConfigSchema,
  SimpleCommitSchema,
  SSHInfoSchema,
  StaticConfigSchema,
  SystemInfoSchema,
  type UserKey,
  UserKeySchema,
  UserSchema,
  WebsiteSchema,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { NullTimestampSchema } from '/@/api/neoshowcase/protobuf/null_pb'

/**
 * In-memory backend for tests: a Connect transport backed by the fixtures below.
 * Mutations update the fixtures, and every request is recorded in `requests`.
 */

const ts = (iso: string): Timestamp => {
  const ms = Date.parse(iso)
  return create(TimestampSchema, { seconds: BigInt(Math.floor(ms / 1000)), nanos: 0 })
}
const nullTs = (iso: string) => create(NullTimestampSchema, { timestamp: ts(iso), valid: true })

// --- Fixtures

export const ids = {
  me: 'user-admin',
  other: 'user-other',
  repoGithub: 'repo-github',
  repoGitea: 'repo-gitea',
  appRuntime: 'app-runtime',
  appStatic: 'app-static',
  appError: 'app-error',
  buildOk: 'build-ok',
  buildFailed: 'build-failed',
  buildRunning: 'build-running',
  buildRunning2: 'build-running-2',
}

export const hashes = {
  a: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
  b: 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb',
  c: 'cccccccccccccccccccccccccccccccccccccccc',
}

const createFixtures = () => {
  const users = [
    create(UserSchema, { id: ids.me, name: 'admin-user', admin: true, avatarUrl: 'https://example.com/a.png' }),
    create(UserSchema, { id: ids.other, name: 'other-user', admin: false, avatarUrl: 'https://example.com/b.png' }),
  ]

  const repositories: Repository[] = [
    create(RepositorySchema, {
      id: ids.repoGithub,
      name: 'github-repo',
      url: 'https://github.com/traPtitech/NeoShowcase.git',
      htmlUrl: 'https://github.com/traPtitech/NeoShowcase',
      authMethod: Repository_AuthMethod.NONE,
      ownerIds: [ids.me],
    }),
    create(RepositorySchema, {
      id: ids.repoGitea,
      name: 'gitea-repo',
      url: 'git@git.trap.jp:traP/gitea-repo.git',
      htmlUrl: 'https://git.trap.jp/traP/gitea-repo',
      authMethod: Repository_AuthMethod.SSH,
      ownerIds: [ids.me, ids.other],
    }),
  ]

  const applications: Application[] = [
    create(ApplicationSchema, {
      id: ids.appRuntime,
      name: 'runtime-app',
      repositoryId: ids.repoGithub,
      refName: 'main',
      commit: hashes.a,
      deployType: DeployType.RUNTIME,
      running: true,
      container: Application_ContainerState.RUNNING,
      currentBuild: ids.buildOk,
      createdAt: ts('2024-01-01T00:00:00Z'),
      updatedAt: ts('2024-01-02T00:00:00Z'),
      config: create(ApplicationConfigSchema, {
        buildConfig: {
          case: 'runtimeBuildpack',
          value: create(BuildConfigRuntimeBuildpackSchema, {
            context: '.',
            runtimeConfig: create(RuntimeConfigSchema, {
              entrypoint: '',
              command: './app',
              useMariadb: true,
            }),
          }),
        },
        autoShutdown: create(AutoShutdownConfigSchema, {
          enabled: true,
          startup: AutoShutdownConfig_StartupBehavior.LOADING_PAGE,
        }),
      }),
      websites: [
        create(WebsiteSchema, {
          id: 'website-1',
          fqdn: 'runtime-app.trap.show',
          pathPrefix: '/',
          stripPrefix: false,
          https: true,
          h2c: false,
          httpPort: 8080,
          authentication: AuthenticationType.SOFT,
        }),
      ],
      portPublications: [
        create(PortPublicationSchema, {
          internetPort: 20000,
          applicationPort: 5000,
          protocol: PortPublicationProtocol.TCP,
        }),
      ],
      ownerIds: [ids.me],
      latestBuildStatus: BuildStatus.SUCCEEDED,
    }),
    create(ApplicationSchema, {
      id: ids.appStatic,
      name: 'static-app',
      repositoryId: ids.repoGithub,
      refName: 'main',
      commit: hashes.b,
      deployType: DeployType.STATIC,
      running: true,
      container: Application_ContainerState.RUNNING,
      currentBuild: ids.buildOk,
      createdAt: ts('2024-02-01T00:00:00Z'),
      updatedAt: ts('2024-02-02T00:00:00Z'),
      config: create(ApplicationConfigSchema, {
        buildConfig: {
          case: 'staticCmd',
          value: create(BuildConfigStaticCmdSchema, {
            baseImage: 'node:22',
            buildCmd: 'npm run build',
            staticConfig: create(StaticConfigSchema, { artifactPath: 'dist', spa: true }),
          }),
        },
      }),
      websites: [],
      portPublications: [],
      ownerIds: [ids.me],
      latestBuildStatus: BuildStatus.SUCCEEDED,
    }),
    create(ApplicationSchema, {
      id: ids.appError,
      name: 'broken-app',
      repositoryId: ids.repoGitea,
      refName: 'develop',
      commit: hashes.c,
      deployType: DeployType.RUNTIME,
      running: false,
      container: Application_ContainerState.ERRORED,
      containerMessage: 'container exited with code 1',
      currentBuild: ids.buildFailed,
      createdAt: ts('2024-03-01T00:00:00Z'),
      updatedAt: ts('2024-03-02T00:00:00Z'),
      config: create(ApplicationConfigSchema, {
        buildConfig: {
          case: 'runtimeBuildpack',
          value: create(BuildConfigRuntimeBuildpackSchema, {
            context: '.',
            runtimeConfig: create(RuntimeConfigSchema, { entrypoint: '', command: '' }),
          }),
        },
      }),
      websites: [],
      portPublications: [],
      ownerIds: [ids.other],
      latestBuildStatus: BuildStatus.FAILED,
    }),
  ]

  const builds: Build[] = [
    create(BuildSchema, {
      id: ids.buildOk,
      applicationId: ids.appRuntime,
      commit: hashes.a,
      status: BuildStatus.SUCCEEDED,
      queuedAt: ts('2024-01-02T00:00:00Z'),
      startedAt: nullTs('2024-01-02T00:00:05Z'),
      updatedAt: nullTs('2024-01-02T00:01:00Z'),
      finishedAt: nullTs('2024-01-02T00:01:00Z'),
      retriable: true,
    }),
    create(BuildSchema, {
      id: ids.buildFailed,
      applicationId: ids.appRuntime,
      commit: hashes.b,
      status: BuildStatus.FAILED,
      queuedAt: ts('2024-01-01T00:00:00Z'),
      startedAt: nullTs('2024-01-01T00:00:05Z'),
      updatedAt: nullTs('2024-01-01T00:00:30Z'),
      finishedAt: nullTs('2024-01-01T00:00:30Z'),
      retriable: true,
    }),
  ]

  builds.push(
    create(BuildSchema, {
      id: ids.buildRunning,
      applicationId: ids.appStatic,
      commit: hashes.b,
      status: BuildStatus.BUILDING,
      queuedAt: ts('2024-02-02T00:00:00Z'),
      startedAt: nullTs('2024-02-02T00:00:05Z'),
      retriable: false,
    }),
    create(BuildSchema, {
      id: ids.buildRunning2,
      applicationId: ids.appStatic,
      commit: hashes.b,
      status: BuildStatus.BUILDING,
      queuedAt: ts('2024-02-01T00:00:00Z'),
      startedAt: nullTs('2024-02-01T00:00:05Z'),
      retriable: false,
    }),
  )

  const commits = [
    create(SimpleCommitSchema, {
      hash: hashes.a,
      authorName: 'alice',
      commitDate: ts('2024-01-02T00:00:00Z'),
      message: 'feat: add feature A',
    }),
    create(SimpleCommitSchema, {
      hash: hashes.b,
      authorName: 'bob',
      commitDate: ts('2024-01-01T00:00:00Z'),
      message: 'fix: repair B',
    }),
    create(SimpleCommitSchema, {
      hash: hashes.c,
      authorName: 'carol',
      commitDate: ts('2024-03-01T00:00:00Z'),
      message: 'chore: tidy C',
    }),
  ]

  const userKeys: UserKey[] = [
    create(UserKeySchema, {
      id: 'key-1',
      userId: ids.me,
      name: 'laptop',
      publicKey: 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIexample laptop',
      createdAt: ts('2024-01-01T00:00:00Z'),
    }),
  ]

  const envVars = new Map<string, Map<string, string>>([[ids.appRuntime, new Map([['FOO', 'bar']])]])

  const systemInfo = create(SystemInfoSchema, {
    publicKey: 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIsystem',
    ssh: create(SSHInfoSchema, { host: 'ssh.trap.show', port: 2222 }),
    domains: [
      create(AvailableDomainSchema, { domain: '*.trap.show', excludeDomains: [], authAvailable: true }),
      create(AvailableDomainSchema, { domain: 'example.com', excludeDomains: [], authAvailable: false }),
    ],
    ports: [
      create(AvailablePortSchema, { startPort: 20000, endPort: 20100, protocol: PortPublicationProtocol.TCP }),
      create(AvailablePortSchema, { startPort: 30000, endPort: 30100, protocol: PortPublicationProtocol.UDP }),
    ],
    additionalLinks: [create(AdditionalLinkSchema, { name: 'Docs', url: 'https://docs.example.com' })],
    version: 'v0.0.0-test',
    revision: 'abcdef0',
  })

  return { users, repositories, applications, builds, commits, userKeys, envVars, systemInfo }
}

export let fixtures = createFixtures()

export type RecordedRequest = { method: string; message: Record<string, unknown> }
/** Every request seen since the last `resetMockBackend()`. */
export const requests: RecordedRequest[] = []
export const requestsTo = (method: string) => requests.filter((r) => r.method === method).map((r) => r.message)

/** Methods that currently fail, with the number of calls still to fail (Infinity: all). */
const failing = new Map<string, number>()

/**
 * Makes `method` (e.g. 'getUserKeys') fail with `message` (default code: Unavailable) until `recover(method)`,
 * or for the next `times` calls.
 */
export const failMethod = (
  method: string,
  message = 'backend is down',
  times = Number.POSITIVE_INFINITY,
  code = Code.Unavailable,
) => {
  failing.set(method, times)
  failureMessages.set(method, message)
  failureCodes.set(method, code)
}
export const recover = (method: string) => failing.delete(method)

/** Number of upcoming `getOutputStream` connections that break after sending their lines. */
let outputStreamDrops = 0
export const dropOutputStream = (times = 1) => {
  outputStreamDrops = times
}

const failureMessages = new Map<string, string>()
const failureCodes = new Map<string, Code>()

/** Unary methods whose calls wait for a gate; `once`: only the next call waits. */
const held = new Map<string, { gate: Promise<void>; once: boolean }>()
const hold = (method: string, once: boolean) => {
  let release!: () => void
  const gate = new Promise<void>((resolve) => {
    release = resolve
  })
  held.set(method, { gate, once })
  return () => {
    if (held.get(method)?.gate === gate) held.delete(method)
    release()
  }
}
/** Holds back the responses of the unary `method` until the returned `release()` is called. */
export const holdMethod = (method: string) => hold(method, false)
/** Holds back only the next call of the unary `method` until the returned `release()` is called. */
export const holdNext = (method: string) => hold(method, true)

export const resetMockBackend = () => {
  fixtures = createFixtures()
  requests.length = 0
  failing.clear()
  failureMessages.clear()
  failureCodes.clear()
  held.clear()
  outputStreamDrops = 0
}

// --- Service implementation

const findOrThrow = <T extends { id: string }>(list: T[], id: string, what: string): T => {
  const found = list.find((x) => x.id === id)
  if (!found) throw new Error(`${what} not found: ${id}`)
  return found
}

const waitUntilAborted = (signal: AbortSignal) =>
  new Promise<void>((resolve) => {
    if (signal.aborted) resolve()
    else signal.addEventListener('abort', () => resolve(), { once: true })
  })

const routes = (router: ConnectRouter) => {
  const impl: ServiceImpl<typeof APIService> = {
    getSystemInfo: () => fixtures.systemInfo,
    generateKeyPair: () =>
      create(GenerateKeyPairResponseSchema, { keyId: 'generated', publicKey: 'ssh-ed25519 AAAAgenerated' }),
    getMe: () => fixtures.users[0],
    getUsers: () => create(GetUsersResponseSchema, { users: fixtures.users }),
    createUserKey: (req) => {
      const key = create(UserKeySchema, {
        id: `key-${fixtures.userKeys.length + 1}`,
        userId: ids.me,
        name: req.name,
        publicKey: req.publicKey,
        createdAt: ts('2024-06-01T00:00:00Z'),
      })
      fixtures.userKeys.push(key)
      return key
    },
    getUserKeys: () => create(GetUserKeysResponseSchema, { keys: fixtures.userKeys }),
    deleteUserKey: (req) => {
      fixtures.userKeys = fixtures.userKeys.filter((k) => k.id !== req.keyId)
      return {}
    },

    getRepositories: (_req) => {
      return create(GetRepositoriesResponseSchema, { repositories: fixtures.repositories })
    },
    getRepositoryCommits: (req) => {
      return create(GetRepositoryCommitsResponseSchema, {
        commits: fixtures.commits.filter((c) => req.hashes.includes(c.hash)),
      })
    },
    getRepository: (req) => findOrThrow(fixtures.repositories, req.repositoryId, 'repository'),
    getRepositoryRefs: () =>
      create(GetRepositoryRefsResponseSchema, {
        refs: [
          create(GitRefSchema, { refName: 'main', commit: hashes.a }),
          create(GitRefSchema, { refName: 'develop', commit: hashes.c }),
        ],
      }),
    updateRepository: (req) => {
      const repo = findOrThrow(fixtures.repositories, req.id, 'repository')
      if (req.name !== undefined) repo.name = req.name
      if (req.url !== undefined) repo.url = req.url
      if (req.ownerIds) repo.ownerIds = req.ownerIds.ownerIds
      return {}
    },
    refreshRepository: () => ({}),
    deleteRepository: () => ({}),
    createRepository: (req) => {
      const repo = create(RepositorySchema, {
        id: 'repo-new',
        name: req.name,
        url: req.url,
        htmlUrl: req.url,
        ownerIds: [ids.me],
      })
      fixtures.repositories.push(repo)
      return repo
    },

    createApplication: (req) => {
      const app = clone(ApplicationSchema, findOrThrow(fixtures.applications, ids.appRuntime, 'application'))
      app.id = 'app-new'
      app.name = req.name
      app.repositoryId = req.repositoryId
      app.refName = req.refName
      app.ownerIds = [ids.me]
      fixtures.applications.push(app)
      return app
    },
    getApplications: (req) => {
      const apps =
        req.scope === GetApplicationsRequest_Scope.REPOSITORY
          ? fixtures.applications.filter((a) => a.repositoryId === req.repositoryId)
          : fixtures.applications
      return create(GetApplicationsResponseSchema, { applications: apps })
    },
    getApplication: (req) => findOrThrow(fixtures.applications, req.id, 'application'),
    updateApplication: (req) => {
      const app = findOrThrow(fixtures.applications, req.id, 'application')
      if (req.name !== undefined) app.name = req.name
      if (req.refName !== undefined) app.refName = req.refName
      if (req.config) app.config = req.config
      if (req.websites) {
        app.websites = req.websites.websites.map((w, i) =>
          create(WebsiteSchema, { ...w, $typeName: 'neoshowcase.protobuf.Website', id: `website-new-${i}` }),
        )
      }
      if (req.portPublications) app.portPublications = req.portPublications.portPublications
      if (req.ownerIds) app.ownerIds = req.ownerIds.ownerIds
      return {}
    },
    deleteApplication: () => ({}),
    startApplication: (req) => {
      findOrThrow(fixtures.applications, req.id, 'application').running = true
      return {}
    },
    stopApplication: (req) => {
      findOrThrow(fixtures.applications, req.id, 'application').running = false
      return {}
    },

    getAvailableMetrics: () => create(AvailableMetricsSchema, { metricsNames: ['cpu', 'memory'] }),
    getApplicationMetrics: (_req) => {
      const now = Math.floor(Date.now() / 1000)
      return create(ApplicationMetricsSchema, {
        metrics: Array.from({ length: 5 }, (_, i) =>
          create(ApplicationMetricSchema, {
            time: create(TimestampSchema, { seconds: BigInt(now - (5 - i) * 60) }),
            value: i * 10,
          }),
        ),
      })
    },
    // one line a minute older than `before` per call: 'previous log line', then 'older log line 1', 2, ...
    getOutput: (req) => {
      const n = requestsTo('getOutput').indexOf(req as unknown as Record<string, unknown>) // by arrival, even if held
      const before = req.before?.seconds ?? BigInt(Math.floor(Date.now() / 1000))
      return create(ApplicationOutputsSchema, {
        outputs: [
          create(ApplicationOutputSchema, {
            time: create(TimestampSchema, { seconds: before - 60n }),
            log: n === 0 ? 'previous log line' : `older log line ${n}`,
          }),
        ],
      })
    },
    // yields a couple of lines, then keeps the stream open until the client aborts
    getOutputStream: async function* (_req, ctx) {
      yield create(ApplicationOutputSchema, { time: ts('2024-01-02T00:00:01Z'), log: 'streamed line 1' })
      yield create(ApplicationOutputSchema, { time: ts('2024-01-02T00:00:02Z'), log: 'streamed line 2' })
      if (outputStreamDrops > 0) {
        outputStreamDrops--
        throw new ConnectError('stream broke', Code.Unavailable)
      }
      await waitUntilAborted(ctx.signal)
    },

    getEnvVars: (req) =>
      create(ApplicationEnvVarsSchema, {
        variables: [
          ...[...(fixtures.envVars.get(req.id) ?? [])].map(([key, value]) =>
            create(ApplicationEnvVarSchema, { applicationId: req.id, key, value, system: false }),
          ),
          create(ApplicationEnvVarSchema, { applicationId: req.id, key: 'NS_APP_ID', value: req.id, system: true }),
        ],
      }),
    setEnvVar: (req) => {
      const vars = fixtures.envVars.get(req.applicationId) ?? new Map()
      vars.set(req.key, req.value)
      fixtures.envVars.set(req.applicationId, vars)
      return {}
    },
    deleteEnvVar: (req) => {
      fixtures.envVars.get(req.applicationId)?.delete(req.key)
      return {}
    },

    getAllBuilds: () => create(GetBuildsResponseSchema, { builds: fixtures.builds }),
    getBuilds: (req) =>
      create(GetBuildsResponseSchema, { builds: fixtures.builds.filter((b) => b.applicationId === req.id) }),
    getBuild: (req) => findOrThrow(fixtures.builds, req.buildId, 'build'),
    retryCommitBuild: () => ({}),
    cancelBuild: () => ({}),
    getBuildArtifact: () => create(ArtifactContentSchema, { filename: 'artifact.tar', content: new Uint8Array() }),
    getBuildLog: () => create(BuildLogSchema, { log: new TextEncoder().encode('step 1: build\nstep 2: done\n') }),
    getBuildLogStream: async function* (req, ctx) {
      yield create(BuildLogSchema, { log: new TextEncoder().encode(`streamed build line 1 of ${req.buildId}\n`) })
      await waitUntilAborted(ctx.signal)
    },
  }

  // record every request (also the streaming ones) before dispatching to the implementation
  const recorded = Object.fromEntries(
    Object.entries(impl).map(([method, fn]) => [
      method,
      (req: object, ctx: unknown) => {
        requests.push({ method, message: req as Record<string, unknown> })
        const remaining = failing.get(method)
        if (remaining !== undefined) {
          if (remaining <= 1) failing.delete(method)
          else failing.set(method, remaining - 1)
          throw new ConnectError(
            failureMessages.get(method) ?? 'backend is down',
            failureCodes.get(method) ?? Code.Unavailable,
          )
        }
        const call = () => (fn as (req: object, ctx: unknown) => unknown)(req, ctx)
        const hold = held.get(method)
        if (!hold) return call()
        if (hold.once) held.delete(method)
        return hold.gate.then(call)
      },
    ]),
  )
  router.service(APIService, recorded as ServiceImpl<typeof APIService>)
}

export const mockTransport = createRouterTransport(routes)
