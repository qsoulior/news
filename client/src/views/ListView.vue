<script setup lang="ts">
import { ref } from "vue"
import { NFlex, NCollapseTransition, NButton, NText, NDivider } from "naive-ui"
import { IconFilter, IconFilterDismiss } from "@/components/icons"
import ListFilter from "@/components/ListFilter.vue"
import ListSearch from "@/components/ListSearch.vue"
import ListSort from "@/components/ListSort.vue"

const isFilterShown = ref(false)
</script>

<template>
  <n-flex vertical size="large">
    <ListSearch />
    <n-collapse-transition :show="isFilterShown">
      <ListFilter />
    </n-collapse-transition>
    <n-flex align="center" justify="space-between">
      <ListSort />
      <n-flex align="center">
        <n-text>Результатов: 1234</n-text>
        <n-button tertiary @click="isFilterShown = !isFilterShown">
          <template #icon>
            <IconFilterDismiss v-if="isFilterShown" />
            <IconFilter v-else />
          </template>
        </n-button>
      </n-flex>
    </n-flex>
    <n-divider style="margin: 0" />
  </n-flex>
</template>
