<script setup lang="ts">
import { NFlex, NSkeleton, NCard } from "naive-ui"
</script>

<template>
  <n-card size="small">
    <n-flex vertical>
      <n-skeleton :sharp="false" style="width: 60%; margin: 0.2em 0" />
      <n-skeleton :sharp="false" style="margin: 0.2em 0" />
      <n-flex justify="space-between" style="margin: 0.2em 0">
        <n-skeleton :sharp="false" style="width: 30%" />
        <n-skeleton :sharp="false" style="width: 20%" />
      </n-flex>
    </n-flex>
  </n-card>
</template>
