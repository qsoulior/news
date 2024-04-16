<script setup lang="ts">
import { reactive } from "vue"
import { NFlex, NButton, NSelect, NIcon, type SelectOption } from "naive-ui"
import { IconArrow } from "@/components/icons"

interface Sort {
  type: "relevance" | "date"
  ascending: boolean
}

const sort = defineModel<Sort>("value", {
  default: () =>
    reactive({
      type: "relevance",
      ascending: true
    })
})

const options: SelectOption[] = [
  {
    label: "По релевантности",
    value: "relevance"
  },
  {
    label: "По дате",
    value: "date"
  }
]
</script>

<template>
  <n-flex :wrap="false">
    <n-select v-model:value="sort.type" :options="options" placeholder="Сортировка" style="width: 12em" />
    <n-button tertiary title="Изменить направление сортировки" @click="sort.ascending = !sort.ascending">
      <template #icon>
        <n-icon :style="sort.ascending ? '' : 'transform: rotate(180deg)'">
          <IconArrow />
        </n-icon>
      </template>
    </n-button>
  </n-flex>
</template>
