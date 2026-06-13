<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import {
  AnalyzeScanDirectory,
  AnalyzeGetLargeFiles,
  AnalyzeDeletePath,
  AnalyzeOpenInFinder,
  AnalyzePickDirectory,
  AnalyzeListExternalVolumes,
  AnalyzeScanExternalVolume
} from '../../../wailsjs/go/main/App'
import { validatePath } from '../../utils/validation'
import { handleError } from '../../utils/errorHandler'
import PageHeader from '../shared/PageHeader.vue'
import TextField from '../shared/TextField.vue'
import AppButton from '../shared/AppButton.vue'
import CheckboxRow from '../shared/CheckboxRow.vue'
import ActionBar from '../shared/ActionBar.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import EmptyState from '../shared/EmptyState.vue'
import InfoRow from '../shared/InfoRow.vue'

const LARGE_FILE_LIMIT = 20

const scanPath = ref('')
const scanning = ref(false)
const scanProgress = ref('')
const scanResult = ref(null)
const largeFiles = ref([])
const loading = ref(false)
const showDeleteDialog = ref(false)
const deleteTargets = ref([])
const externalVolumes = ref([])
const selectedExternalVolume = ref('')
const showExternalEmpty = ref(false)

const selectedFiles = computed(() => largeFiles.value.filter((file) => file.selected))

const deleteDialogTitle = computed(() =>
  deleteTargets.value.length === 1 ? 'Delete Item' : 'Delete Items'
)

const deleteDialogMessage = computed(() => {
  if (deleteTargets.value.length === 1) {
    return 'Delete this item permanently?'
  }
  return `Delete ${deleteTargets.value.length} items permanently?`
})

let unsubscribeProgress = null

onMounted(() => {
  scanPath.value = '/Users'
  unsubscribeProgress = EventsOn('analyze:progress', (data) => {
    scanProgress.value = data.message || 'Scanning'
  })
  loadExternalVolumes()
})

onBeforeUnmount(() => {
  if (unsubscribeProgress) {
    EventsOff('analyze:progress')
  }
})

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function fetchLargeFiles(path) {
  const files = await AnalyzeGetLargeFiles(path, LARGE_FILE_LIMIT)
  return (files || []).map((file) => ({ ...file, selected: false }))
}

async function applyScanResult(path, result) {
  if (!result) {
    throw new Error('Scan returned no results')
  }

  scanPath.value = path
  scanResult.value = result
  largeFiles.value = await fetchLargeFiles(path)
}

async function runScan(fetchResult, path, errorContext) {
  try {
    const result = await fetchResult()
    await applyScanResult(path, result)
    return true
  } catch (error) {
    handleError(error, errorContext)
    return false
  }
}

async function scanDirectoryAndLoadFiles() {
  return runScan(
    () => AnalyzeScanDirectory(scanPath.value),
    scanPath.value,
    'Disk Analysis',
  )
}

async function scan() {
  const validation = validatePath(scanPath.value, true)
  if (!validation.valid) {
    handleError(new Error(validation.error), 'Path Validation')
    return
  }

  scanning.value = true
  scanProgress.value = 'Starting scan'
  scanResult.value = null
  largeFiles.value = []

  const success = await scanDirectoryAndLoadFiles()
  if (!success) {
    scanResult.value = null
  }

  scanProgress.value = ''
  scanning.value = false
}

async function loadExternalVolumes() {
  try {
    externalVolumes.value = await AnalyzeListExternalVolumes()
    if (externalVolumes.value.length > 0 && !selectedExternalVolume.value) {
      selectedExternalVolume.value = externalVolumes.value[0].path
    }
  } catch (error) {
    handleError(error, 'External Volumes')
  }
}

async function scanExternalVolume() {
  await loadExternalVolumes()
  if (externalVolumes.value.length === 0) {
    showExternalEmpty.value = true
    return
  }

  const target = selectedExternalVolume.value || externalVolumes.value[0].path
  scanning.value = true
  scanProgress.value = 'Scanning external volume'

  const success = await runScan(
    () => AnalyzeScanExternalVolume(target),
    target,
    'External Volume',
  )
  if (success) {
    showExternalEmpty.value = false
  }

  scanning.value = false
}

function requestDelete(paths) {
  const targets = Array.isArray(paths) ? paths : [paths]

  for (const path of targets) {
    const validation = validatePath(path)
    if (!validation.valid) {
      handleError(new Error(validation.error), 'Path Validation')
      return
    }
  }

  deleteTargets.value = targets
  showDeleteDialog.value = true
}

function requestDeleteSelected() {
  if (selectedFiles.value.length === 0) {
    handleError(new Error('Select at least one file'), 'Delete')
    return
  }
  requestDelete(selectedFiles.value.map((file) => file.path))
}

async function confirmDelete() {
  loading.value = true
  const targets = [...deleteTargets.value]
  const targetSet = new Set(targets)

  try {
    for (const path of targets) {
      await AnalyzeDeletePath(path)
    }

    largeFiles.value = largeFiles.value.filter((file) => !targetSet.has(file.path))
    await scanDirectoryAndLoadFiles()

    const message = targets.length === 1 ? 'Item deleted' : `${targets.length} items deleted`
    window.dispatchEvent(new CustomEvent('show-toast', {
      detail: { message, type: 'success' },
    }))
  } catch (error) {
    handleError(error, 'Delete')
    await scanDirectoryAndLoadFiles()
  } finally {
    loading.value = false
    deleteTargets.value = []
  }
}

function toggleFile(file) {
  file.selected = !file.selected
}

function selectAll() {
  largeFiles.value.forEach((file) => {
    file.selected = true
  })
}

function deselectAll() {
  largeFiles.value.forEach((file) => {
    file.selected = false
  })
}

async function openInFinder(path) {
  try {
    await AnalyzeOpenInFinder(path)
  } catch (error) {
    handleError(error, 'Open in Finder')
  }
}

async function browseDirectory() {
  try {
    const selectedPath = await AnalyzePickDirectory(scanPath.value)
    if (selectedPath) {
      scanPath.value = selectedPath
    }
  } catch (error) {
    handleError(error, 'Choose Folder')
  }
}
</script>

<template>
  <div class="analyze-tab">
    <PageHeader
      title="Disk Space Analyzer"
      subtitle="Visualize and analyze disk usage"
    />

    <ConfirmDialog
      v-model:show="showDeleteDialog"
      :title="deleteDialogTitle"
      :message="deleteDialogMessage"
      :items="deleteTargets"
      confirm-text="Delete"
      cancel-text="Cancel"
      destructive
      @confirm="confirmDelete"
    />

    <div class="scan-controls">
      <TextField
        v-model="scanPath"
        placeholder="/Users"
        mono
        :disabled="scanning"
      />
      <AppButton variant="secondary" :disabled="scanning" @click="browseDirectory">
        Choose Folder
      </AppButton>
      <AppButton variant="primary" :loading="scanning" :disabled="scanning" @click="scan">
        Scan
      </AppButton>
    </div>

    <div class="external-controls">
      <select v-model="selectedExternalVolume" class="external-controls__select" :disabled="scanning">
        <option v-for="volume in externalVolumes" :key="volume.path" :value="volume.path">
          {{ volume.name }}
        </option>
      </select>
      <AppButton variant="secondary" :disabled="scanning" @click="scanExternalVolume">
        Scan Volume
      </AppButton>
    </div>

    <EmptyState v-if="showExternalEmpty" message="No external volumes are available" />

    <LoadingPanel v-if="scanning" :message="scanProgress || 'Scanning'" />

    <div v-else-if="scanResult" class="analyze-results">
      <div class="summary-grid">
        <div class="summary-card">
          <span class="summary-card__label">Total Size</span>
          <span class="summary-card__value">{{ formatBytes(scanResult.totalSize || 0) }}</span>
        </div>
        <div class="summary-card">
          <span class="summary-card__label">Total Items</span>
          <span class="summary-card__value">{{ (scanResult.totalItems || 0).toLocaleString() }}</span>
        </div>
        <div class="summary-card summary-card--wide">
          <InfoRow label="Scanned Path" :value="scanPath" mono />
        </div>
      </div>

      <div v-if="largeFiles.length > 0" class="files-section">
        <div class="files-section__header">
          <h2 class="files-section__title">Largest Files</h2>
          <div class="files-section__controls">
            <AppButton variant="ghost" @click="selectAll">Select All</AppButton>
            <AppButton variant="ghost" @click="deselectAll">Deselect All</AppButton>
          </div>
        </div>

        <div class="files-list">
          <CheckboxRow
            v-for="file in largeFiles"
            :key="file.path"
            :title="file.name"
            :description="file.path"
            :size="formatBytes(file.size)"
            :checked="file.selected"
            @toggle="toggleFile(file)"
          >
            <template #trailing>
              <AppButton variant="ghost" @click="openInFinder(file.path)">Open</AppButton>
              <AppButton variant="ghost" :disabled="loading" @click="requestDelete(file.path)">Delete</AppButton>
            </template>
          </CheckboxRow>
        </div>

        <ActionBar :summary="`${selectedFiles.length} selected`">
          <AppButton
            variant="danger"
            :disabled="selectedFiles.length === 0 || loading"
            @click="requestDeleteSelected"
          >
            Delete Selected
          </AppButton>
        </ActionBar>
      </div>

      <EmptyState v-else message="No large files found in this directory" />
    </div>

    <EmptyState
      v-else
      message="Enter a path, choose a folder, or click Scan to analyze disk usage"
    />
  </div>
</template>

<style scoped>
.analyze-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.scan-controls {
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.external-controls {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
}

.external-controls__select {
  min-height: 32px;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  padding: 0 var(--space-2);
}

.analyze-results {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: var(--space-4);
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-2);
}

.summary-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.summary-card--wide {
  grid-column: 1 / -1;
  padding: var(--space-2);
}

.summary-card__label {
  font-size: var(--font-size-caption);
  font-weight: 500;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.summary-card__value {
  font-size: var(--font-size-metric);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-text-primary);
}

.files-section {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.files-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  margin-bottom: var(--space-3);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.files-section__title {
  margin: 0;
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.files-section__controls {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  flex-shrink: 0;
}

.files-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: var(--space-2);
}

.files-list :deep(.checkbox-row__trailing) {
  display: flex;
  gap: var(--space-1);
}
</style>
