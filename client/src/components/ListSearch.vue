<script setup lang="ts">
import { NInput, NInputGroup, NButton, NIcon } from "naive-ui"
import { IconSearch } from "@/components/icons"

const search = defineModel<string>("value", {
  default: ""
})
const emit = defineEmits<{
  submit: [value: string]
}>()

function submit() {
  emit("submit", search.value)
}
</script>

<template>
  <n-input-group>
    <n-input v-model:value="search" placeholder="Поиск" @keyup.enter="submit" />
    <n-button @click="submit">
      <template #icon>
        <n-icon>
          <IconSearch />
        </n-icon>
      </template>
    </n-button>
  </n-input-group>
</template>
