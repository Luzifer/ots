<template>
  <!-- Creation disabled -->
  <div
    v-if="!showCreateForm"
    class="card border-info-subtle mb-3"
  >
    <div
      class="card-header bg-info-subtle"
      v-html="$t('title-secret-create-disabled')"
    />
    <div
      class="card-body"
      v-html="$t('text-secret-create-disabled')"
    />
  </div>

  <!-- Creation possible -->
  <div
    v-else
    class="card border-primary-subtle mb-3"
  >
    <div
      class="card-header bg-primary-subtle"
      v-html="$t('title-new-secret')"
    />
    <div class="card-body">
      <form
        class="row"
        @submit.prevent="createSecret"
      >
        <div class="col-12 mb-3">
          <label for="createSecretData">{{ $t('label-secret-data') }}</label>
          <grow-area
            id="createSecretData"
            v-model="secret"
            class="form-control"
            :rows="2"
            @paste-file="handlePasteFile"
          />
        </div>
        <div
          v-if="!customize.disableFileAttachment"
          class="col-12 mb-3"
        >
          <label for="createSecretFiles">{{ $t('label-secret-files') }}</label>
          <input
            id="createSecretFiles"
            ref="createSecretFiles"
            class="form-control"
            type="file"
            multiple
            :accept="customize.acceptedFileTypes"
            @change="handleSelectFiles"
          >
          <div class="form-text">
            {{ $t('text-max-filesize', { maxSize: bytesToHuman(maxFileSize) }) }}
          </div>
          <div
            v-if="invalidFilesSelected"
            class="alert alert-danger"
          >
            {{ $t('text-invalid-files-selected') }}
          </div>
          <div
            v-else-if="maxFileSizeExceeded"
            class="alert alert-danger"
          >
            {{ $t('text-max-filesize-exceeded', { curSize: bytesToHuman(fileSize), maxSize: bytesToHuman(maxFileSize) }) }}
          </div>
          <FilesDisplay
            v-if="attachedFiles.length > 0"
            class="mt-3"
            :can-delete="true"
            :track-download="false"
            :files="attachedFiles"
            @file-clicked="deleteFile"
          />
        </div>
        <div class="col-md-6 col-12 order-2 order-md-1">
          <button
            type="submit"
            class="btn btn-success"
            :disabled="!canCreate"
          >
            <template v-if="!createRunning">
              {{ $t('btn-create-secret') }}
            </template>
            <template v-else>
              <i class="fa-solid fa-spinner fa-spin-pulse" />
              {{ $t('btn-create-secret-processing') }}
            </template>
          </button>
        </div>
        <div
          v-if="!customize.disableExpiryOverride"
          class="col-md-6 col-12 order-1 order-md-2"
        >
          <div class="row mb-3 justify-content-end">
            <label
              class="col-md-6 col-form-label text-md-end"
              for="createSecretExpiry"
            >{{ $t('label-expiry') }}</label>
            <div class="col-md-6">
              <select
                v-model="selectedExpiry"
                class="form-select"
              >
                <option
                  v-for="opt in expiryChoices"
                  :key="opt.value || 'null'"
                  :value="opt.value"
                >
                  {{ opt.text }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import {
  bytesToHuman,
  durationToSeconds,
} from '../helpers'
import appCrypto from '../crypto.ts'
import { defineComponent } from 'vue'
import FilesDisplay from './fileDisplay.vue'
import GrowArea from './growarea.vue'
import OTSMeta from '../ots-meta'

const defaultExpiryChoices = [
  90 * 86400, // 90 days
  30 * 86400, // 30 days
  7 * 86400, // 7 days
  3 * 86400, // 3 days
  24 * 3600, // 1 day
  12 * 3600, // 12 hours
  4 * 3600, // 4 hours
  60 * 60, // 1 hour
  30 * 60, // 30 minutes
  5 * 60, // 5 minutes
]

const defaultExpiryChoicesHuman = [
  '90d',
  '30d',
  '7d',
  '3d',
  '24h', // or 1d, equivalent
  '12h',
  '4h',
  '1h',
  '30m',
  '5m',
]

/*
 * We define an internal max file-size which cannot get exceeded even
 * though the server might accept more: at around 70 MiB the base64
 * encoding broke and nothing works anymore. This might be fixed by
 * changing how the base64 implementation works (maybe use a WASM
 * object?) or switching to a browser-native implementation in case
 * that will appear somewhen in the future but for now we just "fix"
 * the issue by disallowing bigger files.
 */
const internalMaxFileSize = 64 * 1024 * 1024 // 64 MiB

const passwordCharset = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
const passwordLength = 20

export default defineComponent({
  components: { FilesDisplay, GrowArea },

  computed: {
    canCreate(): boolean {
      return (this.secret.trim().length > 0 || this.selectedFileMeta.length > 0) && !this.maxFileSizeExceeded && !this.invalidFilesSelected
    },

    customize(): any {
      return window.OTSCustomize || {}
    },

    expiryChoices(): Record<string, string | null>[] {
      if (this.$root.customize.expiryChoicesHuman) {
        return this.expiryChoicesHuman
      }

      const choices = [{ text: this.$t('expire-default'), value: null as string | null }]

      for (const choice of this.customize.expiryChoices || defaultExpiryChoices) {
        if (window.maxSecretExpire > 0 && choice > window.maxSecretExpire) {
          continue
        }

        const option = { text: '', value: choice }
        if (choice >= 86400) {
          option.text = this.$t('expire-n-days', Math.round(choice / 86400))
        } else if (choice >= 3600) {
          option.text = this.$t('expire-n-hours', Math.round(choice / 3600))
        } else if (choice >= 60) {
          option.text = this.$t('expire-n-minutes', Math.round(choice / 60))
        } else {
          option.text = this.$t('expire-n-seconds', choice)
        }

        choices.push(option)
      }

      return choices
    },

    expiryChoicesHuman() {
      const choices = []
      if (!this.hasValidDefaultExpiryHuman()) {
        choices.push({ text: this.$t('expire-default'), value: null })
      }

      for (const choice of this.$root.customize.expiryChoicesHuman || defaultExpiryChoicesHuman) {
        const option = { value: choice }

        const unit = choice.slice(-1)
        const amount = parseInt(choice.slice(0, -1), 10)

        option.text = this._getTextForAmount(unit, amount)

        choices.push(option)
      }

      return choices
    },

    invalidFilesSelected(): boolean {
      if (this.customize.acceptedFileTypes === '') {
        // No limitation configured, no need to check
        return false
      }

      const accepted = this.customize.acceptedFileTypes.split(',')
      for (const fm of this.selectedFileMeta) {
        let isAccepted = false

        for (const a of accepted) {
          isAccepted ||= this.isAcceptedBy(fm, a)
        }

        if (!isAccepted) {
          // Well we only needed one rejected
          return true
        }
      }

      // We found no reason to reject: This is fine!
      return false
    },

    isSecureEnvironment(): boolean {
      return Boolean(window.crypto.subtle)
    },

    maxFileSize(): number {
      return this.customize.maxAttachmentSizeTotal === 0 ? internalMaxFileSize : Math.min(internalMaxFileSize, this.customize.maxAttachmentSizeTotal)
    },

    maxFileSizeExceeded(): boolean {
      return this.fileSize > this.maxFileSize
    },

    showCreateForm(): boolean {
      return this.canWrite && this.isSecureEnvironment
    },
  },

  created(): void {
    this.checkWriteAccess()

    this.$root.$watch(
      'customize',
      newVal => {
        if (newVal) {
          this.initExpiry()
        }
      },
      { immediate: true },
    )
  },

  data() {
    return {
      attachedFiles: [],
      canWrite: null,
      createRunning: false,
      expiryInitialized: false,
      fileSize: 0,
      secret: '',
      securePassword: null,
      selectedExpiry: null,
      selectedFileMeta: [],
    }
  },

  emits: ['error', 'navigate'],

  methods: {
    _getTextForAmount(unit, amount) {
      switch (unit) {
      case 'd':
        return this.$tc('expire-n-days', amount)
      case 'h':
        return this.$tc('expire-n-hours', amount)
      case 'm':
        return this.$tc('expire-n-minutes', amount)
      case 's':
        return this.$tc('expire-n-seconds', amount)
      }

      return amount
    },

    bytesToHuman,

    checkWriteAccess(): Promise<void> {
      return fetch('api/isWritable', {
        credentials: 'same-origin',
        method: 'GET',
        redirect: 'error',
      })
        .then(resp => {
          if (resp.status !== 204) {
            throw new Error(`unexpected status: ${resp.status}`)
          }
          this.canWrite = true
        })
        .catch(() => {
          this.canWrite = false
        })
    },

    // createSecret executes the secret creation after encrypting the secret
    createSecret(): void {
      if (!this.canCreate) {
        return
      }

      // Encoding large files takes a while, prevent duplicate click on "create"
      this.createRunning = true

      this.securePassword = [...window.crypto.getRandomValues(new Uint8Array(passwordLength))]
        .map(n => passwordCharset[n % passwordCharset.length])
        .join('')

      const meta = new OTSMeta()
      meta.secret = this.secret

      if (this.attachedFiles.length > 0) {
        for (const f of this.attachedFiles) {
          meta.files.push(f.fileObj)
        }
      }

      meta.serialize()
        .then(secret => appCrypto.enc(secret, this.securePassword))
        .then(secret => {
          let reqURL = 'api/create'
          if (this.selectedExpiry !== null) {
            reqURL = `api/create?expire=${durationToSeconds(this.selectedExpiry)}`
          }

          return fetch(reqURL, {
            body: JSON.stringify({ secret }),
            headers: {
              'content-type': 'application/json',
            },
            method: 'POST',
          })
            .then(resp => {
              if (resp.status !== 201) {
                // Server says "no"
                this.$emit('error', this.$t('alert-something-went-wrong'))
                return
              }

              resp.json()
                .then(data => {
                  this.$emit('navigate', {
                    path: '/display-secret-url',
                    query: {
                      expiresAt: data.expires_at,
                      secretId: data.secret_id,
                      securePassword: this.securePassword,
                    },
                  })
                })
            })
            .catch(() => {
              // Network error
              this.$emit('error', this.$t('alert-something-went-wrong'))
            })
        })
    },

    deleteFile(fileId: string): void {
      this.attachedFiles = [...this.attachedFiles].filter(file => file.id !== fileId)
      this.updateFileMeta()
    },

    handlePasteFile(file: File): void {
      this.attachedFiles.push({
        fileObj: file,
        id: window.crypto.randomUUID(),
        name: file.name,
        size: file.size,
        type: file.type,
      })
      this.updateFileMeta()
    },

    handleSelectFiles(): void {
      for (const file of this.$refs.createSecretFiles.files) {
        this.attachedFiles.push({
          fileObj: file,
          id: window.crypto.randomUUID(),
          name: file.name,
          size: file.size,
          type: file.type,
        })
      }
      this.updateFileMeta()

      this.$refs.createSecretFiles.value = ''
    },

    hasValidDefaultExpiryHuman() {
      const defaultExpiry = this.$root.customize.defaultExpiryHuman || false
      if (defaultExpiry === false) {
        return false
      }

      if (!this.$root.customize.expiryChoicesHuman) {
        return false
      }

      return this.$root.customize.expiryChoicesHuman.includes(defaultExpiry)
    },

    initExpiry() {
      const match = document.cookie.match(/(?:^|;\s*)selectedExpiry=([^;]*)/)
      this.selectedExpiry = match
        ? decodeURIComponent(match[1])
        : this.$root.customize?.defaultExpiryHuman || null

      if (!this.$root.customize?.expiryChoicesHuman) {
        return
      }
      if (!this.$root.customize?.expiryChoicesHuman.includes(this.selectedExpiry)) {
        this.selectedExpiry = null
      }
    },

    isAcceptedBy(fileMeta: any, accept: string): boolean {
      if (/^(?:[a-z]+|\*)\/(?:[a-zA-Z0-9.+_-]+|\*)$/.test(accept)) {
        // That's likely supposed to be a mime-type
        return RegExp(`^${accept.replaceAll('*', '.*')}$`).test(fileMeta.type)
      } else if (/^\.[a-z.]+$/.test(accept)) {
        // That should be a file extension
        return fileMeta.name.endsWith(accept)
      }

      // What exactly is it then? At least it can't accept anything.
      return false
    },

    updateFileMeta(): void {
      let cumSize = 0
      for (const f of this.attachedFiles) {
        cumSize += f.size
      }

      this.fileSize = cumSize
      this.selectedFileMeta = this.attachedFiles.map(file => ({
        name: file.name,
        type: file.type,
      }))
    },
  },

  name: 'AppCreate',

  watch: {
    selectedExpiry(newVal) {
      if (!this.expiryInitialized) {
        this.expiryInitialized = true

        return
      }

      document.cookie = `selectedExpiry=${newVal || ''}; path=/; max-age=${60 * 60 * 24 * 365}`
    },
  },
})
</script>
