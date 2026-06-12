<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    value: string | number
    accent?: 'green' | 'blue' | 'gold' | 'stone'
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
</script>

<template>
  <div
    class="relative overflow-hidden rounded bg-surface border border-white/8 px-5 py-4 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5"
    :class="accents[props.accent]"
  >
    <p class="text-[0.65rem] uppercase tracking-widest text-muted mb-2">{{ label }}</p>
    <p class="font-pixel text-[2.2rem] leading-none text-white">{{ value }}</p>
  </div>
</template>
