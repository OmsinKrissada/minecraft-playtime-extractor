<script setup lang="ts">
import type { APIResponse } from '@/assets/types'
import { BAR_COLORS } from '@/assets/utils'
import { computed } from 'vue'
const props = defineProps<{ players: APIResponse[] }>()

let max = computed(() => props.players[0]?.playtime_hr ?? 1)

function barHeight(h: number) {
  return Math.max((h / max.value) * 100, 2).toFixed(1)
}

function label(h: number) {
  return h < 1 ? (h * 60).toFixed(0) + 'm' : h.toFixed(1) + 'h'
}
</script>

<template>
  <div class="flex items-end gap-2 h-44" role="img" aria-label="Bar chart of player playtime">
    <div
      v-for="(player, i) in players"
      :key="player.name"
      class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end"
    >
      <span class="font-pixel text-[0.9rem] text-white/80">{{ label(player.playtime_hr) }}</span>
      <div class="w-full h-full flex items-end">
        <div
          class="w-full rounded-t-[2px] min-h-1 hover:opacity-80 transition-opacity cursor-default"
          :style="{
            height: barHeight(player.playtime_hr) + '%',
            background: BAR_COLORS[i % BAR_COLORS.length],
          }"
          :title="`${player.name}: ${label(player.playtime_hr)}`"
        ></div>
      </div>
      <span class="text-[0.6rem] text-muted text-center truncate w-full">{{ player.name }}</span>
    </div>
  </div>
</template>
