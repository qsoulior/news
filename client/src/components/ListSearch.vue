<script setup lang="ts">
import { NInput, NInputGroup, NButton, NIcon } from "naive-ui"
import { IconSearch } from "@/components/icons"
</script>

<template>
  <n-input-group>
    <n-input placeholder="Поиск" />
    <n-button>
      <template #icon>
        <n-icon>
          <IconSearch />
        </n-icon>
      </template>
    </n-button>
  </n-input-group>
</template>
