<script setup lang="ts">
import { NFlex, NButton, NSelect, NIcon, type SelectOption } from "naive-ui"
import { IconArrow } from "@/components/icons"

type Type = "relevance" | "date"
const type = defineModel<Type>("type", {
  default: "relevance"
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

const ascending = defineModel<boolean>("ascending", {
  default: true
})
</script>

<template>
  <n-flex :wrap="false">
    <n-select v-model:value="type" :options="options" placeholder="Сортировка" style="width: 12em" />
    <n-button tertiary title="Изменить направление сортировки" @click="ascending = !ascending">
      <template #icon>
        <n-icon :style="ascending ? 'transform: rotate(180deg)' : ''">
          <IconArrow />
        </n-icon>
      </template>
    </n-button>
  </n-flex>
</template>
