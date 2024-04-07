<script setup lang="ts">
import { ref } from "vue"
import { NFlex, NButton, NSelect, NText, NDatePicker, type SelectOption } from "naive-ui"
import { IconPlus } from "@/components/icons"

const sourceValue = ref<string | null>(null)
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

const dateShortcuts = {
  Сегодня: () => new Date().getTime()
}
</script>

<template>
  <n-flex vertical>
    <n-flex justify="space-between">
      <n-select
        v-model:value="sourceValue"
        :options="sourceOptions"
        multiple
        max-tag-count="responsive"
        placeholder="Источники"
        style="max-width: 20em"
      />
      <n-flex align="center">
        <n-text>с</n-text>
        <n-date-picker
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
        filterable
        multiple
        clearable
        tag
        max-tag-count="responsive"
        placeholder="Тэги"
        :show="false"
        style="max-width: 40em"
      >
        <template #arrow>
          <IconPlus />
        </template>
      </n-select>
      <n-flex :wrap="false">
        <n-button>Сбросить</n-button>
        <n-button>Применить</n-button>
      </n-flex>
    </n-flex>
  </n-flex>
</template>
