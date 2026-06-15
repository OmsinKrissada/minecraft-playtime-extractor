<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const props = withDefaults(
  defineProps<{
    label: string
    value: string | number
    accent?: 'green' | 'blue' | 'gold' | 'stone'
    tooltip?: string
  }>(),
  {
    accent: 'green',
  },
)

const accents: Record<string, string> = {
  green: 'after:bg-grass',
  blue: 'after:bg-diamond',
  gold: 'after:bg-gold',
  stone: 'after:bg-stone',
}

const tooltipOpen = ref(false)
const tooltipRef = ref<HTMLElement | null>(null)

function toggleTooltip() {
  tooltipOpen.value = !tooltipOpen.value
}

function handleOutside(e: MouseEvent | TouchEvent) {
  if (tooltipRef.value && !tooltipRef.value.contains(e.target as Node)) {
    tooltipOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleOutside)
  document.addEventListener('touchstart', handleOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleOutside)
  document.removeEventListener('touchstart', handleOutside)
})
</script>

<template>
  <div
    class="relative rounded bg-surface border border-white/8 px-5 py-4 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5"
    :class="accents[props.accent]"
  >
    <div class="flex items-center gap-1.5 mb-2">
      <p class="text-[0.65rem] uppercase tracking-widest text-muted">{{ label }}</p>
      <div v-if="props.tooltip" ref="tooltipRef" class="group relative flex items-center">
        <svg
          class="w-4 h-4 text-muted/50 hover:text-muted/80 transition-colors cursor-pointer"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          @click.stop="toggleTooltip"
        >
          <path
            fill-rule="evenodd"
            d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 01.67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 11-.671-1.34l.041-.022zM12 9a.75.75 0 100-1.5.75.75 0 000 1.5z"
            clip-rule="evenodd"
          />
        </svg>
        <div
          class="pointer-events-none absolute bottom-full right-0 mb-2 z-10 transition-opacity duration-150 w-max max-w-[min(14rem,calc(100vw-2.5rem))]"
          :class="tooltipOpen ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'"
        >
          <div
            class="bg-surface border border-white/12 rounded px-2.5 py-1.5 text-[0.65rem] text-muted/90 tracking-wide shadow-lg"
          >
            {{ props.tooltip }}
          </div>
          <div class="flex justify-end" style="padding-right: 0.375rem">
            <div class="w-2 h-2 bg-surface border-r border-b border-white/12 rotate-45 -mt-1" />
          </div>
        </div>
      </div>
    </div>
    <p class="font-pixel text-[2.2rem] leading-none text-white">{{ value }}</p>
  </div>
</template>
