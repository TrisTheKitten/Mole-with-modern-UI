<script setup>
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { ref } from 'vue'
import { useInstallerStore } from '../../stores/installer'
import PageHeader from '../shared/PageHeader.vue'
import AppButton from '../shared/AppButton.vue'
import ActionBar from '../shared/ActionBar.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import EmptyState from '../shared/EmptyState.vue'
import InfoRow from '../shared/InfoRow.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import MessageBanner from '../shared/MessageBanner.vue'

const store = useInstallerStore()
const { files, scanning, removing, progressMessage, result, error } = storeToRefs(store)
const showConfirmDialog = ref(false)
const selectedPaths = ref([])
const activePath = ref('')

const selectedFiles = computed(() => files.value.filter((file) => selectedPaths.value.includes(file.path)))
const selectedSize = computed(() => selectedFiles.value.reduce((sum, file) => sum + file.size, 0))
const totalSize = computed(() => files.value.reduce((sum, file) => sum + file.size, 0))
const activeFile = computed(() => files.value.find((file) => file.path === activePath.value) || files.value[0] || null)
const allSelected = computed(() => files.value.length > 0 && selectedPaths.value.length === files.value.length)
const partiallySelected = computed(() => selectedPaths.value.length > 0 && !allSelected.value)
const selectionSummary = computed(() => `${selectedFiles.value.length} selected · ${formatBytes(selectedSize.value)}`)

onMounted(async () => {
  store.setupEventListeners()
  await store.scan()
  activePath.value = files.value[0]?.path || ''
})

function toggleFile(file) {
  activePath.value = file.path
  if (selectedPaths.value.includes(file.path)) {
    selectedPaths.value = selectedPaths.value.filter((path) => path !== file.path)
    return
  }
  selectedPaths.value = [...selectedPaths.value, file.path]
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedPaths.value = []
    return
  }
  selectedPaths.value = files.value.map((file) => file.path)
}

function showDetails(file) {
  activePath.value = file.path
}

function handleRowKeydown(event, file) {
  if (event.key === ' ' || event.key === 'Enter') {
    event.preventDefault()
    toggleFile(file)
  }
}

async function refreshInstallers() {
  await store.scan()
  selectedPaths.value = selectedPaths.value.filter((path) =>
    files.value.some((file) => file.path === path)
  )
  activePath.value = files.value.some((file) => file.path === activePath.value)
    ? activePath.value
    : files.value[0]?.path || ''
}

function requestRemove() {
  if (selectedFiles.value.length === 0) return
  showConfirmDialog.value = true
}

async function confirmRemove() {
  const paths = selectedFiles.value.map((file) => file.path)
  await store.remove(paths)
  selectedPaths.value = selectedPaths.value.filter((path) =>
    files.value.some((file) => file.path === path)
  )
  activePath.value = files.value.some((file) => file.path === activePath.value)
    ? activePath.value
    : files.value[0]?.path || ''
}

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)))
  return `${(bytes / Math.pow(1024, index)).toFixed(2)} ${units[index]}`
}

function fileName(path) {
  return path.split('/').filter(Boolean).pop() || path
}

function formatDate(value) {
  if (!value) return 'Unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Unknown'
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}
</script>

<template>
  <div class="installer-tab">
    <PageHeader title="Installer Cleanup" subtitle="Remove installer downloads" />

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Remove Installers"
      :message="`Remove ${selectedFiles.length} files and free ${formatBytes(selectedSize)}?`"
      confirm-text="Remove Files"
      cancel-text="Cancel"
      destructive
      @confirm="confirmRemove"
    />

    <MessageBanner v-if="error" type="error" :message="error" />
    <LoadingPanel v-if="scanning || removing" :message="progressMessage || 'Scanning'" />

    <div v-else class="installer-content">
      <div class="summary-row">
        <span>{{ files.length }} files · {{ formatBytes(totalSize) }}</span>
        <div class="toolbar__actions">
          <AppButton variant="secondary" @click="refreshInstallers">Refresh</AppButton>
        </div>
      </div>

      <EmptyState v-if="files.length === 0" message="No installer files found" />

      <template v-else>
        <div class="selection-bar">
          <label class="select-all">
            <input
              type="checkbox"
              :checked="allSelected"
              :indeterminate.prop="partiallySelected"
              :aria-checked="partiallySelected ? 'mixed' : String(allSelected)"
              @change="toggleSelectAll"
            >
            <span>Select All</span>
          </label>
          <span>{{ selectionSummary }}</span>
        </div>

        <div class="installer-layout">
          <div class="installer-list">
            <div
              v-for="file in files"
              :key="file.path"
              class="installer-row"
              :class="{
                'installer-row--selected': selectedPaths.includes(file.path),
                'installer-row--active': activePath === file.path
              }"
              role="checkbox"
              tabindex="0"
              :aria-checked="selectedPaths.includes(file.path)"
              :aria-label="fileName(file.path)"
              @click="toggleFile(file)"
              @keydown="handleRowKeydown($event, file)"
            >
              <span class="installer-row__check" aria-hidden="true">
                <span v-if="selectedPaths.includes(file.path)" class="installer-row__check-mark"></span>
              </span>
              <span class="installer-row__main">
                <span class="installer-row__title">{{ fileName(file.path) }}</span>
                <span class="installer-row__path">{{ file.path }}</span>
              </span>
              <span class="installer-row__meta">{{ formatBytes(file.size) }}</span>
              <AppButton variant="ghost" @click.stop="showDetails(file)">Details</AppButton>
            </div>
          </div>

          <aside v-if="activeFile" class="details-panel" aria-label="Installer details">
            <h2 class="details-panel__title">Details</h2>
            <InfoRow label="Name" :value="fileName(activeFile.path)" />
            <InfoRow label="Source" :value="activeFile.source" />
            <InfoRow label="Size" :value="formatBytes(activeFile.size)" />
            <InfoRow label="Modified" :value="formatDate(activeFile.lastModified)" />
            <InfoRow label="Path" :value="activeFile.path" mono />
          </aside>
        </div>

        <ActionBar :summary="selectionSummary">
          <AppButton variant="danger" :disabled="selectedFiles.length === 0" @click="requestRemove">
            Remove Files
          </AppButton>
        </ActionBar>
      </template>

      <MessageBanner
        v-if="result && result.errors && result.errors.length"
        type="error"
        :message="result.errors.join('; ')"
      />
    </div>
  </div>
</template>

<style scoped>
.installer-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.installer-content,
.installer-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.installer-content {
  flex: 1;
  min-height: 0;
}

.summary-row,
.toolbar__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.summary-row {
  justify-content: space-between;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-2);
}

.selection-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}

.select-all {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text-primary);
  font-size: var(--font-size-body);
  font-weight: 500;
  cursor: pointer;
}

.select-all input {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent);
}

.installer-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 280px;
  gap: var(--space-3);
  flex: 1;
  min-height: 0;
}

.installer-list {
  min-height: 0;
  overflow-y: auto;
  padding-bottom: var(--space-2);
}

.installer-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  text-align: left;
  cursor: pointer;
}

.installer-row--selected {
  background: var(--color-accent-subtle);
  border-color: var(--color-accent-border);
}

.installer-row--active {
  border-color: var(--color-border-strong);
}

.installer-row__check {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  background: var(--color-bg-app);
  border: 1.5px solid var(--color-border-strong);
  border-radius: 4px;
}

.installer-row--selected .installer-row__check {
  background: var(--color-accent);
  border-color: var(--color-accent);
}

.installer-row__check-mark {
  width: 8px;
  height: 8px;
  background: var(--color-selection-text);
  border-radius: 2px;
}

.installer-row__main {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.installer-row__title {
  font-weight: 600;
}

.installer-row__path,
.installer-row__meta {
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}

.installer-row__path {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.installer-row__meta {
  flex-shrink: 0;
  font-variant-numeric: tabular-nums;
}

.details-panel {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  align-self: start;
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.details-panel__title {
  margin: 0 0 var(--space-1);
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

@media (max-width: 900px) {
  .installer-layout {
    grid-template-columns: 1fr;
  }
}
</style>
