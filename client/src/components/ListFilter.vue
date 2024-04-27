<script setup lang="ts">
import { NFlex, NButton, NSelect, NText, NDatePicker, NIcon, NButtonGroup, type SelectOption } from "naive-ui"
import { IconPlus, IconReset, IconSubmit } from "@/components/icons"
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
        class="filter__sources"
        v-model:value="filter.sources"
        :options="sourceOptions"
        multiple
        max-tag-count="responsive"
        placeholder="Источники"
      />
      <n-flex class="filter__dates" align="center" justify="space-between" :wrap="false">
        <n-date-picker
          class="filter__date"
          v-model:value="filter.dateStart"
          clearable
          type="date"
          placeholder="дд.мм.гггг"
          format="dd.MM.yyyy"
          :actions="null"
          :shortcuts="dateShortcuts"
        />
        <n-text>/</n-text>
        <n-date-picker
          class="filter__date"
          v-model:value="filter.dateEnd"
          clearable
          type="date"
          placeholder="дд.мм.гггг"
          format="dd.MM.yyyy"
          :actions="null"
          :shortcuts="dateShortcuts"
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
      <n-button-group>
        <n-button @click="reset" title="Сбросить фильтры">
          <template #icon>
            <n-icon>
              <IconReset />
            </n-icon>
          </template>
        </n-button>
        <n-button @click="submit" title="Применить фильтры">
          <template #icon>
            <n-icon>
              <IconSubmit />
            </n-icon>
          </template>
        </n-button>
      </n-button-group>
    </n-flex>
  </n-flex>
</template>

<style scoped>
.filter__sources {
  width: 20em;
  flex-grow: 1;
}

.filter__dates {
  flex-grow: 1;
}

.filter__date {
  width: 10em;
  flex-grow: 1;
}

@media screen and (min-width: 768px) {
  .filter__sources {
    flex-grow: 0;
  }

  .filter__dates {
    flex-grow: 0;
  }
}
</style>
