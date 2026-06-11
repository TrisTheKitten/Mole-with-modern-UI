<script setup>
import PageHeader from '../shared/PageHeader.vue'
import EmptyState from '../shared/EmptyState.vue'
</script>

<template>
  <div class="touchid-tab">
    <PageHeader title="Touch ID" subtitle="Coming soon" />
    <EmptyState message="Touch ID is coming soon" />
  </div>
</template>

<style scoped>
.touchid-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}
</style>
