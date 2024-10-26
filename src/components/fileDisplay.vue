<template>
  <div class="list-group mb-3">
    <a
      v-for="file in files"
      :key="file.id"
      class="cursor-pointer list-group-item list-group-item-action font-monospace d-flex align-items-center"
      :href="file.url"
      :download="file.name"
      @click="handleClick(file)"
    >
      <i :class="fasFileType(file.type)" />
      <span>{{ file.name }}</span>
      <span class="ms-auto">{{ bytesToHuman(file.size) }}</span>
      <template v-if="trackDownload">
        <i
          v-if="!hasDownloaded[file.id]"
          class="fas fa-fw fa-download ms-2 text-warning"
        />
        <i
          v-else
          class="fas fa-fw fa-circle-check ms-2 text-success"
        />
      </template>
      <template v-if="canDelete">
        <i
          class="fas fa-fw fa-trash ms-2 text-danger"
        />
      </template>
    </a>
  </div>
</template>

<script>
import { bytesToHuman } from '../helpers'

export default {
  data() {
    return {
      hasDownloaded: {},
    }
  },

  methods: {
    bytesToHuman,

    fasFileType(type) {
      return [
        'fas',
        'fa-fw',
        'me-2',
        ...[
          { icon: ['fa-file-pdf'], match: /application\/pdf/ },
          { icon: ['fa-file-audio'], match: /^audio\// },
          { icon: ['fa-file-image'], match: /^image\// },
          { icon: ['fa-file-lines'], match: /^text\// },
          { icon: ['fa-file-video'], match: /^video\// },
          { icon: ['fa-file-zipper'], match: /^application\/(gzip|x-tar|zip)$/ },
          { icon: ['fa-file-circle-question'], match: /.*/ },
        ].filter(el => el.match.test(type))[0].icon,
      ].join(' ')
    },

    handleClick(file) {
      this.$set(this.hasDownloaded, file.id, true)
      this.$emit('fileClicked', file.id)
    },
  },

  name: 'AppFileDisplay',

  props: {
    canDelete: {
      default: false,
      required: false,
      type: Boolean,
    },

    files: {
      required: true,
      type: Array,
    },

    trackDownload: {
      default: true,
      required: false,
      type: Boolean,
    },
  },
}
</script>
