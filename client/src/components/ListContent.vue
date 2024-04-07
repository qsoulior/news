<script setup lang="ts">
import { ref } from "vue"
import { NFlex, NEmpty, NPagination } from "naive-ui"
import ListContentItem from "@/components/ListContentItem.vue"

const count = ref(10)
</script>

<template>
  <n-flex vertical>
    <n-empty v-if="count == 0" description="Новости не найдены" />
    <n-flex v-else vertical size="large" align="center">
      <n-flex vertical style="width: 100%">
        <ListContentItem
          v-for="i in count"
          :key="i"
          title="Заголовок"
          description="Описание"
          :published-at="new Date()"
          source="РИА Новости"
        />
      </n-flex>
      <n-pagination :page-count="100" />
    </n-flex>
  </n-flex>
</template>
