<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import type { APIResponse } from './assets/types.ts'
import PlayerRow from './components/PlayerRow.vue'
import PlayerChart from './components/PlayerChart.vue'
import StatCard from './components/StatCard.vue'
import { timeAgo } from './assets/utils.ts'

const API = 'api/'

let players = ref<APIResponse[]>([])
let totalHours = computed(() => players.value.reduce((s, p) => s + p.playtime_hr, 0))
let avgHours = computed(() => totalHours.value / players.value.length)
let maxHours = computed(() => players.value[0]?.playtime_hr ?? 1)
let updatedAt = ref(
  new Date().toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  }),
)
let refreshing = ref(false)

async function refresh() {
  refreshing.value = true
  players.value = await fetch(API).then((r) => r.json())
  updatedAt.value = new Date().toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })
  refreshing.value = false
}

let intervalId: number

onMounted(() => {
  refresh()
  intervalId = setInterval(refresh, 60_000)
})

onUnmounted(() => {
  clearInterval(intervalId)
})
</script>

<template>
  <!-- Glow blobs -->
  <div
    class="fixed w-[500px] h-[500px] rounded-full bg-grass opacity-[0.12] blur-[120px] -top-36 -left-24 pointer-events-none z-0"
    aria-hidden="true"
  ></div>
  <div
    class="fixed w-[400px] h-[400px] rounded-full bg-diamond opacity-[0.12] blur-[120px] -bottom-24 -right-24 pointer-events-none z-0"
    aria-hidden="true"
  ></div>

  <div class="relative z-10 max-w-[860px] mx-auto px-6 pt-10 pb-16">
    <!-- Header -->
    <header class="flex items-start gap-5 mb-12">
      <img
        class="shrink-0 w-[72px] h-[72px] border-2 border-white/8 rounded bg-surface"
        style="image-rendering: pixelated"
        viewBox="0 0 64 64"
        xmlns="http://www.w3.org/2000/svg"
        aria-hidden="true"
        src="/star-tech-logo.png"
      />

      <div>
        <h1
          class="font-pixel text-5xl leading-none tracking-wide text-grass [text-shadow:3px_3px_0_rgba(0,0,0,0.5)] mb-1"
        >
          STAR <span class="text-white">TECH</span>
        </h1>
        <p class="text-[0.75rem] text-muted tracking-wide">
          <!-- Minecraft Server · Playtime Leaderboard -->
          <!-- Currently Theta 1 Hotfix 3 -->
          World started {{ timeAgo('2026-05-17T07:00:00+07:00', true) }} · Currently Theta 1 Hotfix
          3
        </p>

        <div class="flex items-center gap-2 mt-3 flex-wrap">
          <!-- Live pill -->
          <span
            class="inline-flex items-center gap-1.5 text-[0.7rem] text-grass border border-grass/30 bg-grass/8 px-2.5 py-1 rounded-[2px] tracking-wider"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-grass animate-pulse-dot"></span>
            LIVE
          </span>
        </div>
      </div>
    </header>

    <!-- Stat cards -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-10">
      <StatCard label="Total players" :value="players.length" accent="green" />
      <StatCard label="Total hours" :value="totalHours.toFixed(0)" accent="blue" />
      <StatCard label="Top player" :value="players[0]?.name ?? '—'" accent="gold" />
      <StatCard label="Avg playtime" :value="avgHours.toFixed(1) + 'h'" accent="stone" />
    </div>

    <!-- Chart -->
    <section class="mb-10">
      <div class="flex items-center gap-2.5 mb-4">
        <span class="text-[0.65rem] uppercase tracking-[0.15em] text-muted">Hours played</span>
        <div class="flex-1 h-px bg-white/8"></div>
      </div>
      <div class="bg-surface border border-white/8 rounded p-5">
        <PlayerChart :players="players" />
      </div>
    </section>

    <!-- Leaderboard -->
    <div class="flex items-center gap-2.5 mb-4">
      <span class="text-[0.65rem] uppercase tracking-[0.15em] text-muted">Leaderboard</span>
      <div class="flex-1 h-px bg-white/8"></div>
    </div>
    <div class="flex flex-col gap-2">
      <PlayerRow v-for="(player, i) in players" :player="player" :index="i" :max="maxHours" />
    </div>
  </div>

  <!-- Footer -->
  <footer
    class="relative z-10 text-center text-[0.65rem] text-muted tracking-wider px-6 py-6 border-t border-white/8"
  >
    <!-- Updated {{ updatedAt }} &nbsp;·&nbsp; -->
    API available at
    <a
      href="https://mc.krissada.com/star_tech/playtime/api"
      target="_blank"
      rel="noopener"
      class="text-grass hover:underline"
      >mc.krissada.com/star_tech/playtime/api</a
    >
  </footer>
</template>
