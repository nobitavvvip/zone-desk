<template>
  <div class="desktop-icon" @dblclick="$emit('dblclick')">
    <div class="icon-wrapper">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
      </svg>
    </div>
    <span class="icon-label">{{ app.name }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  app: { id: string; name: string; icon: string }
}>()

defineEmits<{
  dblclick: []
}>()
</script>

<style scoped>
.desktop-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  width: 80px;
  cursor: pointer;
  border-radius: 6px;
  user-select: none;
  transition: all 0.2s ease;
}

.desktop-icon:hover {
  background: var(--hover-bg-subtle);
}

.icon-wrapper {
  width: 56px;
  height: 56px;
  background: var(--hover-bg-subtle);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  transition: all 0.2s ease;
}

.desktop-icon:hover .icon-wrapper {
  background: var(--hover-bg);
  color: var(--accent);
}

.icon-wrapper svg {
  width: 32px;
  height: 32px;
}

.icon-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
  text-align: center;
  word-break: break-all;
  line-height: 1.3;
}
</style>
