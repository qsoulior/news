<script setup lang="ts">
import { NFlex, NEmpty, NPagination } from "naive-ui"
import ListContentItem from "@/components/ListContentItem.vue"
import ListContentSkeleton from "@/components/ListContentSkeleton.vue"
import { type NewsHead } from "@/entities/news"

defineProps<{
  news: NewsHead[]
  loading: boolean
}>()

const page = defineModel<number>("page")
</script>

<template>
  <n-flex vertical>
    <n-empty v-if="!loading && news.length == 0" description="Новости не найдены" />
    <n-flex v-else vertical size="large" align="center">
      <n-flex v-if="loading" vertical style="width: 100%">
        <ListContentSkeleton v-for="i in 10" :key="i" />
      </n-flex>
      <n-flex v-else vertical style="width: 100%">
        <ListContentItem
          v-for="item in news"
          :key="item.id"
          :title="item.title"
          :description="item.description"
          :source="item.source"
          :published-at="item.publishedAt"
        />
      </n-flex>
      <n-pagination v-model:page="page" :page-count="100" />
    </n-flex>
  </n-flex>
</template>
