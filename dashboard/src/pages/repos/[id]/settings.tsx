import { type RouteSectionProps, useParams } from '@solidjs/router'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { SideView } from '/@/components/layouts/SideView'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { SideLink } from '/@/components/UI/SideLink'

export default (props: RouteSectionProps) => {
  const params = useParams()
  const base = () => `/repos/${params.id}/settings`

  return (
    <MainViewContainer>
      <SideView.Container>
        <SideView.Side>
          <div class="sticky top-0 flex w-full flex-col">
            <SideLink
              href={base()}
              end
              icon={<div class="i-material-symbols:browse-activity-outline shrink-0 text-2xl/6" />}
            >
              General
            </SideLink>
            <SideLink
              href={`${base()}/authorization`}
              icon={<div class="i-material-symbols:conversion-path shrink-0 text-2xl/6" />}
            >
              Authorization
            </SideLink>
            <SideLink
              href={`${base()}/owners`}
              icon={<div class="i-material-symbols:person-outline shrink-0 text-2xl/6" />}
            >
              Owners
            </SideLink>
          </div>
        </SideView.Side>
        <SideView.Main>
          <PageBoundary fallback={<SettingSkeleton />}>{props.children}</PageBoundary>
        </SideView.Main>
      </SideView.Container>
    </MainViewContainer>
  )
}
