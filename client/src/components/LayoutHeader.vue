<script setup lang="ts">
import { NLayoutHeader, NFlex, NButton, NIcon, NText } from "naive-ui"
import { IconGithub, IconSun, IconMoon } from "@/components/icons"
import { computed } from "vue"

type Theme = "light" | "dark"
const theme = defineModel<Theme>("theme")
const isThemeDark = computed(() => theme.value == "dark")
function setThemeLight() {
  theme.value = "light"
}
function setThemeDark() {
  theme.value = "dark"
}
</script>

<template>
  <n-layout-header bordered>
    <n-flex align="center" justify="space-between">
      <router-link :to="{ name: 'list' }" style="text-decoration: none">
        <n-text>Новости</n-text>
      </router-link>
      <n-flex size="small">
        <n-button v-if="isThemeDark" class="n-button--square" quaternary circle @click="setThemeLight">
          <template #icon>
            <n-icon>
              <IconSun />
            </n-icon>
          </template>
        </n-button>
        <n-button v-else class="n-button--square" quaternary circle @click="setThemeDark">
          <template #icon>
            <n-icon>
              <IconMoon />
            </n-icon>
          </template>
        </n-button>
        <n-button
          class="n-button--square"
          quaternary
          circle
          tag="a"
          href="https://github.com/qsoulior/news"
          target="_blank"
        >
          <template #icon>
            <n-icon>
              <IconGithub />
            </n-icon>
          </template>
        </n-button>
      </n-flex>
    </n-flex>
  </n-layout-header>
</template>

<style scoped>
.n-layout-header {
  padding: 0.5rem 2rem;
}

.n-button--square {
  border-radius: 3px;
}
</style>
