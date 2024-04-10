<script setup lang="ts">
import { ref } from "vue"
import { NFlex, NButton, NIcon, NModal, NQrCode, NDropdown, useMessage, type DropdownOption } from "naive-ui"
import { IconQR, IconLink, IconOpen } from "@/components/icons"

const props = defineProps<{
  link: string
}>()

const currentLink = ref(window.location.href)
const options: DropdownOption[] = [
  {
    label: "Локальная новость",
    key: "internal"
  },
  {
    label: "Новость в источнике",
    key: "external"
  }
]

const message = useMessage()
async function copyToClipboard(link: string) {
  try {
    await navigator.clipboard.writeText(link)
    message.success("Ссылка скопирована в буфер обмена")
  } catch (err) {
    message.error("Не удалось скопировать ссылку")
  }
}

async function onSelectOptionLink(key: string) {
  const clipboardLink = key == "internal" ? currentLink.value : props.link
  return copyToClipboard(clipboardLink)
}

const modalLink = ref("")
const isModalShown = ref(false)
function onSelectOptionQR(key: string) {
  modalLink.value = key == "internal" ? currentLink.value : props.link
  isModalShown.value = true
}
</script>

<template>
  <n-flex vertical>
    <n-button tertiary title="Открыть новость в источнике" tag="a" :href="link" target="_blank">
      <template #icon>
        <n-icon>
          <IconOpen />
        </n-icon>
      </template>
    </n-button>
    <n-dropdown trigger="click" :options="options" @select="onSelectOptionLink">
      <n-button tertiary title="Скопировать ссылку">
        <template #icon>
          <n-icon>
            <IconLink />
          </n-icon>
        </template>
      </n-button>
    </n-dropdown>
    <n-dropdown trigger="click" :options="options" @select="onSelectOptionQR">
      <n-button tertiary title="Показать QR-код">
        <template #icon>
          <n-icon>
            <IconQR />
          </n-icon>
        </template>
      </n-button>
    </n-dropdown>
  </n-flex>
  <n-modal v-model:show="isModalShown">
    <n-qr-code :value="modalLink" :size="200" />
  </n-modal>
</template>
