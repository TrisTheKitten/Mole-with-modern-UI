<script setup>
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useInstallerStore } from '../../stores/installer'
import PageHeader from '../shared/PageHeader.vue'
import AppButton from '../shared/AppButton.vue'
import ActionBar from '../shared/ActionBar.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import EmptyState from '../shared/EmptyState.vue'
import InfoRow from '../shared/InfoRow.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import MessageBanner from '../shared/MessageBanner.vue'

const MIN_BAR_PERCENT = 6

const store = useInstallerStore()
const { files, scanning, removing, progressMessage, result, error } = storeToRefs(store)
const showConfirmDialog = ref(false)
const selectedPaths = ref([])
const activePath = ref('')

const displayFiles = computed(() => [...files.value].sort((a, b) => b.size - a.size))
const maxSize = computed(() => files.value.reduce((max, file) => Math.max(max, file.size), 0))
const selectedFiles = computed(() => files.value.filter((file) => selectedPaths.value.includes(file.path)))
const selectedSize = computed(() => selectedFiles.value.reduce((sum, file) => sum + file.size, 0))
const totalSize = computed(() => files.value.reduce((sum, file) => sum + file.size, 0))
const activeFile = computed(
  () => files.value.find((file) => file.path === activePath.value) || displayFiles.value[0] || null
)
const allSelected = computed(() => files.value.length > 0 && selectedPaths.value.length === files.value.length)
const partiallySelected = computed(() => selectedPaths.value.length > 0 && !allSelected.value)
const countLabel = computed(() => `${files.value.length} ${files.value.length === 1 ? 'installer' : 'installers'}`)
const selectionSummary = computed(() =>
  selectedFiles.value.length === 0
    ? 'No files selected'
    : `${selectedFiles.value.length} selected · ${formatBytes(selectedSize.value)}`
)

onMounted(async () => {
  store.setupEventListeners()
  await store.scan()
  activePath.value = displayFiles.value[0]?.path || ''
})

function isSelected(file) {
  return selectedPaths.value.includes(file.path)
}

function toggleFile(file) {
  activePath.value = file.path
  if (isSelected(file)) {
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

function pruneSelectionToFiles() {
  selectedPaths.value = selectedPaths.value.filter((path) => files.value.some((file) => file.path === path))
  activePath.value = files.value.some((file) => file.path === activePath.value)
    ? activePath.value
    : displayFiles.value[0]?.path || ''
}

async function refreshInstallers() {
  await store.scan()
  pruneSelectionToFiles()
}

function requestRemove() {
  if (selectedFiles.value.length === 0) return
  showConfirmDialog.value = true
}

async function confirmRemove() {
  const paths = selectedFiles.value.map((file) => file.path)
  await store.remove(paths)
  pruneSelectionToFiles()
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

function fileExt(path) {
  const name = fileName(path)
  const dot = name.lastIndexOf('.')
  if (dot <= 0 || dot === name.length - 1) return 'FILE'
  return name.slice(dot + 1).toUpperCase().slice(0, 4)
}

function sizePercent(file) {
  if (!maxSize.value || !file.size) return 0
  return Math.max(MIN_BAR_PERCENT, (file.size / maxSize.value) * 100)
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
    <PageHeader title="Installer Cleanup" subtitle="Remove leftover installer downloads to reclaim space" />

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
      <EmptyState v-if="files.length === 0" message="No installer files found in your Downloads">
        <AppButton variant="secondary" @click="refreshInstallers">
          <i class="pi pi-refresh" aria-hidden="true" /> Rescan
        </AppButton>
      </EmptyState>

      <template v-else>
        <section class="overview">
          <div class="overview__metric">
            <span class="overview__value">{{ formatBytes(totalSize) }}</span>
            <span class="overview__label">Reclaimable · {{ countLabel }}</span>
          </div>
          <button type="button" class="ghost-button" @click="refreshInstallers">
            <i class="pi pi-refresh ghost-button__icon" aria-hidden="true" />
            <span>Refresh</span>
          </button>
        </section>

        <div class="selection-bar">
          <label class="select-all">
            <input
              type="checkbox"
              class="select-all__input"
              :checked="allSelected"
              :indeterminate.prop="partiallySelected"
              :aria-checked="partiallySelected ? 'mixed' : String(allSelected)"
              @change="toggleSelectAll"
            >
            <span>Select all</span>
          </label>
          <span class="selection-bar__summary" :class="{ 'selection-bar__summary--active': selectedFiles.length }">
            {{ selectionSummary }}
          </span>
        </div>

        <div class="installer-layout">
          <div class="installer-list" role="list">
            <article
              v-for="(file, index) in displayFiles"
              :key="file.path"
              class="installer-row"
              :class="{
                'installer-row--selected': isSelected(file),
                'installer-row--active': activeFile && activeFile.path === file.path
              }"
              :style="{ animationDelay: `${index * 40}ms` }"
              role="checkbox"
              tabindex="0"
              :aria-checked="isSelected(file)"
              :aria-label="fileName(file.path)"
              @click="toggleFile(file)"
              @keydown="handleRowKeydown($event, file)"
            >
              <span class="installer-row__check" aria-hidden="true">
                <svg v-if="isSelected(file)" class="installer-row__tick" viewBox="0 0 12 12" fill="none">
                  <path
                    d="M2 6L5 9L10 3"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </span>

              <span class="file-chip" aria-hidden="true">{{ fileExt(file.path) }}</span>

              <span class="installer-row__main">
                <span class="installer-row__title" :title="fileName(file.path)">{{ fileName(file.path) }}</span>
                <span class="installer-row__path" :title="file.path">{{ file.path }}</span>
                <span class="size-bar" aria-hidden="true">
                  <span class="size-bar__fill" :style="{ width: `${sizePercent(file)}%` }" />
                </span>
              </span>

              <span class="installer-row__size">{{ formatBytes(file.size) }}</span>

              <button
                type="button"
                class="installer-row__info"
                aria-label="View details"
                @click.stop="showDetails(file)"
              >
                <i class="pi pi-info-circle" aria-hidden="true" />
              </button>
            </article>
          </div>

          <aside v-if="activeFile" class="details-panel" aria-label="Installer details">
            <div class="details-panel__head">
              <span class="file-chip file-chip--lg" aria-hidden="true">{{ fileExt(activeFile.path) }}</span>
              <div class="details-panel__head-text">
                <span class="details-panel__name" :title="fileName(activeFile.path)">{{ fileName(activeFile.path) }}</span>
                <span class="details-panel__size">
                  {{ formatBytes(activeFile.size) }}
                  <span class="details-panel__size-label">on disk</span>
                </span>
              </div>
            </div>
            <div class="details-panel__rows">
              <InfoRow label="Source" :value="activeFile.source" />
              <InfoRow label="Modified" :value="formatDate(activeFile.lastModified)" />
              <InfoRow label="Path" :value="activeFile.path" mono />
            </div>
          </aside>
        </div>

        <ActionBar :summary="selectionSummary">
          <AppButton variant="danger" :disabled="selectedFiles.length === 0" @click="requestRemove">
            <i class="pi pi-trash" aria-hidden="true" /> Remove Files
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

.installer-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.overview {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  margin-bottom: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.overview__metric {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.overview__value {
  font-size: var(--font-size-metric);
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--color-accent);
  font-variant-numeric: tabular-nums;
}

.overview__label {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
}

.ghost-button {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  min-height: 32px;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  background: var(--color-bg-app);
  color: var(--color-text-primary);
  font-family: inherit;
  font-size: var(--font-size-body);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.ghost-button:hover {
  background: var(--color-bg-elevated);
}

.ghost-button:active {
  opacity: 0.85;
}

.ghost-button__icon {
  font-size: 13px;
  color: var(--color-accent);
}

.selection-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding-bottom: var(--space-3);
  margin-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
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

.select-all__input {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent);
  cursor: pointer;
}

.selection-bar__summary {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}

.selection-bar__summary--active {
  color: var(--color-accent);
  font-weight: 600;
}

.installer-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 288px;
  gap: var(--space-3);
  flex: 1;
  min-height: 0;
}

.installer-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  min-height: 0;
  overflow-y: auto;
  padding: 2px 2px var(--space-2);
}

.installer-row {
  display: grid;
  grid-template-columns: auto auto minmax(0, 1fr) auto auto;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  animation: row-reveal 0.32s ease both;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.installer-row:hover {
  background: var(--color-bg-elevated);
  border-color: var(--color-border-strong);
}

.installer-row--selected,
.installer-row--selected:hover {
  background: var(--color-accent-subtle);
  border-color: var(--color-accent-border);
}

.installer-row--active {
  box-shadow: inset 3px 0 0 var(--color-accent);
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
  color: var(--color-selection-text);
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.installer-row--selected .installer-row__check {
  background: var(--color-accent);
  border-color: var(--color-accent);
}

.installer-row__tick {
  width: 12px;
  height: 12px;
}

.file-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 38px;
  height: 38px;
  padding: 0 var(--space-1);
  background: var(--color-bg-inset);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-accent);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  font-variant-numeric: tabular-nums;
}

.file-chip--lg {
  width: 48px;
  height: 48px;
  font-size: 12px;
}

.installer-row__main {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.installer-row__title {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.installer-row__path {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.size-bar {
  display: block;
  width: 100%;
  max-width: 220px;
  height: 4px;
  margin-top: 3px;
  background: var(--color-bg-inset);
  border-radius: 999px;
  overflow: hidden;
}

.size-bar__fill {
  display: block;
  height: 100%;
  min-width: 4px;
  background: var(--color-accent);
  border-radius: 999px;
  opacity: 0.55;
  transition: width var(--transition-fast), opacity var(--transition-fast);
}

.installer-row--selected .size-bar__fill {
  opacity: 1;
}

.installer-row__size {
  flex-shrink: 0;
  min-width: 72px;
  text-align: right;
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  font-variant-numeric: tabular-nums;
}

.installer-row__info {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  font-size: 15px;
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.installer-row__info:hover {
  background: var(--color-bg-app);
  color: var(--color-accent);
}

.details-panel {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  align-self: start;
  position: sticky;
  top: 0;
  padding: var(--space-4);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.details-panel__head {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.details-panel__head-text {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  min-width: 0;
}

.details-panel__name {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.35;
  word-break: break-word;
}

.details-panel__size {
  display: flex;
  align-items: baseline;
  gap: var(--space-2);
  font-size: var(--font-size-metric);
  font-weight: 600;
  color: var(--color-accent);
  font-variant-numeric: tabular-nums;
  line-height: 1.1;
}

.details-panel__size-label {
  font-size: var(--font-size-caption);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.details-panel__rows {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

@keyframes row-reveal {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .installer-row {
    animation: none;
  }
}

@media (max-width: 900px) {
  .installer-layout {
    grid-template-columns: 1fr;
  }

  .details-panel {
    position: static;
  }
}
</style>
