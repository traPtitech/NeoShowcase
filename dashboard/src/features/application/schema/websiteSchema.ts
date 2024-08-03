import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import {
  AuthenticationType,
  type AvailableDomain,
  type CreateWebsiteRequest,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { stringBooleanSchema } from '/@/libs/schemaUtil'

// KobalteのRadioGroupではstringしか扱えないためform内では文字列として持つ
const authenticationSchema = v.pipe(
  v.union([
    v.literal(`${AuthenticationType.OFF}`),
    v.literal(`${AuthenticationType.SOFT}`),
    v.literal(`${AuthenticationType.HARD}`),
  ]),
  v.transform((input): AuthenticationType => {
    switch (input) {
      case `${AuthenticationType.OFF}`: {
        return AuthenticationType.OFF
      }
      case `${AuthenticationType.SOFT}`: {
        return AuthenticationType.SOFT
      }
      case `${AuthenticationType.HARD}`: {
        return AuthenticationType.HARD
      }
      default: {
        const _unreachable: never = input
        throw new Error('unknown website AuthenticationType')
      }
    }
  }),
)

export const createWebsiteSchema = v.pipe(
  v.variant('state', [
    v.pipe(
      v.object({
        state: v.union([v.literal('noChange'), v.literal('readyToChange'), v.literal('added')]),
        subdomain: v.optional(v.string()),
        domain: v.string(),
        pathPrefix: v.string(),
        stripPrefix: v.boolean(),
        https: stringBooleanSchema,
        h2c: v.boolean(),
        httpPort: v.pipe(v.number(), v.integer()),
        authentication: authenticationSchema,
      }),
      // wildcard domainが選択されている場合サブドメインは空であってはならない
      v.forward(
        v.partialCheck(
          [['subdomain'], ['domain']],
          (input) => {
            if (input.domain?.startsWith('*')) return input.subdomain !== ''
            return true
          },
          'Please Enter Subdomain Name',
        ),
        ['subdomain'],
      ),
    ),
    v.object({
      // 削除するwebsite設定の中身はチェックしない
      state: v.literal('readyToDelete'),
    }),
  ]),
  v.transform((input): PartialMessage<CreateWebsiteRequest> => {
    // 削除するwebsite設定の中身はチェックしない
    if (input.state === 'readyToDelete') return {}

    // wildcard domainならsubdomainとdomainを結合
    const fqdn = input.domain.startsWith('*')
      ? `${input.subdomain}${input.domain.replace(/\*/g, '')}`
      : // non-wildcard domainならdomainをそのまま使う
      input.domain

    return {
      fqdn,
      authentication: input.authentication,
      h2c: input.h2c,
      httpPort: input.httpPort,
      https: input.https,
      // バックエンド側では `/${prefix}` で持っている, フォーム内部では'/'を除いて持つ
      pathPrefix: `/${input.pathPrefix}`,
      stripPrefix: input.stripPrefix,
    }
  }),
)

export type CreateWebsiteInput = v.InferInput<typeof createWebsiteSchema>

export const parseCreateWebsiteInput = (input: unknown) => {
  const result = v.parse(createWebsiteSchema, input)
  return result
}

const extractSubdomain = (
  fqdn: string,
  availableDomains: AvailableDomain[],
): {
  subdomain: string
  domain: string
} => {
  const nonWildcardDomains = availableDomains.filter((d) => !d.domain.startsWith('*'))
  const wildcardDomains = availableDomains.filter((d) => d.domain.startsWith('*'))

  const matchNonWildcardDomain = nonWildcardDomains.find((d) => fqdn === d.domain)
  if (matchNonWildcardDomain !== undefined) {
    return {
      subdomain: '',
      domain: matchNonWildcardDomain.domain,
    }
  }

  const matchDomain = wildcardDomains.find((d) => fqdn.endsWith(d.domain.replace(/\*/g, '')))
  if (matchDomain === undefined) {
    const fallbackDomain = availableDomains.at(0)
    if (fallbackDomain === undefined) throw new Error('No domain available')
    return {
      subdomain: '',
      domain: fallbackDomain.domain,
    }
  }
  return {
    subdomain: fqdn.slice(0, -matchDomain.domain.length + 1),
    domain: matchDomain.domain,
  }
}

export const createWebsiteInitialValues = (domain: AvailableDomain): CreateWebsiteInput => ({
  state: 'added',
  domain: domain.domain,
  subdomain: '',
  pathPrefix: '',
  stripPrefix: false,
  https: 'true',
  h2c: false,
  httpPort: 80,
  authentication: `${AuthenticationType.OFF}`,
})

export const websiteMessageToSchema = (website: Website): CreateWebsiteInput => {
  const availableDomains = systemInfo()?.domains ?? []

  const { domain, subdomain } = extractSubdomain(website.fqdn, availableDomains)

  return {
    state: 'noChange',
    domain,
    subdomain,
    // バックエンド側では `/${prefix}` で持っている, フォーム内部では'/'を除いて持つ
    pathPrefix: website.pathPrefix.slice(1),
    stripPrefix: website.stripPrefix,
    https: website.https ? 'true' : 'false',
    h2c: website.h2c,
    httpPort: website.httpPort,
    authentication: `${website.authentication}`,
  }
}
