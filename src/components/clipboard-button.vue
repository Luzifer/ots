<template>
  <button
    v-if="hasClipboard"
    :class="{'btn': true, 'btn-primary': !copyToClipboardSuccess, 'btn-success': copyToClipboardSuccess}"
    :disabled="!content"
    @click="copy"
  >
    <i class="fas fa-clipboard" />
  </button>
</template>
<script>
export default {
  computed: {
    hasClipboard() {
      return Boolean(navigator.clipboard && navigator.clipboard.writeText)
    },
  },

  data() {
    return {
      copyToClipboardSuccess: false,
    }
  },

  methods: {
    copy() {
      navigator.clipboard.writeText(this.content)
        .then(() => {
          this.copyToClipboardSuccess = true
          window.setTimeout(() => {
            this.copyToClipboardSuccess = false
          }, 500)
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
}
</script>
