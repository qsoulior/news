<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import { NFlex, NCollapseTransition, NButton, NIcon, NText, NDivider } from "naive-ui"
import type { NewsHead } from "@/entities/news"
import { IconFilter, IconFilterDismiss } from "@/components/icons"
import ListSort from "@/components/ListSort.vue"
import ListFilter from "@/components/ListFilter.vue"
import ListSearch from "@/components/ListSearch.vue"
import ListContent from "@/components/ListContent.vue"
import { getNewsHead } from "@/services/news"
import { getQueryStr, getQueryStrs, getQueryInt } from "@/router/query"

const LIMIT = 10
const route = useRoute()

interface Props {
  page?: number
}

const props = withDefaults(defineProps<Props>(), {
  page: 1
})

function setLocalParams() {
  const sortItem = localStorage.getItem("sort")
  if (sortItem != null) {
    const { type = "date", ascending = false } = JSON.parse(sortItem)
    sort.type = type
    sort.ascending = ascending
  }
}

function setRouteParams() {
  const query = route.query

  searchValue.value = getQueryStr(query, "q") ?? ""
  filter.tags = getQueryStrs(query, "tags[]")
  filter.sources = getQueryStrs(query, "sources[]")
  filter.dateStart = getQueryInt(query, "date_start")
  filter.dateEnd = getQueryInt(query, "date_end")

  const sortQuery = getQueryInt(query, "sort")
  if (sortQuery != null) {
    const sortValue = sortMap.get(sortQuery)
    if (sortValue != undefined) {
      sort.type = sortValue.type
      sort.ascending = sortValue.ascending
    }
  }
}

const router = useRouter()
function onUpdatePage(page: number) {
  const query = { ...route.query }
  query["page"] = page == 1 ? [] : page.toString()
  router.push({ name: "list", query: query })
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
const searchValue = ref<string>("")
function onSubmitSearch(search: string) {
  // update query
  const query = { ...route.query }
  query["q"] = search != "" ? search : []

  router.replace({ query: query })
}

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
  // update query
  const query = { ...route.query }
  query["sources[]"] = filter.sources
  query["tags[]"] = filter.tags
  query["date_start"] = filter.dateStart?.toString() ?? []
  query["date_end"] = filter.dateEnd?.toString() ?? []

  router.replace({ query: query })
}

// sort
interface Sort {
  type: "relevance" | "date"
  ascending: boolean
}

const sort = reactive<Sort>({
  type: "date",
  ascending: false
})

enum SortOption {
  SortDateDesc,
  SortDateAsc,
  SortRelevanceDesc,
  SortRelevanceAsc
}

const sortMap = new Map<SortOption, Sort>([
  [SortOption.SortDateDesc, { type: "date", ascending: false }],
  [SortOption.SortDateAsc, { type: "date", ascending: true }],
  [SortOption.SortRelevanceDesc, { type: "relevance", ascending: false }],
  [SortOption.SortRelevanceAsc, { type: "relevance", ascending: true }]
])

watch(sort, (sort) => {
  localStorage.setItem("sort", JSON.stringify(sort))
  // update query
  const query = { ...route.query }

  for (const [opt, val] of sortMap) {
    if (val.type == sort.type && val.ascending == sort.ascending) {
      query["sort"] = opt.toString()
      break
    }
  }

  router.replace({ query: query })
})

onMounted(() => {
  setLocalParams()
  setRouteParams()
})
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <ListSearch v-model:value="searchValue" @submit="onSubmitSearch" />
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
