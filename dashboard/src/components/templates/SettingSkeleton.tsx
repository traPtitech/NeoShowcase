import type { Component } from 'solid-js'
import { Button } from '../UI/Button'
import Skeleton from '../UI/Skeleton'
import { DataTable } from '../layouts/DataTable'
import FormBox from '../layouts/FormBox'
import { FormItem } from './FormItem'

const SettingSkeleton: Component = () => {
  return (
    <DataTable.Container>
      <DataTable.Title>
        <Skeleton>Config Title</Skeleton>
      </DataTable.Title>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title={<Skeleton>Form Label</Skeleton>}>
            <Skeleton width={Number.NaN} height={48} />
          </FormItem>
          <FormItem title={<Skeleton>Second Form Label</Skeleton>}>
            <Skeleton width={Number.NaN} height={48} />
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <Skeleton>
            <Button size="small" variants="primary">
              save
            </Button>
          </Skeleton>
        </FormBox.Actions>
      </FormBox.Container>
    </DataTable.Container>
  )
}

export default SettingSkeleton
