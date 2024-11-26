<script setup lang="ts">
import TheWelcome from '../components/TheWelcome.vue'

import { onMounted } from 'vue';

onMounted(async () => {
  console.log('mounted')

  const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqY2J3bGtyIiwianRpIjoiMjJmZGY1ZWQtY2YxOC00MDI3LTk3NzgtZTkzNTI1M2IwMGFiIiwiaHR0cDovL3d3dy5tb3hmaWVsZC5jb20vd3MvMjAxNi8wOC9pZGVudGl0eS9jbGFpbXMvVXNlcklkIjoiMjc0MzM5IiwiaHR0cDovL3d3dy5tb3hmaWVsZC5jb20vd3MvMjAxNi8wOC9pZGVudGl0eS9jbGFpbXMvRW1haWxDb25maXJtZWQiOiJUcnVlIiwiYWR1bHQiOiJhYWMxMDcwYy03NGVmLTRkYWItOTVlOS01NzcyZjE4NzI1MzMiLCJleHAiOjE3MzI0MDA2MDgsImlzcyI6Imh0dHBzOi8vbW94ZmllbGQtYXBpLmF6dXJld2Vic2l0ZXMubmV0LyIsImF1ZCI6InVzciJ9.SEbg4D2sso2J_AyHMnWSAeNRub5gcEplxGOO84HBXQ0';

  const resp = await fetch('/api/moxfield/v3/decks', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  const decks = await resp.json()

  console.log(decks);
})
</script>

<template>
  <main>
    <TheWelcome />
  </main>
</template>
