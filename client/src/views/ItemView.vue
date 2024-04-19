<script setup lang="ts">
import { onMounted, ref } from "vue"
import ItemContent from "@/components/ItemContent.vue"
import ItemEmpty from "@/components/ItemEmpty.vue"
import ItemSkeleton from "@/components/ItemSkeleton.vue"
import type { News } from "@/entities/news"
import { getNews } from "@/services/news"

interface Props {
  id: string
}

const props = defineProps<Props>()

const news = ref<News>()
const loading = ref(false)

async function loadNews(id: string) {
  loading.value = true
  news.value = await getNews(id)
  loading.value = false
}

onMounted(() => {
  loadNews(props.id)
})
</script>

<template>
  <ItemSkeleton v-if="loading" />
  <ItemContent
    v-else-if="news"
    :title="news.title"
    :description="news.description"
    :link="news.link"
    :source="news.source"
    :publishedAt="news.publishedAt"
    :authors="news.authors"
    :tags="news.tags"
    :categories="news.categories"
    :content="news.content"
  />
  <ItemEmpty v-else />
</template>
