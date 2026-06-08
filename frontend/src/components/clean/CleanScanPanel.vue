<script setup>
defineProps({
  mode: {
    type: String,
    required: true,
    validator: (value) => ['idle', 'scanning'].includes(value),
  },
  progress: {
    type: Number,
    default: 0,
  },
  statusMessage: {
    type: String,
    default: 'Scanning your Mac',
  },
})

defineEmits(['scan'])
</script>

<template>
  <div class="scan-panel">
    <div v-if="mode === 'idle'" class="scan-panel__idle">
      <button
        type="button"
        class="scan-button"
        aria-label="Scan your Mac"
        @click="$emit('scan')"
      >
        <span class="scan-button__rings" aria-hidden="true">
          <span class="scan-button__ring" />
          <span class="scan-button__ring" />
          <span class="scan-button__ring" />
        </span>
        <span class="scan-button__core">
          <i class="pi pi-search scan-button__icon" aria-hidden="true" />
        </span>
      </button>
      <h2 class="scan-panel__title">Scan your Mac</h2>
      <p class="scan-panel__subtitle">Find junk files to free up disk space</p>
    </div>

    <div v-else class="scan-panel__scanning">
      <div class="scan-radar" aria-hidden="true">
        <svg class="scan-radar__track" viewBox="0 0 160 160">
          <circle
            cx="80"
            cy="80"
            r="70"
            fill="none"
            stroke="var(--color-border)"
            stroke-width="3"
          />
          <circle
            cx="80"
            cy="80"
            r="70"
            fill="none"
            stroke="var(--color-accent)"
            stroke-width="3"
            stroke-linecap="round"
            :stroke-dasharray="440"
            :stroke-dashoffset="440 - (440 * progress) / 100"
            transform="rotate(-90 80 80)"
          />
        </svg>
        <span class="scan-radar__sweep" />
        <span class="scan-radar__percent">{{ Math.round(progress) }}%</span>
      </div>
      <p class="scan-panel__status">{{ statusMessage }}</p>
    </div>
  </div>
</template>

<style scoped>
.scan-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 0;
  padding: var(--space-6);
}

.scan-panel__idle {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.scan-button {
  position: relative;
  width: 160px;
  height: 160px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 50%;
}

.scan-button__rings {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.scan-button__ring {
  position: absolute;
  inset: 0;
  border: 2px solid var(--color-accent);
  border-radius: 50%;
  opacity: 0;
  animation: scan-ring-pulse 3s ease-out infinite;
}

.scan-button__ring:nth-child(2) {
  animation-delay: 1s;
}

.scan-button__ring:nth-child(3) {
  animation-delay: 2s;
}

.scan-button__core {
  position: absolute;
  inset: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #d4b896, var(--color-accent));
  box-shadow:
    0 0 0 1px rgba(196, 165, 116, 0.35),
    0 8px 32px rgba(196, 165, 116, 0.25);
  transition: transform var(--transition-fast), box-shadow var(--transition-fast);
}

.scan-button:hover .scan-button__core {
  transform: scale(1.04);
  box-shadow:
    0 0 0 1px rgba(196, 165, 116, 0.5),
    0 12px 40px rgba(196, 165, 116, 0.35);
}

.scan-button:active .scan-button__core {
  transform: scale(0.98);
}

.scan-button__icon {
  font-size: 36px;
  color: var(--color-selection-text);
}

.scan-panel__title {
  margin: var(--space-6) 0 var(--space-2);
  font-size: var(--font-size-hero-metric);
  font-weight: 600;
  color: var(--color-text-primary);
}

.scan-panel__subtitle {
  margin: 0;
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  max-width: 320px;
  line-height: 1.5;
}

.scan-panel__scanning {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.scan-radar {
  position: relative;
  width: 160px;
  height: 160px;
}

.scan-radar__track {
  width: 100%;
  height: 100%;
}

.scan-radar__track circle:last-child {
  transition: stroke-dashoffset 0.25s ease;
}

.scan-radar__sweep {
  position: absolute;
  inset: 10px;
  border-radius: 50%;
  background: conic-gradient(
    from 0deg,
    transparent 0deg,
    transparent 300deg,
    rgba(196, 165, 116, 0.35) 360deg
  );
  animation: scan-radar-sweep 2s linear infinite;
  mask: radial-gradient(circle, transparent 58%, black 59%);
  -webkit-mask: radial-gradient(circle, transparent 58%, black 59%);
}

.scan-radar__percent {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-loader-bright);
}

.scan-panel__status {
  margin: var(--space-6) 0 0;
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  min-height: 1.3em;
}

@keyframes scan-ring-pulse {
  0% {
    transform: scale(0.85);
    opacity: 0.55;
  }

  100% {
    transform: scale(1.55);
    opacity: 0;
  }
}

@keyframes scan-radar-sweep {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .scan-button__ring,
  .scan-radar__sweep {
    animation: none;
  }

  .scan-radar__track circle:last-child {
    transition: none;
  }
}
</style>
