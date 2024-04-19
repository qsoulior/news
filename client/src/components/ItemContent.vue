<script setup lang="ts">
import { computed } from "vue"
import { NFlex, NText, NH2, NTag, NDivider, NIcon, NImage, NEllipsis } from "naive-ui"
import { IconCalendar, IconPerson } from "@/components/icons"
import ItemContentControl from "@/components/ItemContentControl.vue"
import { getSourceImg, getSourceName } from "@/services/news"

const props = defineProps<{
  title: string
  description: string
  link: string
  source: string
  publishedAt?: Date
  authors: string[]
  tags: string[]
  categories: string[]
  content: string
}>()

const sourceImg = computed(() => getSourceImg(props.source))
const sourceName = computed(() => getSourceName(props.source))
const contents = computed(() => props.content.split("\n"))
const authorsStr = computed(() => props.authors.join(", "))
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <n-flex align="center" justify="space-between">
      <n-flex size="small" align="center">
        <n-image :src="sourceImg" width="18" preview-disabled :alt="source" />
        <n-text>{{ sourceName }}</n-text>
      </n-flex>
      <n-flex v-if="publishedAt" size="small" align="center">
        <n-icon :size="20">
          <IconCalendar />
        </n-icon>
        <n-text>{{ publishedAt.toLocaleString() }}</n-text>
      </n-flex>
    </n-flex>
    <n-flex :wrap="false" justify="space-between">
      <n-flex vertical size="large">
        <n-flex vertical>
          <n-h2 style="margin: 0">{{ title }}</n-h2>
          <n-text depth="3">{{ description }}</n-text>
        </n-flex>
      </n-flex>
      <ItemContentControl :link="link" />
    </n-flex>
    <n-flex>
      <n-tag class="item-tag" v-for="(tag, i) in tags" :key="i">
        <n-ellipsis :tooltip="{ delay: 500, contentClass: 'item-tag__tooltip' }">
          {{ tag }}
        </n-ellipsis>
      </n-tag>
    </n-flex>
    <n-divider style="margin: 1em 0" />
    <n-flex vertical :size="20">
      <n-text v-for="(text, i) in contents" :key="i">{{ text }}</n-text>
    </n-flex>
    <n-divider style="margin: 1em 0" />
    <n-flex align="center" justify="space-between">
      <n-flex v-if="authors.length > 0" size="small" align="center">
        <n-icon :size="20">
          <IconPerson />
        </n-icon>
        <n-text>{{ authorsStr }}</n-text>
      </n-flex>
      <n-flex size="small" align="center">
        <n-tag v-for="(category, i) in categories" :key="i" :bordered="false">{{ category }}</n-tag>
      </n-flex>
    </n-flex>
  </n-flex>
</template>

<style scoped>
.item-tag {
  min-width: 0;
}

.item-tag :deep(.n-tag__content) {
  overflow: hidden;
}

:global(.item-tag__tooltip) {
  max-width: calc(100vw - 4rem);
}
</style>
