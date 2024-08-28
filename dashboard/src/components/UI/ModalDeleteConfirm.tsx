import type { ParentComponent } from 'solid-js'

const ModalDeleteConfirm: ParentComponent = (props) => {
  return (
    <div class="overflow-wrap-anywhere h3-regular flex w-full items-center gap-2 rounded-lg bg-ui-secondary px-5 py-4 text-text-black">
      {props.children}
    </div>
  )
}

export default ModalDeleteConfirm
