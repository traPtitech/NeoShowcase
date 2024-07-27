import { safeParse } from 'valibot'
import { describe, expect, test } from 'vitest'
import { AuthenticationType, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { createOrUpdateApplicationSchema } from './applicationSchema'

const validator = (input: unknown) => safeParse(createOrUpdateApplicationSchema, input)

describe('Create Application Schema', () => {
  const baseConfig = {
    deployConfig: {
      type: 'runtime',
      value: {
        runtime: {
          useMariadb: false,
          useMongodb: false,
          entrypoint: '.',
          command: "echo 'test'",
        },
      },
    },
    buildConfig: {
      type: 'buildpack',
      value: {
        buildpack: {
          context: '',
        },
      },
    },
  }

  const baseWebsite = {
    state: 'added',
    subdomain: 'example',
    domain: 'example.com',
    pathPrefix: '',
    stripPrefix: false,
    https: true,
    h2c: false,
    httpPort: 80,
    authentication: AuthenticationType.OFF,
  }

  const basePortPublication = {
    internetPort: 80,
    applicationPort: 3000,
    protocol: PortPublicationProtocol.TCP,
  }

  const base = {
    type: 'create',
    form: {
      name: 'test application',
      repositoryId: 'testRepoId',
      refName: 'main',
      config: baseConfig,
      websites: [baseWebsite],
      portPublications: [basePortPublication],
      startOnCreate: true,
    },
  }

  test('ok: valid input (runtime config)', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          config: {
            deployConfig: {
              type: 'runtime',
              value: {
                runtime: {
                  useMariadb: false,
                  useMongodb: false,
                  entrypoint: '.',
                  command: "echo 'test'",
                },
              },
            },
            buildConfig: {
              type: 'buildpack',
              value: {
                buildpack: {
                  context: '',
                },
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (static config)', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          config: {
            deployConfig: {
              type: 'static',
              value: {
                static: {
                  artifactPath: '.',
                  spa: 'false',
                },
              },
            },
            buildConfig: {
              type: 'buildpack',
              value: {
                buildpack: {
                  context: '',
                },
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (buildpack config)', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          config: {
            deployConfig: {
              type: 'runtime',
              value: {
                runtime: {
                  useMariadb: false,
                  useMongodb: false,
                  entrypoint: '.',
                  command: "echo 'test'",
                },
              },
            },
            buildConfig: {
              type: 'buildpack',
              value: {
                buildpack: {
                  context: '',
                },
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (dockerfile config)', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          config: {
            deployConfig: {
              type: 'runtime',
              value: {
                runtime: {
                  useMariadb: false,
                  useMongodb: false,
                  entrypoint: '.',
                  command: "echo 'test'",
                },
              },
            },
            buildConfig: {
              type: 'dockerfile',
              value: {
                dockerfile: {
                  dockerfileName: 'Dockerfile',
                  context: '',
                },
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (cmd config)', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          config: {
            deployConfig: {
              type: 'runtime',
              value: {
                runtime: {
                  useMariadb: false,
                  useMongodb: false,
                  entrypoint: '.',
                  command: "echo 'test'",
                },
              },
            },
            buildConfig: {
              type: 'cmd',
              value: {
                cmd: {
                  baseImage: 'node:22-alpine',
                  buildCmd: 'npm run build',
                },
              },
            },
          },
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty name', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          name: '',
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Application Name',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'name',
          }),
        ],
      }),
    ])
  })

  test('ng: empty refname', () => {
    expect(
      validator({
        ...base,
        form: {
          ...base.form,
          refName: '',
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Branch Name',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'refName',
          }),
        ],
      }),
    ])
  })

  // TODO: add ng config test
})

describe('Update Application Schema', () => {
  test('ok: valid input (update general config)', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: 'testAppId',
          name: 'test application',
          repositoryId: 'testRepoId',
          refName: 'main',
        },
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ng: empty name', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: 'testAppId',
          name: '',
          repositoryId: 'testRepoId',
          refName: 'main',
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Application Name',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'name',
          }),
        ],
      }),
    ])
  })

  test('ng: empty refname', () => {
    expect(
      validator({
        type: 'update',
        form: {
          id: 'testAppId',
          name: 'test application',
          repositoryId: 'testRepoId',
          refName: '',
        },
      }).issues,
    ).toEqual([
      expect.objectContaining({
        message: 'Enter Branch Name',
        path: [
          expect.objectContaining({
            key: 'form',
          }),
          expect.objectContaining({
            key: 'refName',
          }),
        ],
      }),
    ])
  })
})
