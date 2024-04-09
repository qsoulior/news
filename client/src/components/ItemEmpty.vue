<script setup lang="ts">
import { NResult, NButton, NIcon } from "naive-ui"
import { IconHome } from "@/components/icons"
</script>

<template>
  <n-result status="404" title="404 Not Found" description="Новость не найдена" style="margin: 2rem 0">
    <template #footer>
      <router-link :to="{ name: 'list' }" #="{ navigate, href }" custom>
        <n-button tertiary tag="a" :href="href" @click="navigate">
          <template #default> Вернуться на главную страницу </template>
          <template #icon>
            <n-icon>
              <IconHome />
            </n-icon>
          </template>
        </n-button>
      </router-link>
    </template>
  </n-result>
</template>
