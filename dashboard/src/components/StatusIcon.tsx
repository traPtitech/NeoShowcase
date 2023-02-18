import { JSXElement } from 'solid-js'
import { AiFillCheckCircle } from 'solid-icons/ai'
import { vars } from '/@/theme.css'
import { IoReloadCircle } from 'solid-icons/io'
import { BiSolidErrorCircle } from 'solid-icons/bi'

export type Status = 'running' | 'static' | 'deploying' | 'error'

interface Props {
  status: Status
}

export const StatusIcon = ({ status }: Props): JSXElement => {
  switch (status) {
    case 'running':
      return <AiFillCheckCircle size={20} color={vars.icon.success1} />
    case 'static':
      return <AiFillCheckCircle size={20} color={vars.icon.success2} />
    case 'deploying':
      return <IoReloadCircle size={20} color={vars.icon.pending} />
    case 'error':
      return <BiSolidErrorCircle size={20} color={vars.icon.error} />
  }
}
