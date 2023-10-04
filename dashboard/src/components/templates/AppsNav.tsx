import AddIcon from '/@/assets/icons/24/add.svg'
import { Component } from 'solid-js'
import { Button } from '../UI/Button'
import { Nav } from './Nav'

export const AppsNav: Component = () => {
  return (
    <Nav
      title="Apps"
      action={
        <Button color="primary" size="medium" leftIcon={<AddIcon />}>
          Add New App
        </Button>
      }
    />
  )
}
