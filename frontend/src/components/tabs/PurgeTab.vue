<script setup>
import PageHeader from '../shared/PageHeader.vue'
import EmptyState from '../shared/EmptyState.vue'
</script>

<template>
  <div class="purge-tab">
    <PageHeader title="Purge Artifacts" subtitle="Coming soon" />
    <EmptyState message="Purge is coming soon" />
  </div>
</template>

<style scoped>
.purge-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}
</style>
