<template>
  <button
    v-if="!customize.disableQRSupport"
    id="secret-url-qrcode"
    ref="qrButton"
    class="btn btn-secondary"
    :disabled="!qrDataURL"
  >
    <i class="fas fa-qrcode" />
  </button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Popover } from 'bootstrap'
import qrcode from 'qrcode'

export default defineComponent({
  computed: {
    customize(): any {
      return window.OTSCustomize || {}
    },
  },

  data() {
    return {
      qrDataURL: null,
    }
  },

  methods: {
    generateQR(): void {
      if (window.OTSCustomize.disableQRSupport) {
        return
      }

      qrcode.toDataURL(this.qrContent, { width: 200 })
        .then(url => {
          this.qrDataURL = url
        })
    },
  },

  mounted(): void {
    this.generateQR()
  },

  name: 'AppQRButton',

  props: {
    qrContent: {
      required: true,
      type: String,
    },
  },

  watch: {
    qrContent() {
      this.generateQR()
    },

    qrDataURL(to: string): void {
      if (this.popover) {
        this.popover.dispose()
      }

      this.popover = new Popover(this.$refs.qrButton, {
        content: () => {
          const img = document.createElement('img')
          img.src = to
          return img
        },

        html: true,
        placement: 'left',
        trigger: 'focus',
      })
    },
  },
})
</script>
