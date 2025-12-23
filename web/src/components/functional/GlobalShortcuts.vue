<template>
    <TerminalLogModal :open="isTerminalLogModalOpen" :on-open-change="handleTerminalLogModalOpenChange" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import TerminalLogModal from '@/components/modal/system/TerminalLogModal.vue';
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts';

const isTerminalLogModalOpen = ref(false);

const handleTerminalLogModalOpenChange = (open: boolean) => {
    isTerminalLogModalOpen.value = open;
};

useKeyboardShortcuts({
    'ctrl+t': () => {
        isTerminalLogModalOpen.value = true;
    },
});
</script>
