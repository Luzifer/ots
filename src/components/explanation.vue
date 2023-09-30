<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="card border-primary-subtle mb-3">
    <div
      class="card-header bg-primary-subtle"
      v-html="$t('title-explanation')"
    />
    <div class="card-body">
      <ul>
        <li
          v-for="(explanation, idx) in $t('items-explanation')"
          :key="`idx${idx}`"
        >
          {{ explanation }}
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
export default { name: 'AppExplanation' }
</script>
