<script setup lang="ts">
import { ref } from "vue"
import { type GlobalTheme, NConfigProvider, NGlobalStyle } from "naive-ui"
import LayoutContainer from "@/components/LayoutContainer.vue"

const theme = ref<GlobalTheme | null>()
</script>

<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <LayoutContainer v-model:theme="theme" />
  </n-config-provider>
</template>
