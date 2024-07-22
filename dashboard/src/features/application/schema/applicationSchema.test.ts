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
    name: 'test application',
    repositoryId: 'testRepoId',
    refName: 'main',
    config: baseConfig,
    websites: [baseWebsite],
    portPublications: [basePortPublication],
    startOnCreate: true,
  }

  test('ok: valid input (runtime config)', () => {
    expect(
      validator({
        ...base,
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
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (static config)', () => {
    expect(
      validator({
        ...base,
        config: {
          deployConfig: {
            type: 'static',
            value: {
              static: {
                artifactPath: '.',
                spa: false,
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
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (buildpack config)', () => {
    expect(
      validator({
        ...base,
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
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (dockerfile config)', () => {
    expect(
      validator({
        ...base,
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
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  test('ok: valid input (cmd config)', () => {
    expect(
      validator({
        ...base,
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
      }),
    ).toEqual(expect.objectContaining({ success: true }))
  })

  // test("ng: empty name", () => {
  // 	expect(
  // 		validator({
  // 			...base,
  // 			name: "",
  // 		}).issues,
  // 	).toEqual([
  // 		expect.objectContaining({
  // 			message: "Enter Repository Name",
  // 			path: [
  // 				expect.objectContaining({
  // 					key: "name",
  // 				}),
  // 			],
  // 		}),
  // 	]);
  // });

  // test("ng: empty url", () => {
  // 	expect(
  // 		validator({
  // 			...base,
  // 			url: "",
  // 		}).issues,
  // 	).toEqual([
  // 		expect.objectContaining({
  // 			message: "Enter Repository URL",
  // 			path: [
  // 				expect.objectContaining({
  // 					key: "url",
  // 				}),
  // 			],
  // 		}),
  // 	]);
  // });

  // test("ng: auth method is basic, but the URL starts with 'http://'", () => {
  // 	expect(
  // 		validator({
  // 			...base,
  // 			url: "http://example.com/test/test.git",
  // 			auth: {
  // 				method: "basic",
  // 				value: {
  // 					basic: {
  // 						username: "test name",
  // 						password: "test password",
  // 					},
  // 				},
  // 			},
  // 		}).issues,
  // 	).toEqual([
  // 		expect.objectContaining({
  // 			message:
  // 				"Basic認証を使用する場合、URLはhttps://から始まる必要があります",
  // 			path: [
  // 				expect.objectContaining({
  // 					key: "url",
  // 				}),
  // 			],
  // 		}),
  // 	]);
  // });
})

// describe("Update Repository Schema", () => {
// 	const base = {
// 		id: "testRepositoryId",
// 		type: "update",
// 		name: "test repository",
// 		url: "https://example.com/test/test.git",
// 		auth: {
// 			method: "none",
// 			value: {
// 				none: {},
// 			},
// 		},
// 		ownerIds: ["owner1"],
// 	};

// 	test("ok: valid input", () => {
// 		expect(validator(base)).toEqual(expect.objectContaining({ success: true }));
// 	});

// 	test("ok: update name", () => {
// 		expect(
// 			validator({
// 				id: base.id,
// 				type: base.type,
// 				name: base.name,
// 			}),
// 		).toEqual(expect.objectContaining({ success: true }));
// 	});

// 	test("ok: update auth config", () => {
// 		expect(
// 			validator({
// 				id: base.id,
// 				type: base.type,
// 				url: base.url,
// 				auth: base.auth,
// 			}),
// 		).toEqual(expect.objectContaining({ success: true }));
// 	});

// 	test("ok: update ownerIds", () => {
// 		expect(
// 			validator({
// 				id: base.id,
// 				type: base.type,
// 				ownerIds: base.ownerIds,
// 			}),
// 		).toEqual(expect.objectContaining({ success: true }));
// 	});

// 	test("ng: empty id", () => {
// 		expect(
// 			validator({
// 				...base,
// 				id: undefined,
// 			}).issues,
// 		).toEqual([
// 			expect.objectContaining({
// 				path: [
// 					expect.objectContaining({
// 						key: "id",
// 					}),
// 				],
// 			}),
// 		]);
// 	});

// 	test("ng: auth method is basic, but the URL starts with 'http://'", () => {
// 		expect(
// 			validator({
// 				...base,
// 				url: "http://example.com/test/test.git",
// 				auth: {
// 					method: "basic",
// 					value: {
// 						basic: {
// 							username: "test name",
// 							password: "test password",
// 						},
// 					},
// 				},
// 			}).issues,
// 		).toEqual([
// 			expect.objectContaining({
// 				message:
// 					"Basic認証を使用する場合、URLはhttps://から始まる必要があります",
// 				path: [
// 					expect.objectContaining({
// 						key: "url",
// 					}),
// 				],
// 			}),
// 		]);
// 	});
// });
