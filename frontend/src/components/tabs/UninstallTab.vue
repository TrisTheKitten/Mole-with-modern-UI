<script setup>
import { ref, onMounted, computed } from 'vue'
import { UninstallScanApps, UninstallApps, UninstallGetRelatedFiles, UninstallPreview } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'
import { handleError } from '../../utils/errorHandler'
import PageHeader from '../shared/PageHeader.vue'
import CheckboxRow from '../shared/CheckboxRow.vue'
import ActionBar from '../shared/ActionBar.vue'
import AppButton from '../shared/AppButton.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import TextField from '../shared/TextField.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import UninstallTimeline from '../uninstall/UninstallTimeline.vue'

const BYTES_PER_MEGABYTE = 1024 * 1024
const BYTES_PER_GIGABYTE = 1024 * BYTES_PER_MEGABYTE
const LARGE_APP_BYTES = 1024 * BYTES_PER_MEGABYTE
const MEDIUM_APP_BYTES = 250 * BYTES_PER_MEGABYTE

const sourceFilters = [
  { value: 'all', label: 'All Sources' },
  { value: 'homebrew', label: 'Homebrew' },
  { value: 'manual', label: 'Manual' },
]

const ageFilters = [
  { value: 'all', label: 'All Ages' },
  { value: 'recent', label: 'Recent' },
  { value: 'month', label: '< 1 Month' },
  { value: 'older', label: 'Older' },
]

const sizeFilters = [
  { value: 'all', label: 'All Sizes' },
  { value: 'large', label: 'Large' },
  { value: 'medium', label: 'Medium' },
  { value: 'small', label: 'Small' },
]

const apps = ref([])
const loading = ref(false)
const phase = ref('idle')
const percent = ref(0)
const stepGroups = ref([])
const searchQuery = ref('')
const sourceFilter = ref('all')
const ageFilter = ref('all')
const sizeFilter = ref('all')
const showRelatedFiles = ref(false)
const relatedFiles = ref([])
const selectedApp = ref(null)
const showConfirmDialog = ref(false)
const summary = ref(null)
const dryRun = ref(false)
const preview = ref(null)

const filteredApps = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return apps.value.filter((app) =>
    matchesSearch(app, query) &&
    matchesSource(app) &&
    matchesAge(app) &&
    matchesSize(app)
  )
})

const selectedApps = computed(() => apps.value.filter((app) => app.selected))

const hasActiveFilters = computed(() =>
  searchQuery.value.trim() ||
  sourceFilter.value !== 'all' ||
  ageFilter.value !== 'all' ||
  sizeFilter.value !== 'all'
)

const relatedFilesMessage = computed(() => {
  if (!relatedFiles.value.length) return 'No related files found.'
  return relatedFiles.value.join('\n')
})

const STEP_LABELS = {
  locate: 'Find files',
  remove: 'Remove app',
  cleanup: 'Clean up',
}

function upsertStep(appName, stepId, status, detail) {
  if (!appName || !stepId) return
  let group = stepGroups.value.find((item) => item.app === appName)
  if (!group) {
    stepGroups.value.push({ app: appName, state: 'working', steps: [] })
    group = stepGroups.value[stepGroups.value.length - 1]
  }
  let step = group.steps.find((item) => item.id === stepId)
  if (!step) {
    group.steps.push({ id: stepId, label: STEP_LABELS[stepId] || stepId, detail: '', status: 'pending' })
    step = group.steps[group.steps.length - 1]
  }
  step.status = status
  step.detail = detail
  group.state = computeGroupState(group)
}

function computeGroupState(group) {
  if (group.steps.some((step) => step.status === 'failed')) return 'failed'
  const cleanup = group.steps.find((step) => step.id === 'cleanup')
  if (cleanup && cleanup.status === 'success') return 'done'
  return 'working'
}

function summarize({ appsRemoved, failedCount, spaceFreed }) {
  if (appsRemoved === 0) {
    return {
      tone: 'failed',
      title: 'Couldn’t uninstall',
      detail: 'Nothing was removed — check the steps above',
    }
  }
  const detail = [`Removed ${appsRemoved} ${appsRemoved === 1 ? 'app' : 'apps'}`]
  if (failedCount > 0) detail.push(`${failedCount} failed`)
  if (spaceFreed > 0) detail.push(`${formatSize(spaceFreed)} freed`)
  return {
    tone: failedCount > 0 ? 'partial' : 'success',
    title: failedCount > 0 ? 'Finished with issues' : 'All clean',
    detail: detail.join(' · '),
  }
}

onMounted(async () => {
  await scanApps(false)

  EventsOn('uninstall:progress', (data) => {
    if (phase.value !== 'running') return
    if (typeof data.percent === 'number') percent.value = data.percent
    if (data.step && data.status) {
      upsertStep(data.app, data.step, data.status, data.message)
    }
  })

  EventsOn('uninstall:complete', (data) => {
    if (phase.value !== 'running') return
    percent.value = 100
    summary.value = summarize({
      appsRemoved: data.appsRemoved || 0,
      failedCount: data.errors?.length || 0,
      spaceFreed: data.spaceFreed || 0,
    })
    phase.value = 'done'
  })
})

async function scanApps(forceRescan) {
  loading.value = true
  try {
    const data = await UninstallScanApps(forceRescan)
    apps.value = data.map((app) => ({ ...app, selected: false }))
  } catch (error) {
    handleError(error, 'Scan Apps')
  } finally {
    loading.value = false
  }
}

async function showRelated(app) {
  selectedApp.value = app
  try {
    const files = await UninstallGetRelatedFiles(getAppIdentifier(app))
    relatedFiles.value = files
    showRelatedFiles.value = true
  } catch (error) {
    handleError(error, 'Related Files')
  }
}

function requestUninstall() {
  if (selectedApps.value.length === 0) {
    handleError(new Error('Select at least one app'), 'Uninstall')
    return
  }
  if (dryRun.value) {
    uninstall()
    return
  }
  showConfirmDialog.value = true
}

async function uninstall() {
  const appIdentifiers = selectedApps.value.map((app) => getAppIdentifier(app))

  if (dryRun.value) {
    try {
      preview.value = await UninstallPreview(appIdentifiers)
    } catch (error) {
      handleError(error, 'Uninstall')
    }
    return
  }

  preview.value = null
  stepGroups.value = []
  summary.value = null
  percent.value = 0
  phase.value = 'running'

  try {
    await UninstallApps(appIdentifiers)
  } catch (error) {
    handleUninstallFailure(error)
  }
}

function handleUninstallFailure(error) {
  if (phase.value === 'done') return
  if (stepGroups.value.length === 0) {
    phase.value = 'idle'
    handleError(error, 'Uninstall')
    return
  }
  const removed = stepGroups.value.filter((group) => group.state === 'done').length
  const failed = stepGroups.value.filter((group) => group.state === 'failed').length
  percent.value = 100
  summary.value = summarize({ appsRemoved: removed, failedCount: failed, spaceFreed: 0 })
  phase.value = 'done'
}

function toggleApp(app) {
  app.selected = !app.selected
}

function formatSize(bytes) {
  if (!bytes) return '0 MB'
  if (bytes >= BYTES_PER_GIGABYTE) {
    return (bytes / BYTES_PER_GIGABYTE).toFixed(2) + ' GB'
  }
  return (bytes / BYTES_PER_MEGABYTE).toFixed(2) + ' MB'
}

async function handleResultDone() {
  phase.value = 'idle'
  stepGroups.value = []
  summary.value = null
  percent.value = 0
  await scanApps(true)
}

function getAppIdentifier(app) {
  return app.path || app.bundleId
}

function getAppKey(app) {
  return app.path || app.bundleId || app.name
}

function matchesSearch(app, query) {
  if (!query) return true
  return [app.name, app.path, app.bundleId]
    .filter(Boolean)
    .some((value) => value.toLowerCase().includes(query))
}

function matchesSource(app) {
  if (sourceFilter.value === 'all') return true
  const isHomebrew = Boolean(app.brewCask)
  return sourceFilter.value === 'homebrew' ? isHomebrew : !isHomebrew
}

function matchesAge(app) {
  if (ageFilter.value === 'all') return true
  if (ageFilter.value === 'recent') return app.age === 'Recent'
  if (ageFilter.value === 'month') return app.age === '< 1 month'
  return !['Recent', '< 1 month'].includes(app.age)
}

function matchesSize(app) {
  if (sizeFilter.value === 'all') return true
  if (sizeFilter.value === 'large') return app.size >= LARGE_APP_BYTES
  if (sizeFilter.value === 'medium') {
    return app.size >= MEDIUM_APP_BYTES && app.size < LARGE_APP_BYTES
  }
  return app.size < MEDIUM_APP_BYTES
}

function clearFilters() {
  searchQuery.value = ''
  sourceFilter.value = 'all'
  ageFilter.value = 'all'
  sizeFilter.value = 'all'
}
</script>

<template>
  <div class="uninstall-tab">
    <PageHeader
      title="Uninstall Applications"
      subtitle="Remove apps and their associated files"
    />

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Confirm Uninstall"
      :message="`Uninstall ${selectedApps.length} apps?`"
      confirm-text="Uninstall"
      cancel-text="Cancel"
      destructive
      @confirm="uninstall"
    />

    <ConfirmDialog
      v-model:show="showRelatedFiles"
      :title="`Related Files — ${selectedApp?.name || ''}`"
      :message="relatedFilesMessage"
      confirm-text="Close"
      @confirm="showRelatedFiles = false"
    />

    <LoadingPanel v-if="loading" message="Scanning apps" />

    <UninstallTimeline
      v-else-if="phase !== 'idle'"
      :groups="stepGroups"
      :percent="percent"
      :phase="phase"
      :summary="summary"
      @done="handleResultDone"
    />

    <div v-else class="uninstall-content">
      <div class="controls">
        <TextField
          v-model="searchQuery"
          placeholder="Search apps"
        />
        <label class="filter-field">
          <span>Source</span>
          <select v-model="sourceFilter">
            <option v-for="filter in sourceFilters" :key="filter.value" :value="filter.value">
              {{ filter.label }}
            </option>
          </select>
        </label>
        <label class="filter-field">
          <span>Age</span>
          <select v-model="ageFilter">
            <option v-for="filter in ageFilters" :key="filter.value" :value="filter.value">
              {{ filter.label }}
            </option>
          </select>
        </label>
        <label class="filter-field">
          <span>Size</span>
          <select v-model="sizeFilter">
            <option v-for="filter in sizeFilters" :key="filter.value" :value="filter.value">
              {{ filter.label }}
            </option>
          </select>
        </label>
        <label class="dry-run-toggle">
          <input v-model="dryRun" type="checkbox">
          <span>Dry Run</span>
        </label>
        <AppButton v-if="hasActiveFilters" variant="ghost" @click="clearFilters">Clear</AppButton>
        <AppButton variant="secondary" @click="scanApps(true)">Refresh</AppButton>
      </div>

      <p class="filter-summary">{{ filteredApps.length }} of {{ apps.length }} apps</p>

      <div v-if="preview" class="preview-panel">
        <h2 class="preview-panel__title">Preview</h2>
        <p v-if="preview.entries.length === 0" class="preview-panel__empty">No changes</p>
        <ul v-else class="preview-panel__list">
          <li v-for="entry in preview.entries" :key="`${entry.action}-${entry.path}`">{{ entry.path || entry.detail }}</li>
        </ul>
      </div>

      <div class="app-list">
        <p v-if="filteredApps.length === 0" class="empty-state">No apps match.</p>
        <CheckboxRow
          v-for="app in filteredApps"
          :key="getAppKey(app)"
          :title="app.name"
          :description="app.path"
          :size="`${formatSize(app.size)} · ${app.age}`"
          :checked="app.selected"
          @toggle="toggleApp(app)"
        >
          <template #trailing>
            <AppButton variant="ghost" @click="showRelated(app)">Files</AppButton>
          </template>
        </CheckboxRow>
      </div>

      <ActionBar :summary="`${selectedApps.length} selected`">
        <AppButton
          variant="primary"
          :disabled="selectedApps.length === 0"
          @click="requestUninstall"
        >
          Uninstall Selected
        </AppButton>
      </ActionBar>
    </div>
  </div>
</template>

<style scoped>
.uninstall-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.uninstall-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.controls {
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
  flex-wrap: wrap;
}

.app-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: var(--space-2);
}

.dry-run-toggle {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
  min-height: 32px;
}

.filter-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  min-width: 120px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}

.filter-field select {
  min-height: 32px;
  padding: var(--space-1) var(--space-2);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-family: inherit;
  font-size: var(--font-size-body);
}

.filter-summary,
.empty-state {
  margin: 0 0 var(--space-2);
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
  font-variant-numeric: tabular-nums;
}

.preview-panel {
  margin-bottom: var(--space-3);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
}

.preview-panel__title {
  margin: 0 0 var(--space-2);
  font-size: var(--font-size-body);
}

.preview-panel__empty,
.preview-panel__list {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}
</style>
