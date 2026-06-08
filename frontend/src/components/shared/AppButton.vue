<script setup>
import { computed } from 'vue'
import LoadingDots from './LoadingDots.vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'secondary', 'ghost', 'danger'].includes(v),
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String,
    default: 'button',
  },
})

defineEmits(['click'])

const loaderTone = computed(() => (['primary', 'danger'].includes(props.variant) ? 'light' : 'brand'))
</script>

<template>
  <button
    :type="type"
    class="app-button"
    :class="[`app-button--${variant}`, { 'app-button--loading': loading }]"
    :disabled="disabled || loading"
    @click="$emit('click', $event)"
  >
    <LoadingDots v-if="loading" size="sm" :tone="loaderTone" />
    <span class="app-button__label" :class="{ 'app-button__label--hidden': loading }">
      <slot />
    </span>
  </button>
</template>

<style scoped>
.app-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  min-height: 32px;
  padding: var(--space-2) var(--space-4);
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  font-family: inherit;
  font-size: var(--font-size-body);
  font-weight: 500;
  line-height: 1;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), opacity var(--transition-fast);
}

.app-button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.app-button--primary {
  background: var(--color-accent);
  color: var(--color-selection-text);
}

.app-button--primary:not(:disabled):hover {
  background: var(--color-accent-pressed);
}

.app-button--primary:not(:disabled):active {
  opacity: 0.85;
}

.app-button--secondary {
  background: var(--color-bg-surface);
  color: var(--color-text-primary);
  border-color: var(--color-border-strong);
}

.app-button--secondary:not(:disabled):hover {
  background: var(--color-bg-elevated);
}

.app-button--secondary:not(:disabled):active {
  opacity: 0.85;
}

.app-button--ghost {
  background: transparent;
  color: var(--color-accent);
  padding: var(--space-1) var(--space-2);
  min-height: 28px;
}

.app-button--ghost:not(:disabled):hover {
  background: var(--color-bg-elevated);
}

.app-button--ghost:not(:disabled):active {
  opacity: 0.85;
}

.app-button--danger {
  background: var(--color-danger);
  color: #fff;
}

.app-button--danger:not(:disabled):hover {
  background: #e03e34;
}

.app-button--danger:not(:disabled):active {
  opacity: 0.85;
}

.app-button__label--hidden {
  visibility: hidden;
  position: absolute;
}
</style>
