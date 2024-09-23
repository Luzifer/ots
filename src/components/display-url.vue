<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="card border-success-subtle mb-3">
    <div
      class="card-header bg-success-subtle"
      v-html="$t('title-secret-created')"
    />
    <div
      v-if="!burned"
      class="card-body"
    >
      <p v-html="$t('text-pre-url')" />
      <div class="input-group mb-3">
        <input
          ref="secretUrl"
          class="form-control"
          type="text"
          readonly
          :value="secretUrl"
          @focus="$refs.secretUrl.select()"
        >
        <app-clipboard-button
          :content="secretUrl"
          :title="$t('tooltip-copy-to-clipboard')"
        />
        <app-qr-button :qr-content="secretUrl" />
        <button
          class="btn btn-danger"
          :title="$t('tooltip-burn-secret')"
          @click="burnSecret"
        >
          <i class="fas fa-fire fa-fw" />
        </button>
      </div>
      <p v-html="$t('text-burn-hint')" />
      <p v-if="expiresAt">
        {{ $t('text-burn-time') }}
        <strong>{{ expiresAt.toLocaleString() }}</strong>
      </p>
    </div>
    <div
      v-else
      class="card-body"
    >
      {{ $t('text-secret-burned') }}
    </div>
  </div>
</template>
<script>
import appClipboardButton from './clipboard-button.vue'
import appQrButton from './qr-button.vue'

export default {
  components: { appClipboardButton, appQrButton },
  computed: {
    secretUrl() {
      return [
        window.location.href.split('#')[0],
        encodeURIComponent([
          this.secretId,
          this.securePassword,
        ].join('|')),
      ].join('#')
    },
  },

  data() {
    return {
      burned: false,
      popover: null,
    }
  },

  methods: {
    burnSecret() {
      return fetch(`api/get/${this.secretId}`)
        .then(() => {
          this.burned = true
        })
    },
  },

  mounted() {
    // Give the interface a moment to transistion and focus
    window.setTimeout(() => this.$refs.secretUrl.focus(), 100)
  },

  name: 'AppDisplayURL',
  props: {
    expiresAt: {
      default: null,
      required: false,
      type: Date,
    },

    secretId: {
      required: true,
      type: String,
    },

    securePassword: {
      required: true,
      type: String,
    },
  },
}
</script>
