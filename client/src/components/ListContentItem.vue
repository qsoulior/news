<script setup lang="ts">
import { NCard, NFlex, NText } from "naive-ui"

const props = defineProps<{
  title: string
  description: string
  publishedAt: Date
  source: string
}>()
</script>

<template>
  <n-card size="small">
    <n-flex vertical size="small">
      <n-text strong>{{ title }}</n-text>
      <n-text>{{ description }}</n-text>
      <n-flex justify="space-between">
        <n-text depth="3">{{ publishedAt.toLocaleString() }}</n-text>
        <n-text depth="3">{{ source }}</n-text>
      </n-flex>
    </n-flex>
  </n-card>
</template>
