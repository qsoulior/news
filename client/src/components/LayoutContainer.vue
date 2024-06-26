<script setup lang="ts">
import { onMounted, ref, watch } from "vue"
import { RouterView } from "vue-router"
import { type GlobalTheme, NLayout, NLayoutContent, darkTheme, useOsTheme, NBackTop } from "naive-ui"
import LayoutHeader from "@/components/LayoutHeader.vue"
import LayoutFooter from "@/components/LayoutFooter.vue"

type ThemeType = "light" | "dark"
const themeType = ref<ThemeType>("light")
function getThemeObj(type: ThemeType) {
  return type == "dark" ? darkTheme : null
}

const themeObj = defineModel<GlobalTheme | null>("theme")
watch(themeType, (value) => {
  themeObj.value = getThemeObj(value)
  localStorage.setItem("theme", value)
})

function loadTheme() {
  const osThemeType = useOsTheme()
  const storageThemeType = localStorage.getItem("theme")
  switch (storageThemeType) {
    case "light":
      themeObj.value = null
      themeType.value = storageThemeType
      break
    case "dark":
      themeObj.value = darkTheme
      themeType.value = storageThemeType
      break
    default:
      themeObj.value = osThemeType.value == "dark" ? darkTheme : null
      themeType.value = osThemeType.value ?? "light"
  }
}

onMounted(() => {
  loadTheme()
})
</script>

<template>
  <n-layout style="height: 100vh">
    <n-back-top />
    <n-layout>
      <LayoutHeader v-model:theme="themeType" />
      <n-layout-content style="min-height: calc(100vh - 51px)" content-class="layout-content" embedded>
        <RouterView />
      </n-layout-content>
    </n-layout>
    <LayoutFooter />
  </n-layout>
</template>

<style scoped>
.n-layout-footer {
  padding: 1rem;
}

:deep(.layout-content) {
  padding: 1rem;
}

@media screen and (min-width: 768px) {
  :deep(.layout-content) {
    padding: 2rem;
  }
}
</style>
