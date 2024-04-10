<script setup lang="ts">
import { ref } from "vue"
import { type GlobalTheme, NConfigProvider, NGlobalStyle, NMessageProvider } from "naive-ui"
import LayoutContainer from "@/components/LayoutContainer.vue"

const theme = ref<GlobalTheme | null>()
</script>

<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <n-message-provider>
      <LayoutContainer v-model:theme="theme" />
    </n-message-provider>
  </n-config-provider>
</template>
