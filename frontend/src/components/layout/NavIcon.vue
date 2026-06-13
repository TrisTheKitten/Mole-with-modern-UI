<script setup>
import { NAV_ICON_CLASSES } from '../../config/navIcons'

defineProps({
  name: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <i :class="['pi', NAV_ICON_CLASSES[name], 'nav-icon']" aria-hidden="true" />
</template>

<style scoped>
.nav-icon {
  font-size: 14px;
  width: 16px;
  text-align: center;
  flex-shrink: 0;
  line-height: 1;
}
</style>
