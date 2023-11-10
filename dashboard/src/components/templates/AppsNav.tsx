import { Component } from 'solid-js'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { Nav } from './Nav'

export const AppsNav: Component = () => {
  return (
    <Nav
      title="Apps"
      action={
        <Button variants="primary" size="medium" leftIcon={<MaterialSymbols>add</MaterialSymbols>}>
          Add New App
        </Button>
      }
    />
  )
}
