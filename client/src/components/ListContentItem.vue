<script setup lang="ts">
import { computed } from "vue"
import { NCard, NFlex, NText, NAvatar } from "naive-ui"
import { getSourceImg, getSourceName } from "@/services/news"

const props = defineProps<{
  id: string
  title: string
  description: string
  publishedAt: Date
  source: string
}>()

const sourceImg = computed(() => getSourceImg(props.source))
const sourceName = computed(() => getSourceName(props.source))
</script>

<template>
  <router-link :to="{ name: 'item', params: { id: id } }" style="text-decoration: none">
    <n-card size="small">
      <n-flex vertical size="small">
        <n-text strong>{{ title }}</n-text>
        <n-text>{{ description }}</n-text>
        <n-flex justify="space-between">
          <n-text depth="3">{{ publishedAt.toLocaleString() }}</n-text>
          <n-flex size="small" align="center">
            <n-avatar :src="sourceImg" :size="18" color="transparent" />
            <n-text depth="3">{{ sourceName }}</n-text>
          </n-flex>
        </n-flex>
      </n-flex>
    </n-card>
  </router-link>
</template>
