<script setup lang="ts">
import { NFlex, NButton, NSelect, NText, NDatePicker, type SelectOption } from "naive-ui"
import { IconPlus } from "@/components/icons"
import { reactive } from "vue"

interface Filter {
  sources: string[]
  dateStart: number | null
  dateEnd: number | null
  tags: string[]
}

const filter = defineModel<Filter>("value", {
  default: () =>
    reactive({
      sources: [],
      dateStart: null,
      dateEnd: null,
      tags: []
    })
})

const emit = defineEmits<{
  submit: [filter: Filter]
}>()

const sourceOptions: SelectOption[] = [
  {
    label: "РИА Новости",
    value: "ria"
  },
  {
    label: "Известия",
    value: "iz"
  },
  {
    label: "Лента.Ру",
    value: "lenta"
  },
  {
    label: "NewsData",
    value: "newsdata"
  }
]

const dateShortcuts: Record<string, () => number> = {
  Сегодня: () => new Date().getTime()
}

function submit() {
  emit("submit", filter.value)
}

function reset() {
  filter.value.sources = []
  filter.value.dateStart = null
  filter.value.dateEnd = null
  filter.value.tags = []
}
</script>

<template>
  <n-flex vertical>
    <n-flex justify="space-between">
      <n-select
        v-model:value="filter.sources"
        :options="sourceOptions"
        multiple
        max-tag-count="responsive"
        placeholder="Источники"
        style="max-width: 20em"
      />
      <n-flex align="center">
        <n-text>с</n-text>
        <n-date-picker
          v-model:value="filter.dateStart"
          clearable
          type="date"
          placeholder="дд.мм.гггг"
          format="dd.MM.yyyy"
          :actions="null"
          :shortcuts="dateShortcuts"
          style="max-width: 10em"
        />
        <n-text>по</n-text>
        <n-date-picker
          v-model:value="filter.dateEnd"
          clearable
          type="date"
          placeholder="дд.мм.гггг"
          format="dd.MM.yyyy"
          :actions="null"
          :shortcuts="dateShortcuts"
          style="max-width: 10em"
        />
      </n-flex>
    </n-flex>
    <n-flex justify="space-between" :wrap="false">
      <n-select
        v-model:value="filter.tags"
        filterable
        multiple
        clearable
        tag
        max-tag-count="responsive"
        placeholder="Тэги"
        :show="false"
        style="min-width: 0"
      >
        <template #arrow>
          <IconPlus />
        </template>
      </n-select>
      <n-flex :wrap="false">
        <n-button @click="reset">Сбросить</n-button>
        <n-button @click="submit">Применить</n-button>
      </n-flex>
    </n-flex>
  </n-flex>
</template>
