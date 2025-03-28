import { match } from 'ts-pattern'
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
  v.transform(
    (input): AuthenticationType =>
      match(input)
        .returnType<AuthenticationType>()
        .with(`${AuthenticationType.OFF}`, () => AuthenticationType.OFF)
        .with(`${AuthenticationType.SOFT}`, () => AuthenticationType.SOFT)
        .with(`${AuthenticationType.HARD}`, () => AuthenticationType.HARD)
        .exhaustive(),
  ),
)

export const createWebsiteSchema = v.pipe(
  v.object({
    subdomain: v.optional(v.string()),
    domain: v.string(),
    pathPrefix: v.string(),
    stripPrefix: v.boolean(),
    https: stringBooleanSchema,
    // Static App ではフォームに表示されないので undefined になってしまう
    // そのままだと invalid になって設定を変更できないのでデフォルト値を設定する
    h2c: v.optional(v.boolean(), false),
    httpPort: v.optional(v.pipe(v.number(), v.integer()), 80),
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
  v.transform((input): CreateWebsiteRequest => {
    // wildcard domainならsubdomainとdomainを結合
    const fqdn = input.domain.startsWith('*')
      ? `${input.subdomain}${input.domain.replace(/\*/g, '')}`
      : // non-wildcard domainならdomainをそのまま使う
        input.domain

    return {
      $typeName: 'neoshowcase.protobuf.CreateWebsiteRequest',
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

export const getInitialValueOfCreateWebsiteForm = (domain: AvailableDomain): CreateWebsiteInput => ({
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
