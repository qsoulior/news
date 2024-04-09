<script setup lang="ts">
import { onMounted, ref } from "vue"
import ItemContent from "@/components/ItemContent.vue"
import ItemEmpty from "@/components/ItemEmpty.vue"
import type { News } from "@/entities/news"
import { getNews } from "@/services/news"

const news = ref<News>()
const loading = ref(false)

let timer = 0
function loadNews() {
  clearTimeout(timer)
  loading.value = true
  timer = setTimeout(async () => {
    news.value = await getNews("")
    loading.value = false
  }, 100)
}

onMounted(() => {
  loadNews()
})
</script>

<template>
  <ItemContent
    v-if="news"
    :title="news.title"
    :description="news.description"
    :source="news.source"
    :publishedAt="news.publishedAt"
    :authors="news.authors"
    :tags="news.tags"
    :categories="news.categories"
    :content="news.content"
  />
  <ItemEmpty v-else />
</template>
