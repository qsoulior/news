<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue"
import { useRoute, useRouter, type LocationQuery } from "vue-router"
import { NFlex, NCollapseTransition, NButton, NIcon, NDivider, useMessage } from "naive-ui"
import type { NewsHead } from "@/entities/news"
import { IconFilter, IconFilterDismiss } from "@/components/icons"
import ListSort from "@/components/ListSort.vue"
import ListFilter from "@/components/ListFilter.vue"
import ListSearch from "@/components/ListSearch.vue"
import ListContent from "@/components/ListContent.vue"
import { getNewsHead, toDateString } from "@/services/news"
import { getQueryStr, getQueryStrs, getQueryInt } from "@/router/query"

const LIMIT = 20
const route = useRoute()
const message = useMessage()

interface Props {
  page?: number
}

const props = withDefaults(defineProps<Props>(), {
  page: 1
})

const fromDateQuery = (query: string | null) => (query != null ? new Date(query).getTime() : null)

function initParams() {
  const query = route.query

  searchText.value = getQueryStr(query, "q") ?? ""
  filter.tags = getQueryStrs(query, "tags[]")
  filter.sources = getQueryStrs(query, "sources[]")
  filter.dateStart = fromDateQuery(getQueryStr(query, "date_from"))
  filter.dateEnd = fromDateQuery(getQueryStr(query, "date_to"))

  initSort()
}

function initSort() {
  const sortQuery = getQueryInt(route.query, "sort")
  if (sortQuery != null) {
    const sortValue = sortMap.get(sortQuery)
    if (sortValue != undefined) {
      sort.type = sortValue.type
      sort.ascending = sortValue.ascending
      localStorage.setItem("sort", JSON.stringify(sort))
      return
    }
  }

  const sortItem = localStorage.getItem("sort")
  if (sortItem == null) return

  const { type = "date", ascending = false } = JSON.parse(sortItem)
  if (type == "date" && !ascending) return

  sort.type = type
  sort.ascending = ascending
  replaceRouteSort({ ...route.query }, { type, ascending })
}

const router = useRouter()
async function onUpdatePage(page: number) {
  const query = { ...route.query }
  query["page"] = page == 1 ? [] : page.toString()
  router.push({ query: query })
}

// search
const searchText = ref<string>("")
async function onSubmitSearch(text: string) {
  const query = { ...route.query }
  delete query["page"]

  query["q"] = text != "" ? text : []

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

const toDateQuery = (timestamp: number | null) => (timestamp != null ? toDateString(new Date(timestamp)) : [])

async function onSubmitFilter(filter: Filter) {
  const query = { ...route.query }
  delete query["page"]

  query["sources[]"] = filter.sources
  query["tags[]"] = filter.tags
  query["date_from"] = toDateQuery(filter.dateStart)
  query["date_to"] = toDateQuery(filter.dateEnd)

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

function getSortOption(sort: Sort): SortOption | undefined {
  for (const [opt, val] of sortMap) {
    if (val.type == sort.type && val.ascending == sort.ascending) {
      return opt
    }
  }
  return undefined
}

async function onUpdateSort(sort: Sort) {
  localStorage.setItem("sort", JSON.stringify(sort))
  const query = { ...route.query }
  delete query["page"]

  replaceRouteSort(query, sort)
}

function replaceRouteSort(query: LocationQuery, sort: Sort) {
  const opt = getSortOption(sort)
  if (opt != undefined) {
    query["sort"] = opt.toString()
  }

  router.replace({ query: query })
}

const news = ref<NewsHead[]>([])
const count = ref(0)
const pageCount = computed(() => Math.ceil(count.value / LIMIT))

const loading = ref(false)
async function getNews(page: number) {
  if (page < 1) {
    news.value = []
    count.value = 0
    return
  }

  const skip = (page - 1) * LIMIT
  loading.value = true

  try {
    const { results, totalCount } = await getNewsHead(
      {
        text: searchText.value,
        sources: filter.sources,
        tags: filter.tags,
        dateFrom: filter.dateStart != null ? new Date(filter.dateStart) : undefined,
        dateTo: filter.dateEnd != null ? new Date(filter.dateEnd) : undefined
      },
      {
        limit: LIMIT,
        skip: skip,
        sort: getSortOption(sort)
      }
    )

    news.value = results
    count.value = totalCount
  } catch (err) {
    if (err instanceof Error) {
      console.error(err)
      message.error("Ошибка получения новостей")
    }
  } finally {
    loading.value = false
  }
}

watch(
  () => route.query,
  () => {
    initParams()
    return getNews(props.page)
  },
  { immediate: true }
)

watch(sort, onUpdateSort)
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <ListSearch v-model:value="searchText" @submit="onSubmitSearch" />
    <n-collapse-transition :show="isFilterShown">
      <ListFilter v-model:value="filter" @submit="onSubmitFilter" />
    </n-collapse-transition>
    <n-flex align="center" justify="space-between">
      <ListSort v-model:value="sort" @update:value="onUpdateSort" />
      <n-flex align="center">
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
    <n-divider v-if="count > 0 && !loading" class="list__divider_text"> Новостей: {{ count }} </n-divider>
    <n-divider v-else class="list__divider_empty" />
    <ListContent :news="news" :loading="loading" :page="page" :page-count="pageCount" @update:page="onUpdatePage" />
  </n-flex>
</template>

<style scoped>
.list__divider_text {
  margin: 0;
  font-size: inherit;
}

.list__divider_empty {
  --dm: calc((1.4em - 1px) / 2);
  margin: var(--dm) 0;
}
</style>
