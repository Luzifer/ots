<template>
  <nav class="navbar navbar-expand-lg bg-primary-subtle">
    <div class="container-fluid">
      <a
        class="navbar-brand"
        href="#"
        @click.prevent="$emit('navigate', '/')"
      >
        <i
          v-if="!appIcon"
          class="fas fa-user-secret me-1"
        />
        <img
          v-else
          class="me-1"
          :src="appIcon"
        >
        <span v-if="!customize.disableAppTitle">{{ customize.appTitle }}</span>
      </a>

      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarSupportedContent"
        aria-controls="navbarSupportedContent"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon" />
      </button>

      <div
        id="navbarSupportedContent"
        class="collapse navbar-collapse"
      >
        <ul class="navbar-nav ms-auto mb-2 mb-lg-0 me-2">
          <li class="nav-item">
            <a
              class="nav-link"
              href="#"
              @click.prevent="$emit('navigate', '/explanation')"
            >
              <i class="fas fa-circle-info" /> {{ $t('btn-show-explanation') }}
            </a>
          </li>
          <li class="nav-item">
            <a
              class="nav-link"
              href="#"
              @click.prevent="$emit('navigate', '/')"
            >
              <i class="fas fa-plus" /> {{ $t('btn-new-secret') }}
            </a>
          </li>
        </ul>
        <form
          v-if="!customize.disableThemeSwitcher"
          class="d-flex align-items-center btn-group"
        >
          <input
            id="theme-light"
            v-model="intTheme"
            type="radio"
            name="theme"
            class="btn-check"
            value="light"
          >
          <label
            class="btn btn-outline-secondary btn-sm"
            for="theme-light"
          >
            <i class="fas fa-sun" />
          </label>

          <input
            id="theme-auto"
            v-model="intTheme"
            type="radio"
            name="theme"
            class="btn-check"
            value="auto"
          >
          <label
            class="btn btn-outline-secondary btn-sm"
            for="theme-auto"
          >
            {{ $t('btn-theme-switcher-auto') }}
          </label>

          <input
            id="theme-dark"
            v-model="intTheme"
            type="radio"
            name="theme"
            class="btn-check"
            value="dark"
          >
          <label
            class="btn btn-outline-secondary btn-sm"
            for="theme-dark"
          >
            <i class="fas fa-moon" />
          </label>
        </form>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  computed: {
    appIcon(): string {
      // Use specified icon or fall back to null
      const appIcon = this.$parent.customize.appIcon || null
      // Use specified icon or fall back to light-mode appIcon (which might be null)
      const darkIcon = this.$parent.customize.appIconDark || appIcon

      return this.$root.theme === 'dark' ? darkIcon : appIcon
    },

    customize(): any {
      return this.$parent.customize || {}
    },
  },

  data() {
    return {
      intTheme: '',
    }
  },

  emits: ['navigate', 'update:theme'],

  mounted(): void {
    this.intTheme = this.theme
  },

  name: 'AppNavbar',

  props: {
    theme: {
      required: true,
      type: String,
    },
  },

  watch: {
    intTheme(to: string, from: string): void {
      if (to === from) {
        return
      }

      this.$emit('update:theme', to)
    },

    theme(to: string, from: string): void {
      if (to === from) {
        return
      }

      this.intTheme = to
    },
  },
})
</script>
