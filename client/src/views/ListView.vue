<script setup lang="ts">
import { ref, watch } from "vue"
import { useRouter } from "vue-router"
import { NFlex, NCollapseTransition, NButton, NText, NDivider } from "naive-ui"
import type { NewsHead } from "@/entities/news"
import { IconFilter, IconFilterDismiss } from "@/components/icons"
import ListSort from "@/components/ListSort.vue"
import ListFilter from "@/components/ListFilter.vue"
import ListSearch from "@/components/ListSearch.vue"
import ListContent from "@/components/ListContent.vue"
import { getNewsHead } from "@/services/news"

interface Props {
  page?: number
}

const props = withDefaults(defineProps<Props>(), {
  page: 1
})

const router = useRouter()
function onUpdatePage(page: number) {
  router.push({ name: "list", params: { page: page == 1 ? "" : page.toString() } })
}

const isFilterShown = ref(false)
const count = ref(10)
const countTotal = ref(1000)

const news = ref<NewsHead[]>([])
const loading = ref(false)

let timer = 0
function loadNews(page: number) {
  clearTimeout(timer)
  loading.value = true
  timer = setTimeout(async () => {
    news.value = await getNewsHead(count.value, 0)
    loading.value = false
  }, 100)
}

watch(
  () => props.page,
  (page) => loadNews(page),
  { immediate: true }
)
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <ListSearch />
    <n-collapse-transition :show="isFilterShown">
      <ListFilter />
    </n-collapse-transition>
    <n-flex align="center" justify="space-between">
      <ListSort />
      <n-flex align="center">
        <n-text v-if="!loading">Результатов: {{ countTotal }}</n-text>
        <n-button tertiary @click="isFilterShown = !isFilterShown">
          <template #icon>
            <IconFilterDismiss v-if="isFilterShown" />
            <IconFilter v-else />
          </template>
        </n-button>
      </n-flex>
    </n-flex>
    <n-divider style="margin: 0" />
    <ListContent
      :news="news"
      :loading="loading"
      :page="page"
      :page-count="countTotal / count"
      @update:page="onUpdatePage"
    />
  </n-flex>
</template>
