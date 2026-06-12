<script setup lang="ts">
import { computed, ref } from 'vue'
import { BAR_COLORS, timeAgo, isOnline, fmtValue, fmtUnit, avatarUrl } from '@/assets/utils.ts'
import type { APIResponse as PlayerResponse } from '@/assets/types'

const props = defineProps<{
  player: PlayerResponse
  index: number
  max: number
}>()

const color = computed(() => BAR_COLORS[props.index % BAR_COLORS.length])
const pct = computed(() =>
  props.max > 0 ? ((props.player.playtime_hr / props.max) * 100).toFixed(1) : 0,
)
const online = computed(() => isOnline(props.player.last_seen))
const imgError = ref(false)

const rankStyles = [
  'font-pixel text-[1.3rem] text-gold drop-shadow-[0_0_8px_rgba(245,197,24,0.5)]',
  'font-pixel text-[1.3rem] text-iron',
  'font-pixel text-[1.3rem] text-[#CD7F32]',
  'font-mono text-[0.8rem] text-muted',
]
const rankStyle = computed(() => (props.index < 3 ? rankStyles[props.index] : rankStyles[3]))
const rankLabel = computed(() =>
  props.index < 3 ? String(props.index + 1) : `#${props.index + 1}`,
)
</script>

<template>
  <div
    class="animate-slide-in bg-surface border border-white/8 rounded px-5 py-4 grid items-center gap-3.5 hover:border-white/15 hover:bg-surface2 transition-colors"
    :style="`grid-template-columns:2rem 2.5rem 1fr auto; animation-delay:${index * 0.06}s`"
  >
    <!-- Rank -->
    <span class="text-center leading-none" :class="rankStyle">{{ rankLabel }}</span>

    <!-- Avatar -->
    <div
      class="w-10 h-10 rounded-[3px] border border-white/8 bg-surface2 flex items-center justify-center overflow-hidden shrink-0"
    >
      <img
        v-if="!imgError"
        :src="avatarUrl(player.uuid)"
        :alt="`${player.name}'s avatar`"
        class="w-full h-full object-cover"
        style="image-rendering: pixelated"
        @error="imgError = true"
      />
      <span v-else class="font-pixel text-lg text-white/60">
        {{ player.name.slice(0, 2).toUpperCase() }}
      </span>
    </div>

    <!-- Name + meta -->
    <div>
      <p class="text-[0.9rem] font-semibold text-white tracking-wide mb-1.5">{{ player.name }}</p>
      <span
        v-if="online"
        class="inline-flex items-center gap-1 text-[0.6rem] text-grass border border-grass/35 bg-grass/10 px-1.5 py-0.5 rounded-[2px] tracking-wide"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-grass animate-pulse-dot"></span>
        ONLINE
      </span>
      <span v-else class="text-[0.65rem] text-muted tracking-wide">{{
        timeAgo(player.last_seen)
      }}</span>
    </div>

    <!-- Playtime -->
    <div class="text-right min-w-[80px]">
      <p class="font-pixel text-[1.5rem] leading-none text-white">
        {{ fmtValue(player.playtime_hr) }}
      </p>
      <p class="text-[0.6rem] text-muted tracking-widest mt-0.5">
        {{ fmtUnit(player.playtime_hr) }}
      </p>
      <div class="h-1 bg-white/6 rounded-[1px] overflow-hidden mt-1.5 w-20 ml-auto">
        <div class="h-full rounded-[1px]" :style="{ width: pct + '%', background: color }"></div>
      </div>
    </div>
  </div>
</template>
