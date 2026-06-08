<script setup>
import { ref, onMounted } from 'vue'

const toasts = ref([])
let nextId = 0

onMounted(() => {
  window.addEventListener('show-toast', (e) => {
    const { message, type = 'info', duration = 5000 } = e.detail
    const id = nextId++

    toasts.value.push({ id, message, type })

    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  })
})

function close(id) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}
</script>

<template>
  <div class="toast-container">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :class="['toast', toast.type]"
      @click="close(toast.id)"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.toast {
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-strong);
  box-shadow: var(--shadow-panel);
  font-size: var(--font-size-body);
  line-height: 1.4;
  cursor: pointer;
  max-width: 400px;
  animation: slideIn 0.3s;
}

.toast.error {
  background: var(--color-danger-surface);
  color: var(--color-danger);
  border-color: var(--color-danger-border);
}

.toast.info {
  background: var(--color-accent-subtle);
  color: var(--color-accent);
  border-color: var(--color-accent-border);
}

.toast.success {
  background: var(--color-success-surface);
  color: var(--color-success);
  border-color: var(--color-success-border);
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>
