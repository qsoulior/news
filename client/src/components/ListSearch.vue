<script setup lang="ts">
import { NInput, NInputGroup, NButton, NIcon } from "naive-ui"
import { IconSearch } from "@/components/icons"

const value = defineModel<string>("value")
</script>

<template>
  <n-input-group>
    <n-input v-model:value="value" placeholder="Поиск" />
    <n-button>
      <template #icon>
        <n-icon>
          <IconSearch />
        </n-icon>
      </template>
    </n-button>
  </n-input-group>
</template>
