<template>
    <TerminalLogModal :open="isTerminalLogModalOpen" :on-open-change="handleTerminalLogModalOpenChange" />
    <QuickAddModal :open="isQuickAddModalOpen" :on-open-change="handleQuickAddModalOpenChange" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import TerminalLogModal from '@/components/modal/system/TerminalLogModal.vue';
import QuickAddModal from '@/components/modal/stream/QuickAddModal.vue';
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts';

const isTerminalLogModalOpen = ref(false);
const isQuickAddModalOpen = ref(false);

const handleTerminalLogModalOpenChange = (open: boolean) => {
    isTerminalLogModalOpen.value = open;
};

const handleQuickAddModalOpenChange = (open: boolean) => {
    isQuickAddModalOpen.value = open;
};

useKeyboardShortcuts({
    'ctrl+t': () => {
        isTerminalLogModalOpen.value = true;
    },
    'ctrl+n': () => {
        isQuickAddModalOpen.value = true;
    },
});
</script>
