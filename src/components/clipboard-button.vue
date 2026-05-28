<template>
  <button
    v-if="hasClipboard"
    :class="{'btn': true, 'btn-primary': !copyToClipboardSuccess, 'btn-success': copyToClipboardSuccess}"
    :disabled="!content"
    @click="copy"
  >
    <i :class="{'fas fa-fw fa-clipboard': !copyToClipboardSuccess, 'fas fa-fw fa-circle-check': copyToClipboardSuccess}" />
  </button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  computed: {
    hasClipboard(): boolean {
      return Boolean(navigator.clipboard && navigator.clipboard.writeText)
    },
  },

  data() {
    return {
      copyToClipboardSuccess: false,
    }
  },

  methods: {
    copy(): void {
      navigator.clipboard.writeText(this.content)
        .then(() => {
          this.copyToClipboardSuccess = true
          window.setTimeout(() => {
            this.copyToClipboardSuccess = false
          }, 1500)
        })
    },
  },

  name: 'AppClipboardButton',

  props: {
    content: {
      default: null,
      required: false,
      type: String,
    },
  },
})
</script>
