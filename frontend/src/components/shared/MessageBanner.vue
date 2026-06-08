<script setup>
defineProps({
  type: {
    type: String,
    default: 'error',
    validator: (v) => ['error', 'success'].includes(v),
  },
  message: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <div class="message-banner" :class="`message-banner--${type}`" role="alert">
    {{ message }}
  </div>
</template>

<style scoped>
.message-banner {
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
  font-size: var(--font-size-body);
  line-height: 1.4;
}

.message-banner--error {
  background: var(--color-danger-surface);
  color: var(--color-danger);
  border: 1px solid var(--color-danger-border);
}

.message-banner--success {
  background: var(--color-success-surface);
  color: var(--color-success);
  border: 1px solid var(--color-success-border);
}
</style>
