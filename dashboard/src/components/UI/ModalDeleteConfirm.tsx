import type { ParentComponent } from 'solid-js'

const ModalDeleteConfirm: ParentComponent = (props) => {
  return (
    <div class="overflow-wrap-anywhere h3-regular bg-ui-secondary text-text-black flex w-full items-center gap-2 rounded-lg px-5 py-4">
      {props.children}
    </div>
  )
}

export default ModalDeleteConfirm
