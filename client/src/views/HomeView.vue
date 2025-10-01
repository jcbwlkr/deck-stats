<script setup lang="ts">
import TheWelcome from '../components/TheWelcome.vue'

import { onMounted } from 'vue';

onMounted(async () => {
  console.log('mounted')

  const token = '';

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
