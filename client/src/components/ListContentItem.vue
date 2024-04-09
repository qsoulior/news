<script setup lang="ts">
import { NCard, NFlex, NText } from "naive-ui"

defineProps<{
  id: string
  title: string
  description: string
  publishedAt: Date
  source: string
}>()
</script>

<template>
  <router-link :to="{ name: 'item', params: { id: id } }" style="text-decoration: none">
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
  </router-link>
</template>
