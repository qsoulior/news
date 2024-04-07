<script setup lang="ts">
import { onMounted, ref, watch } from "vue"
import { RouterView } from "vue-router"
import { type GlobalTheme, NLayout, NLayoutContent, darkTheme, useOsTheme } from "naive-ui"
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
  <n-layout>
    <n-layout style="min-height: 100vh">
      <LayoutHeader v-model:theme="themeType" />
      <n-layout-content content-style="padding: 2rem;" embedded>
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
</style>
