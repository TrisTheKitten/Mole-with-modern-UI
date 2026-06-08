<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  subtitle: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <header class="page-header">
    <h1 class="page-header__title">{{ title }}</h1>
    <p v-if="subtitle" class="page-header__subtitle">{{ subtitle }}</p>
  </header>
</template>

<style scoped>
.page-header {
  margin-bottom: var(--space-4);
}

.page-header__title {
  margin: 0;
  font-size: var(--font-size-title);
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: -0.02em;
}

.page-header__subtitle {
  margin: var(--space-1) 0 0;
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  line-height: 1.4;
}
</style>
