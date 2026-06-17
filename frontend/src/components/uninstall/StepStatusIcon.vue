<script setup>
defineProps({
  status: {
    type: String,
    default: 'pending',
    validator: (v) => ['pending', 'running', 'success', 'failed'].includes(v),
  },
})
</script>

<template>
  <span class="step-icon" :class="`step-icon--${status}`" aria-hidden="true">
    <span v-if="status === 'pending'" key="pending" class="step-icon__pending" />

    <svg v-else-if="status === 'running'" key="running" class="step-icon__spinner" viewBox="0 0 24 24">
      <circle class="step-icon__track" cx="12" cy="12" r="9" />
      <circle class="step-icon__arc" cx="12" cy="12" r="9" />
    </svg>

    <svg v-else-if="status === 'success'" key="success" class="step-icon__mark" viewBox="0 0 24 24">
      <circle class="step-icon__ring" cx="12" cy="12" r="11" />
      <path class="step-icon__draw" d="M6.5 12.5 L10.5 16.5 L17.5 8" />
    </svg>

    <svg v-else key="failed" class="step-icon__mark" viewBox="0 0 24 24">
      <circle class="step-icon__ring" cx="12" cy="12" r="11" />
      <path class="step-icon__draw" d="M8.5 8.5 L15.5 15.5 M15.5 8.5 L8.5 15.5" />
    </svg>
  </span>
</template>

<style scoped>
.step-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  flex-shrink: 0;
}

.step-icon__pending {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1.5px solid var(--color-border-strong);
  background: var(--color-bg-inset);
  transition: border-color var(--transition-fast);
}

.step-icon__spinner {
  width: 22px;
  height: 22px;
  animation: step-spin 0.75s linear infinite;
}

.step-icon__track {
  fill: none;
  stroke: var(--color-bg-elevated);
  stroke-width: 2.5;
}

.step-icon__arc {
  fill: none;
  stroke: var(--color-accent);
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-dasharray: 16 41;
}

.step-icon__mark {
  width: 22px;
  height: 22px;
  animation: step-pop 0.34s cubic-bezier(0.2, 0.7, 0.3, 1.35) both;
}

.step-icon__ring {
  stroke-width: 1.5;
}

.step-icon__draw {
  fill: none;
  stroke-width: 2.2;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-dasharray: 26;
  stroke-dashoffset: 26;
  animation: step-draw 0.28s ease 0.12s forwards;
}

.step-icon--success .step-icon__ring {
  fill: var(--color-success-surface);
  stroke: var(--color-success);
}

.step-icon--success .step-icon__draw {
  stroke: var(--color-success);
}

.step-icon--failed .step-icon__ring {
  fill: var(--color-danger-surface);
  stroke: var(--color-danger);
}

.step-icon--failed .step-icon__draw {
  stroke: var(--color-danger);
}

.step-icon--failed {
  animation: step-shake 0.4s ease 0.16s;
}

@keyframes step-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes step-pop {
  0% {
    transform: scale(0.4);
    opacity: 0;
  }
  60% {
    transform: scale(1.08);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes step-draw {
  to {
    stroke-dashoffset: 0;
  }
}

@keyframes step-shake {
  0%,
  100% {
    transform: translateX(0);
  }
  30% {
    transform: translateX(-2px);
  }
  70% {
    transform: translateX(2px);
  }
}

@media (prefers-reduced-motion: reduce) {
  .step-icon__spinner {
    animation: none;
  }

  .step-icon__arc {
    stroke-dasharray: 28 28;
  }

  .step-icon__mark,
  .step-icon--failed {
    animation: none;
  }

  .step-icon__draw {
    stroke-dashoffset: 0;
    animation: none;
  }
}
</style>
