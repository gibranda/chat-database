<template>
  <div>
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
useHead({
  titleTemplate: '%s - AI Database Agent',
  meta: [
    { name: 'description', content: 'Transform natural language into SQL queries with AI-powered intelligence' }
  ]
})
</script>
