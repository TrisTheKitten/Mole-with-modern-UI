<script setup>
const props = defineProps({
  name: {
    type: String,
    required: true,
  },
})

const iconMap = {
  clean: 'pi-eraser',
  purge: 'pi-folder-minus',
  installer: 'pi-box',
  history: 'pi-history',
  uninstall: 'pi-trash',
  optimize: 'pi-bolt',
  analyze: 'pi-chart-bar',
  status: 'pi-chart-line',
  touchid: 'pi-lock',
  about: 'pi-info-circle',
}
</script>

<template>
  <i :class="['pi', iconMap[name], 'nav-icon']" aria-hidden="true" />
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
