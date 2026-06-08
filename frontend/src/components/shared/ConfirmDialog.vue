<script setup>
import AppButton from './AppButton.vue'

defineProps({
  show: Boolean,
  title: String,
  message: String,
  items: {
    type: Array,
    default: () => [],
  },
  confirmText: {
    type: String,
    default: 'Confirm',
  },
  cancelText: {
    type: String,
    default: 'Cancel',
  },
  destructive: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['confirm', 'cancel', 'update:show'])

function handleConfirm() {
  emit('confirm')
  emit('update:show', false)
}

function handleCancel() {
  emit('cancel')
  emit('update:show', false)
}
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-overlay" @click="handleCancel">
      <div
        class="modal-content"
        :class="{ 'modal-content--with-items': items.length > 0 }"
        role="dialog"
        aria-modal="true"
        @click.stop
      >
        <h3 class="modal-title">{{ title }}</h3>
        <p v-if="message" class="modal-message">{{ message }}</p>
        <ul v-if="items.length > 0" class="modal-items">
          <li v-for="(item, index) in items" :key="index" class="modal-item">{{ item }}</li>
        </ul>
        <div class="modal-actions">
          <AppButton v-if="cancelText" variant="secondary" @click="handleCancel">{{ cancelText }}</AppButton>
          <AppButton :variant="destructive ? 'danger' : 'primary'" @click="handleConfirm">{{ confirmText }}</AppButton>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--color-scrim);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  padding: var(--space-6);
  max-width: 400px;
  width: 90%;
  box-shadow: var(--shadow-modal);
}

.modal-content--with-items {
  max-width: 520px;
}

.modal-title {
  margin: 0 0 var(--space-3);
  font-size: 17px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.modal-message {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-body);
  line-height: 1.5;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.modal-items {
  margin: var(--space-3) 0 0;
  padding: var(--space-2) var(--space-3);
  max-height: 180px;
  overflow-y: auto;
  list-style: none;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.modal-item {
  padding: var(--space-1) 0;
  color: var(--color-text-secondary);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: var(--font-size-caption);
  line-height: 1.4;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.modal-item + .modal-item {
  border-top: 1px solid var(--color-border);
}

.modal-actions {
  display: flex;
  gap: var(--space-2);
  justify-content: flex-end;
  margin-top: var(--space-6);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
  transition: transform var(--transition-fast);
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.96);
}
</style>
