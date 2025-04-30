<template>
  <nav class="navbar navbar-expand-lg bg-primary-subtle">
    <div class="container-fluid">
      <a
        class="navbar-brand"
        href="#"
        @click.prevent="$root.navigate('/')"
      >
        <i
          v-if="!appIcon"
          class="fas fa-user-secret mr-1"
        />
        <img
          v-else
          class="mr-1"
          :src="appIcon"
        >
        <span v-if="!$root.customize.disableAppTitle">{{ $root.customize.appTitle }}</span>
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
              @click.prevent="$root.navigate('/explanation')"
            >
              <i class="fas fa-circle-info" /> {{ $t('btn-show-explanation') }}
            </a>
          </li>
          <li class="nav-item">
            <a
              class="nav-link"
              href="#"
              @click.prevent="$root.navigate('/')"
            >
              <i class="fas fa-plus" /> {{ $t('btn-new-secret') }}
            </a>
          </li>
        </ul>
        <form
          v-if="!$root.customize.disableThemeSwitcher"
          class="d-flex align-items-center btn-group"
        >
          <input
            id="theme-light"
            v-model="$root.theme"
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
            v-model="$root.theme"
            type="radio"
            name="theme"
            class="btn-check"
            value="auto"
          >
          <label
            class="btn btn-outline-secondary btn-sm"
            for="theme-auto"
          >
            Auto
          </label>

          <input
            id="theme-dark"
            v-model="$root.theme"
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

<script>
export default {
  computed: {
    appIcon() {
      // Use specified icon or fall back to null
      const appIcon = this.$root.customize.appIcon || null
      // Use specified icon or fall back to light-mode appIcon (which might be null)
      const darkIcon = this.$root.customize.appIconDark || appIcon

      return this.$root.theme === 'dark' ? darkIcon : appIcon
    },
  },

  name: 'AppNavbar',
}
</script>
