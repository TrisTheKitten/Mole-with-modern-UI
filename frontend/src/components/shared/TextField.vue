<script setup>
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  mono: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="text-field">
    <label v-if="label" class="text-field__label">{{ label }}</label>
    <input
      class="text-field__input"
      :class="{ 'text-field__input--mono': mono }"
      type="text"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="$emit('update:modelValue', $event.target.value)"
    />
  </div>
</template>

<style scoped>
.text-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
  min-width: 0;
}

.text-field__label {
  font-size: var(--font-size-caption);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.text-field__input {
  width: 100%;
  min-height: 32px;
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-family: inherit;
  font-size: var(--font-size-body);
  transition: border-color var(--transition-fast);
}

.text-field__input--mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.text-field__input:focus {
  outline: none;
  border-color: var(--color-accent);
}

.text-field__input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.text-field__input::placeholder {
  color: var(--color-text-tertiary);
}
</style>
