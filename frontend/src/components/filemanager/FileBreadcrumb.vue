<template>
  <div class="breadcrumb">
    <template v-for="(segment, index) in segments" :key="index">
      <span
        class="crumb"
        :class="{ active: index === segments.length - 1 }"
        @click="$emit('navigate', segment.path)"
      >
        {{ segment.name }}
      </span>
      <span v-if="index < segments.length - 1" class="separator">&gt;</span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  path: string
}>()

defineEmits<{
  navigate: [path: string]
}>()

const segments = computed(() => {
  const parts = props.path.split('/').filter(Boolean)
  const result = [{ name: '/', path: '/' }]
  let current = ''
  for (const part of parts) {
    current += '/' + part
    result.push({ name: part, path: current })
  }
  return result
})
</script>

<style scoped>
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 14px;
  color: var(--text-muted);
}

.crumb {
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
}

.crumb:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.crumb.active {
  color: var(--text-primary);
  font-weight: 500;
}

.separator {
  color: #475569;
  font-size: 12px;
}
</style>
