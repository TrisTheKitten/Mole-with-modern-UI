<script setup>
import PageHeader from '../shared/PageHeader.vue'
import EmptyState from '../shared/EmptyState.vue'
</script>

<template>
  <div class="history-tab">
    <PageHeader title="History" subtitle="Coming soon" />
    <EmptyState message="History is coming soon" />
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
</style>
