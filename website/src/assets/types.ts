export interface APIResponse {
  world_run_time: number
  players: Player[]
}

export interface Player {
  name: string
  uuid: string
  playtime_hr: number
  last_seen: string
}
