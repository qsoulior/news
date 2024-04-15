<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue"
import { useRouter } from "vue-router"
import { NFlex, NCollapseTransition, NButton, NIcon, NText, NDivider } from "naive-ui"
import type { NewsHead } from "@/entities/news"
import { IconFilter, IconFilterDismiss } from "@/components/icons"
import ListSort from "@/components/ListSort.vue"
import ListFilter from "@/components/ListFilter.vue"
import ListSearch from "@/components/ListSearch.vue"
import ListContent from "@/components/ListContent.vue"
import { getNewsHead } from "@/services/news"

const LIMIT = 10

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

const news = ref<NewsHead[]>([])
const count = ref(1000)
const loading = ref(false)

let timer = 0
function loadNews(page: number) {
  const skip = (page - 1) * LIMIT

  clearTimeout(timer)
  loading.value = true
  timer = setTimeout(async () => {
    news.value = await getNewsHead(LIMIT, skip)
    loading.value = false
  }, 100)
}

const pageCount = computed(() => Math.ceil(count.value / LIMIT))
watch(
  () => props.page,
  (page) => {
    if (page >= 1 && page <= pageCount.value) {
      loadNews(page)
    }
  },
  { immediate: true }
)

// search
const searchValue = ref<string>()

// filter
const isFilterShown = ref(false)

interface Filter {
  sources: string[]
  dateStart: number | null
  dateEnd: number | null
  tags: string[]
}

const filter = reactive<Filter>({
  sources: [],
  dateStart: null,
  dateEnd: null,
  tags: []
})

function onSubmitFilter(filter: Filter) {
  console.log(filter)
}

// sort
interface Sort {
  type: "relevance" | "date"
  ascending: boolean
}

const sort = reactive<Sort>({
  type: "relevance",
  ascending: true
})

watch(sort, (sort) => {
  localStorage.setItem("sort", JSON.stringify(sort))
})

onMounted(() => {
  const sortItem = localStorage.getItem("sort")
  if (sortItem != null) {
    const { type = "relevance", ascending = true } = JSON.parse(sortItem)
    sort.type = type
    sort.ascending = ascending
  }
})
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <ListSearch v-model:value="searchValue" />
    <n-collapse-transition :show="isFilterShown">
      <ListFilter v-model:value="filter" @submit="onSubmitFilter" />
    </n-collapse-transition>
    <n-flex align="center" justify="space-between">
      <ListSort v-model:value="sort" />
      <n-flex align="center">
        <n-text v-if="!loading">Результатов: {{ count }}</n-text>
        <n-button tertiary title="Показать фильтры" @click="isFilterShown = !isFilterShown">
          <template #icon>
            <n-icon>
              <IconFilterDismiss v-if="isFilterShown" />
              <IconFilter v-else />
            </n-icon>
          </template>
        </n-button>
      </n-flex>
    </n-flex>
    <n-divider style="margin: 0" />
    <ListContent :news="news" :loading="loading" :page="page" :page-count="pageCount" @update:page="onUpdatePage" />
  </n-flex>
</template>
