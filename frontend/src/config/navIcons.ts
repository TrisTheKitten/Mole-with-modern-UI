export const NAV_ICON_CLASSES = {
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
} as const

export type NavIconName = keyof typeof NAV_ICON_CLASSES
