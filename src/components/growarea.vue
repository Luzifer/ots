<template>
  <textarea
    ref="area"
    v-model="data"
    style="resize: none;"
    @paste="handlePaste"
  />
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  created() {
    this.data = this.value
  },

  data() {
    return {
      data: '',
    }
  },

  emits: ['input', 'pasteFile'],

  methods: {
    changeSize(): void {
      if (!this.$refs.area) {
        return
      }

      const verticalBorderSize = this.getStyle('borderTopWidth') + this.getStyle('borderBottomWidth') || 0
      const verticalPaddingSize = this.getStyle('paddingTop') + this.getStyle('paddingBottom') || 0

      const smallestHeight = this.getStyle('lineHeight') * this.rows + verticalBorderSize + verticalPaddingSize
      this.$refs.area.style.height = `${smallestHeight}px`

      const newHeight = this.$refs.area.scrollHeight + verticalBorderSize
      this.$refs.area.style.height = `${newHeight}px`
    },

    getStyle(name: string): number {
      return parseInt(getComputedStyle(this.$refs.area, null)[name])
    },

    handlePaste(evt: ClipboardEvent): void {
      if ([...evt.clipboardData?.items || []]
        .filter(item => item.kind !== 'string')
        .length === 0) {
        return
      }

      /*
       * We have something else than text, prevent using clipboard and
       * pasting and emit an event containing the file data
       */
      evt.stopPropagation()
      evt.preventDefault()

      for (const item of evt.clipboardData?.items || []) {
        if (item.kind === 'string') {
          continue
        }

        this.$emit('pasteFile', item.getAsFile())
      }
    },
  },

  mounted(): void {
    this.changeSize()
  },

  name: 'GrowArea',

  props: {
    rows: {
      default: 4,
      type: Number,
    },

    value: {
      default: '',
      type: String,
    },
  },

  watch: {
    data(to, from) {
      this.changeSize()
      if (to !== from) {
        this.$emit('input', to)
      }
    },

    value(to) {
      if (to !== this.data) {
        this.data = to
      }
    },
  },
})
</script>
