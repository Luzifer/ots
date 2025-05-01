<template>
  <div class="card border-primary-subtle mb-3">
    <div
      class="card-header bg-primary-subtle"
      v-html="$t('title-explanation')"
    />
    <div class="card-body">
      <ul>
        <li
          v-for="(explanation, idx) in $tm('items-explanation')"
          :key="`idx${idx}`"
        >
          {{ explanation }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'AppExplanation',
})
</script>
