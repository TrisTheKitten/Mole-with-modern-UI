<script setup>
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useHistoryStore } from '../../stores/history'
import PageHeader from '../shared/PageHeader.vue'
import AppButton from '../shared/AppButton.vue'
import EmptyState from '../shared/EmptyState.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import MessageBanner from '../shared/MessageBanner.vue'

const store = useHistoryStore()
const { sessions, deletions, logs, loading, error } = storeToRefs(store)

const NON_BREAKING_SPACE = '\u00a0'
const EMPTY_SIZE = 'unknown'
const ACTION_LABELS = [
  ['removed', 'removed'],
  ['trashed', 'trashed'],
  ['skipped', 'skipped'],
  ['failed', 'failed'],
  ['rebuilt', 'rebuilt'],
  ['other', 'other'],
]

onMounted(() => {
  store.load()
})

function actionSummary(actions) {
  const parts = ACTION_LABELS
    .filter(([key]) => actions?.[key] > 0)
    .map(([key, label]) => `${label} ${actions[key]}`)

  return parts.length > 0 ? parts.join(', ') : 'no file actions'
}

function formatSizeKb(sizeKb) {
  if (sizeKb === null || sizeKb === undefined) return EMPTY_SIZE

  return `${sizeKb}${NON_BREAKING_SPACE}KB`
}

function hasHistory() {
  return sessions.value.length > 0 || deletions.value.length > 0
}
</script>

<template>
  <div class="history-tab">
    <div class="history-toolbar">
      <PageHeader title="History" subtitle="Review past Mole operations" />
      <AppButton variant="secondary" :loading="loading" @click="store.load()">Refresh</AppButton>
    </div>

    <MessageBanner v-if="error" type="error" :message="error" />
    <LoadingPanel v-if="loading" message="Loading history" />
    <EmptyState v-else-if="!hasHistory()" message="No operation history found" />

    <div v-else class="history-content">
      <section class="history-section">
        <h2 class="history-section__title">Recent Sessions</h2>
        <EmptyState v-if="sessions.length === 0" message="No sessions found" />
        <div v-else class="history-list">
          <div v-for="session in sessions" :key="`${session.command}-${session.startedAt}`" class="history-row">
            <div class="history-row__main">
              <h3 class="history-row__title">{{ session.command }}</h3>
              <p class="history-row__meta">{{ session.startedAt }} to {{ session.endedAt || 'not ended' }}</p>
              <p class="history-row__meta">{{ actionSummary(session.actions) }}</p>
            </div>
            <span class="history-row__meta">{{ session.items }} items · {{ session.size }}</span>
          </div>
        </div>
      </section>

      <section class="history-section">
        <h2 class="history-section__title">Deletion Audit</h2>
        <EmptyState v-if="deletions.length === 0" message="No deletions found" />
        <div v-else class="history-list">
          <div v-for="deletion in deletions" :key="`${deletion.timestamp}-${deletion.path}`" class="history-row">
            <div class="history-row__main">
              <h3 class="history-row__title">{{ deletion.path }}</h3>
              <p class="history-row__meta">{{ deletion.timestamp }}</p>
            </div>
            <span class="history-row__meta">
              {{ deletion.mode }} · {{ deletion.status }} · {{ formatSizeKb(deletion.sizeKb) }}
            </span>
          </div>
        </div>
      </section>

      <div class="history-logs">
        <p class="history-row__meta">Operations log: {{ logs.operations }}</p>
        <p class="history-row__meta">Deletions log: {{ logs.deletions }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.history-toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-3);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.history-content,
.history-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.history-section__title {
  margin: 0;
  font-size: var(--font-size-body);
  color: var(--color-text-primary);
}

.history-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.history-row__main {
  min-width: 0;
}

.history-row__title {
  margin: 0 0 var(--space-1);
  font-size: var(--font-size-body);
  color: var(--color-text-primary);
  overflow-wrap: anywhere;
}

.history-row__meta {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
  overflow-wrap: anywhere;
}

.history-logs {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}
</style>
