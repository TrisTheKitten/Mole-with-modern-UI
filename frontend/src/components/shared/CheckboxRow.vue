<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    default: '',
  },
  size: {
    type: String,
    default: '',
  },
  checked: {
    type: Boolean,
    default: false,
  },
  muted: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['toggle'])

function handleKeydown(event) {
  if (event.key === ' ' || event.key === 'Enter') {
    event.preventDefault()
    emit('toggle')
  }
}
</script>

<template>
  <div
    class="checkbox-row"
    :class="{ 'checkbox-row--checked': checked, 'checkbox-row--muted': muted }"
    role="checkbox"
    :aria-checked="checked"
    :aria-label="title"
    tabindex="0"
    @click="emit('toggle')"
    @keydown="handleKeydown"
  >
    <span class="checkbox-row__box" aria-hidden="true">
      <svg v-if="checked" class="checkbox-row__check" viewBox="0 0 12 12" fill="none">
        <path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
    </span>
    <div class="checkbox-row__info">
      <span class="checkbox-row__title">{{ title }}</span>
      <span v-if="description" class="checkbox-row__description">{{ description }}</span>
    </div>
    <span v-if="size" class="checkbox-row__size">{{ size }}</span>
    <div v-if="$slots.trailing" class="checkbox-row__trailing" @click.stop>
      <slot name="trailing" />
    </div>
  </div>
</template>

<style scoped>
.checkbox-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-height: 52px;
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.checkbox-row:hover {
  background: var(--color-bg-elevated);
  border-color: var(--color-border-strong);
}

.checkbox-row--checked,
.checkbox-row--checked:hover {
  background: var(--color-accent-subtle);
  border-color: var(--color-accent-border);
}

.checkbox-row--muted {
  opacity: 0.55;
}

.checkbox-row--muted .checkbox-row__size {
  color: var(--color-text-tertiary);
}

.checkbox-row__box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  border: 1.5px solid var(--color-border-strong);
  border-radius: 4px;
  background: var(--color-bg-app);
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.checkbox-row--checked .checkbox-row__box {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: #fff;
}

.checkbox-row__check {
  width: 12px;
  height: 12px;
}

.checkbox-row__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.checkbox-row__title {
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-primary);
  line-height: 1.3;
}

.checkbox-row__description {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.checkbox-row__size {
  flex-shrink: 0;
  font-size: var(--font-size-body);
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  color: var(--color-text-primary);
}

.checkbox-row__trailing {
  flex-shrink: 0;
}
</style>
