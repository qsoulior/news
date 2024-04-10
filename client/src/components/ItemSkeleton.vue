<script setup lang="ts">
import { NFlex, NSkeleton, NDivider } from "naive-ui"
</script>

<template>
  <n-flex vertical size="large" style="max-width: 50em; margin: auto">
    <n-flex align="center" justify="space-between" style="margin: 0.3em 0">
      <n-skeleton :sharp="false" width="20%" />
      <n-skeleton :sharp="false" width="25%" />
    </n-flex>
    <n-skeleton :sharp="false" style="width: 70%; height: 22px; margin: 0.4em 0" />
    <div>
      <n-skeleton text :sharp="false" />
      <n-skeleton text :sharp="false" width="60%" />
    </div>
    <n-flex>
      <n-skeleton :sharp="false" width="6em" size="small" />
      <n-skeleton :sharp="false" width="7em" size="small" />
      <n-skeleton :sharp="false" width="4em" size="small" />
      <n-skeleton :sharp="false" width="5em" size="small" />
    </n-flex>
    <n-divider style="margin: 1em 0" />
    <n-flex vertical :size="20">
      <div>
        <n-skeleton :sharp="false" text :repeat="3" />
        <n-skeleton :sharp="false" text width="40%" />
      </div>
      <div>
        <n-skeleton :sharp="false" text :repeat="2" />
        <n-skeleton :sharp="false" text width="70%" />
      </div>
    </n-flex>
    <n-divider style="margin: 1em 0" />
    <n-flex align="center" justify="space-between">
      <n-skeleton :sharp="false" width="30%" />
      <n-flex size="small" align="center">
        <n-skeleton :sharp="false" width="3em" size="small" />
        <n-skeleton :sharp="false" width="8em" size="small" />
      </n-flex>
    </n-flex>
  </n-flex>
</template>
