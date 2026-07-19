<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import type { APIResponse, Player } from './assets/types.ts'
import PlayerRow from './components/PlayerRow.vue'
import PlayerChart from './components/PlayerChart.vue'
import StatCard from './components/StatCard.vue'
import { timeAgo } from './assets/utils.ts'

const API = import.meta.env.VITE_BASE_API || 'api/'

const worldCreationTime = '2026-05-17T07:00:00+07:00'
const REFRESH_INTERVAL = 30

let players = ref<Player[]>([])
let worldRunTime = ref<number | null>(null)

let totalHours = computed(() => players.value.reduce((s, p) => s + p.playtime_hr, 0))
let maxHours = computed(() => players.value[0]?.playtime_hr ?? 1)
let updatedAt = ref(
  new Date().toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  }),
)
let refreshing = ref(false)
let lastRefreshFailed = ref(false)

async function refresh() {
  refreshing.value = true

  try {
    const response = (await fetch(API).then((r) => r.json())) as APIResponse
    players.value = response.players
    worldRunTime.value = Math.round(response.world_run_time)
    updatedAt.value = new Date().toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
    })
    lastRefreshFailed.value = false
  } catch (err) {
    lastRefreshFailed.value = true
  }

  refreshing.value = false
  refreshCountdown.value = REFRESH_INTERVAL
}

let refreshTimer: number
const refreshCountdown = ref(0)

onMounted(() => {
  refresh()

  refreshTimer = setInterval(() => {
    refreshCountdown.value--
    if (refreshCountdown.value <= 0 && !refreshing.value) {
      refresh()
    }
  }, 1000)
})

onUnmounted(() => {
  clearInterval(refreshTimer)
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
    <header class="flex items-start gap-5 mb-6">
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
          <span v-if="worldCreationTime">
            World started {{ timeAgo(worldCreationTime, true) }} ·
          </span>
          Currently Theta 1 Hotfix 3
        </p>
      </div>
    </header>

    <div class="flex items-center gap-2 mb-6 flex-wrap">
      <!-- Live pill -->
      <span
        class="inline-flex items-center gap-1.5 text-[0.7rem] border px-2.5 py-1 rounded-[2px] tracking-wider"
        :class="{
          'text-red-400 border-red-400/30 bg-red-400/8': lastRefreshFailed,
          'text-grass border-grass/30 bg-grass/8': !lastRefreshFailed,
        }"
      >
        <span
          class="w-1.5 h-1.5 rounded-full animate-pulse-dot"
          :class="{ 'bg-red-400': lastRefreshFailed, 'bg-grass': !lastRefreshFailed }"
        ></span>
        {{ refreshCountdown > 0 ? `Updating in ${refreshCountdown}s` : 'Updating' }}
        <svg
          v-if="refreshCountdown <= 0"
          class="animate-spin-slow"
          :class="{ 'text-red-400': lastRefreshFailed, 'text-grass': !lastRefreshFailed }"
          width="11"
          height="11"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
        >
          <path d="M23 4v6h-6" />
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
        </svg>
      </span>
    </div>

    <!-- Stat cards -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-10">
      <StatCard label="Total players" :value="players.length" accent="green" />
      <StatCard label="World run hours" :value="worldRunTime ?? 'N/A'" accent="blue" />
      <StatCard label="Total hours" :value="totalHours.toFixed(0)" accent="gold" />
      <StatCard
        label="Wage Equivalent"
        :value="(totalHours * 50).toFixed(0) + ' THB'"
        accent="stone"
        tooltip="Calculated from total hours based on minimum hourly rate in Bangkok (50 THB/hr)"
      />
    </div>

    <!-- Chart -->
    <section class="mb-10">
      <div class="flex items-center gap-2.5 mb-4">
        <span class="text-[0.65rem] uppercase tracking-[0.15em] text-muted">Hours played</span>
        <div class="flex-1 h-px bg-white/8"></div>
      </div>
      <div class="bg-surface border border-white/8 rounded p-5 overflow-auto">
        <PlayerChart :players="players" />
      </div>
    </section>

    <!-- Leaderboard -->
    <div class="flex items-center gap-2.5 mb-4">
      <span class="text-[0.65rem] uppercase tracking-[0.15em] text-muted">Leaderboard</span>
      <div class="flex-1 h-px bg-white/8"></div>
    </div>
    <div class="flex flex-col gap-2 mb-10">
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
      >mc.krissada.com/star_tech/playtime/api
    </a>
  </footer>
</template>
