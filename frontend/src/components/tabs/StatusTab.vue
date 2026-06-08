<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import {
  StatusGetMetrics,
  StatusStartMonitoring,
  StatusStopMonitoring
} from '../../../wailsjs/go/main/App'
import { handleError } from '../../utils/errorHandler'
import PageHeader from '../shared/PageHeader.vue'
import AppButton from '../shared/AppButton.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import EmptyState from '../shared/EmptyState.vue'
import StatusHealthSummary from '../status/StatusHealthSummary.vue'
import ResourceSection from '../status/ResourceSection.vue'

const metrics = ref(null)
const monitoring = ref(false)
const loading = ref(false)
const refreshing = ref(false)

const healthScore = computed(() => calculateHealthScore())

const healthStatus = computed(() => {
  const score = healthScore.value
  if (score >= 90) return 'excellent'
  if (score >= 75) return 'good'
  if (score >= 60) return 'fair'
  if (score >= 40) return 'poor'
  return 'critical'
})

const gpuUsage = computed(() => {
  const usage = metrics.value?.gpu?.usage
  if (usage === null || usage === undefined) return null
  const num = Number(usage)
  if (isNaN(num) || num < 0) return null
  return num
})

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatPercent(value) {
  if (value === null || value === undefined) return '0.0%'
  const num = Number(value)
  if (isNaN(num)) return '0.0%'
  return num.toFixed(1) + '%'
}

function formatRate(mbPerSec) {
  if (!mbPerSec || mbPerSec === 0) return '0 MB/s'
  const num = Number(mbPerSec)
  if (isNaN(num)) return '0 MB/s'
  if (num < 1) return (num * 1024).toFixed(2) + ' KB/s'
  return num.toFixed(2) + ' MB/s'
}

function formatTemp(celsius) {
  if (celsius === null || celsius === undefined) return '0.0°C'
  const num = Number(celsius)
  if (isNaN(num)) return '0.0°C'
  return num.toFixed(1) + '°C'
}

function getUsageColor(percent) {
  if (percent >= 90) return 'var(--color-danger)'
  if (percent >= 75) return 'var(--color-warning)'
  return 'var(--color-success)'
}

function getBatteryColor(percent) {
  if (percent < 20) return 'var(--color-danger)'
  if (percent < 40) return 'var(--color-warning)'
  return 'var(--color-success)'
}

function calculateHealthScore() {
  if (!metrics.value) return 0

  let score = 100
  const m = metrics.value

  if (m.cpu?.totalPercent) {
    if (m.cpu.totalPercent > 90) score -= 30
    else if (m.cpu.totalPercent > 70) score -= 20
    else if (m.cpu.totalPercent > 50) score -= 10
  }

  if (m.memory?.percent) {
    if (m.memory.percent > 90) score -= 25
    else if (m.memory.percent > 75) score -= 15
    else if (m.memory.percent > 60) score -= 5
  }

  if (m.disk?.percent) {
    if (m.disk.percent > 95) score -= 20
    else if (m.disk.percent > 85) score -= 10
    else if (m.disk.percent > 75) score -= 5
  }

  if (m.cpu?.temperature) {
    if (m.cpu.temperature > 85) score -= 15
    else if (m.cpu.temperature > 75) score -= 10
    else if (m.cpu.temperature > 65) score -= 5
  }

  if (m.battery?.level !== undefined) {
    if (m.battery.level < 20) score -= 10
    else if (m.battery.level < 40) score -= 5
  }

  return Math.max(0, Math.min(100, score))
}

async function startMonitoring() {
  loading.value = true
  try {
    await StatusStartMonitoring(2)
    monitoring.value = true
    try {
      await fetchMetrics()
    } catch {
      console.warn('Failed to fetch initial metrics, will retry via monitoring events')
    }
  } catch (error) {
    handleError(error, 'Start Monitoring')
    monitoring.value = false
  } finally {
    loading.value = false
  }
}

async function stopMonitoring() {
  loading.value = true
  try {
    await StatusStopMonitoring()
    monitoring.value = false
  } catch (error) {
    handleError(error, 'Stop Monitoring')
  } finally {
    loading.value = false
  }
}

async function fetchMetrics() {
  const timeoutPromise = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('Metrics fetch timeout (10s)')), 10000)
  })

  const data = await Promise.race([StatusGetMetrics(), timeoutPromise])
  if (!data) throw new Error('No metrics data received')
  metrics.value = data
}

async function refreshMetrics() {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await fetchMetrics()
  } catch (error) {
    handleError(error, 'Fetch Metrics')
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  EventsOn('status:update', (data) => {
    metrics.value = data
  })
  startMonitoring()
})

onUnmounted(() => {
  EventsOff('status:update')
  if (monitoring.value) {
    StatusStopMonitoring()
  }
})
</script>

<template>
  <div class="status-tab">
    <div class="status-toolbar">
      <PageHeader
        title="System Status"
        subtitle="Real-time system health monitoring"
      />
      <div class="status-toolbar__actions">
        <AppButton
          v-if="!monitoring"
          variant="primary"
          :loading="loading"
          @click="startMonitoring"
        >
          Start Monitoring
        </AppButton>
        <template v-else>
          <AppButton
            variant="danger"
            :loading="loading"
            @click="stopMonitoring"
          >
            Stop Monitoring
          </AppButton>
          <button
            type="button"
            class="status-toolbar__refresh"
            :disabled="refreshing"
            aria-label="Refresh metrics"
            @click="refreshMetrics"
          >
            <i
              class="pi pi-refresh status-toolbar__refresh-icon"
              :class="{ 'status-toolbar__refresh-icon--spin': refreshing }"
              aria-hidden="true"
            />
          </button>
        </template>
      </div>
    </div>

    <LoadingPanel v-if="loading && !metrics" message="Loading metrics" />

    <LoadingPanel v-else-if="monitoring && !metrics" message="Collecting metrics" />

    <div v-else-if="metrics" class="status-content">
      <StatusHealthSummary :score="healthScore" :status="healthStatus" />

      <div class="status-group">
        <ResourceSection
          v-if="metrics.cpu"
          icon="pi-desktop"
          title="CPU"
          :value="formatPercent(metrics.cpu.totalPercent || 0)"
          :progress="metrics.cpu.totalPercent || 0"
          :progress-color="getUsageColor(metrics.cpu.totalPercent || 0)"
        >
          <template #details>
            <div v-if="metrics.cpu.temperature" class="detail-row">
              <span class="detail-row__label">Temperature</span>
              <span class="detail-row__value">{{ formatTemp(metrics.cpu.temperature) }}</span>
            </div>
            <div v-if="metrics.cpu.cores" class="detail-row">
              <span class="detail-row__label">Cores</span>
              <span class="detail-row__value">{{ metrics.cpu.cores }}</span>
            </div>
          </template>
        </ResourceSection>

        <ResourceSection
          v-if="metrics.memory"
          icon="pi-database"
          title="Memory"
          :value="formatPercent(metrics.memory.percent || 0)"
          :progress="metrics.memory.percent || 0"
          :progress-color="getUsageColor(metrics.memory.percent || 0)"
        >
          <template #details>
            <div class="detail-row">
              <span class="detail-row__label">Used</span>
              <span class="detail-row__value">{{ formatBytes(metrics.memory.used || 0) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-row__label">Total</span>
              <span class="detail-row__value">{{ formatBytes(metrics.memory.total || 0) }}</span>
            </div>
          </template>
        </ResourceSection>

        <ResourceSection
          v-if="metrics.disk"
          icon="pi-server"
          title="Disk"
          :value="formatPercent(metrics.disk.percent || 0)"
          :progress="metrics.disk.percent || 0"
          :progress-color="getUsageColor(metrics.disk.percent || 0)"
        >
          <template #details>
            <div class="detail-row">
              <span class="detail-row__label">Free</span>
              <span class="detail-row__value">{{ formatBytes(metrics.disk.free || 0) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-row__label">Total</span>
              <span class="detail-row__value">{{ formatBytes(metrics.disk.total || 0) }}</span>
            </div>
          </template>
        </ResourceSection>

        <ResourceSection
          v-if="metrics.network"
          icon="pi-globe"
          title="Network"
        >
          <template #details>
            <div class="detail-row">
              <span class="detail-row__label">
                <i class="pi pi-arrow-down detail-row__icon" aria-hidden="true" />
                Download
              </span>
              <span class="detail-row__value">{{ formatRate(metrics.network.download || 0) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-row__label">
                <i class="pi pi-arrow-up detail-row__icon" aria-hidden="true" />
                Upload
              </span>
              <span class="detail-row__value">{{ formatRate(metrics.network.upload || 0) }}</span>
            </div>
          </template>
        </ResourceSection>

        <ResourceSection
          v-if="metrics.battery && metrics.battery.level !== undefined"
          icon="pi-mobile"
          title="Battery"
          :value="formatPercent(metrics.battery.level || 0)"
          :progress="metrics.battery.level || 0"
          :progress-color="getBatteryColor(metrics.battery.level || 0)"
        >
          <template #details>
            <div v-if="metrics.battery.status" class="detail-row">
              <span class="detail-row__label">Status</span>
              <span class="detail-row__value">{{ metrics.battery.status }}</span>
            </div>
            <div v-if="metrics.battery.health" class="detail-row">
              <span class="detail-row__label">Health</span>
              <span class="detail-row__value">{{ metrics.battery.health }}</span>
            </div>
          </template>
        </ResourceSection>

        <ResourceSection
          v-if="metrics.gpu"
          icon="pi-microchip"
          title="GPU"
          :value="gpuUsage !== null ? formatPercent(gpuUsage) : ''"
          :progress="gpuUsage"
          :progress-color="gpuUsage !== null ? getUsageColor(gpuUsage) : undefined"
          :unavailable="gpuUsage === null"
        >
          <template #details>
            <div v-if="metrics.gpu.temperature" class="detail-row">
              <span class="detail-row__label">Temperature</span>
              <span class="detail-row__value">{{ formatTemp(metrics.gpu.temperature) }}</span>
            </div>
            <div v-if="metrics.gpu.memory" class="detail-row">
              <span class="detail-row__label">Memory</span>
              <span class="detail-row__value">{{ formatBytes(metrics.gpu.memory) }}</span>
            </div>
          </template>
        </ResourceSection>
      </div>
    </div>

    <EmptyState
      v-else
      message="Click Start Monitoring to view real-time system metrics"
    />
  </div>
</template>

<style scoped>
.status-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

.status-toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
  padding-bottom: var(--space-3);
  margin-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.status-toolbar :deep(.page-header) {
  margin-bottom: 0;
  flex: 1;
}

.status-toolbar__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-shrink: 0;
  padding-top: var(--space-1);
}

.status-toolbar__refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition:
    background var(--transition-fast),
    border-color var(--transition-fast),
    color var(--transition-fast);
}

.status-toolbar__refresh:not(:disabled):hover {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
  border-color: var(--color-text-tertiary);
}

.status-toolbar__refresh:not(:disabled):active {
  opacity: 0.85;
}

.status-toolbar__refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.status-toolbar__refresh-icon {
  font-size: 13px;
  line-height: 1;
}

.status-toolbar__refresh-icon--spin {
  animation: refresh-spin 0.7s linear infinite;
}

@keyframes refresh-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (prefers-reduced-motion: reduce) {
  .status-toolbar__refresh-icon--spin {
    animation: none;
  }
}

.status-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.status-group {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.status-group > :deep(* + *) {
  border-top: 1px solid var(--color-border);
}

.status-group :deep(.detail-row__label) {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.status-group :deep(.detail-row__icon) {
  font-size: 10px;
  color: var(--color-text-secondary);
}
</style>
