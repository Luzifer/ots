<template>
  <button
    v-if="!$root.customize.disableQRSupport"
    id="secret-url-qrcode"
    ref="qrButton"
    class="btn btn-secondary"
    :disabled="!qrDataURL"
  >
    <i class="fas fa-qrcode" />
  </button>
</template>
<script>
import { Popover } from 'bootstrap'
import qrcode from 'qrcode'

export default {
  data() {
    return {
      qrDataURL: null,
    }
  },

  methods: {
    generateQR() {
      if (this.$root.customize.disableQRSupport) {
        return
      }

      qrcode.toDataURL(this.qrContent, { width: 200 })
        .then(url => {
          this.qrDataURL = url
        })
    },
  },

  mounted() {
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

    qrDataURL(to) {
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
}
</script>
