import { A } from '@solidjs/router'
import { Component } from 'solid-js'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { Nav } from './Nav'

export const AppsNav: Component = () => {
  return (
    <Nav
      title="Apps"
      action={
        <A href="/apps/new">
          <Button variants="primary" size="medium" leftIcon={<MaterialSymbols>add</MaterialSymbols>}>
            Add New App
          </Button>
        </A>
      }
    />
  )
}
