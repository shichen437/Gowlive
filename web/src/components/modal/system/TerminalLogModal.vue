<template>
    <Dialog :open="open" @update:open="handleOpenChange">
        <DialogContent class="sm:max-w-4xl">
            <DialogHeader>
                <DialogTitle>{{ t('system.logs.terminal.title') }}</DialogTitle>
            </DialogHeader>
            <div ref="logContainer"
                class="h-[60vh] overflow-y-auto bg-black text-white p-4 rounded-md font-mono text-sm space-y-1">
                <div v-for="(log, index) in logs" :key="index" class="flex items-start">
                    <span class="text-gray-500 mr-2 whitespace-nowrap">{{ formatDate(log.time) }}</span>
                    <span :class="getLevelColor(log.level)" class="font-bold w-16 text-center mr-2">
                        [{{ log.level }}]</span>
                    <span class="whitespace-pre-wrap break-all">{{ log.msg }}</span>
                </div>
            </div>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { terminalLogs } from '@/api/system/sys_logs';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { formatDate } from '@/utils/convert';

const props = defineProps<{
    open: boolean;
    onOpenChange: (open: boolean) => void;
}>();

const { t } = useI18n();
const logs = ref<any[]>([]);
const since = ref<number | null>(null);
const limit = ref(100);
let timeoutId: ReturnType<typeof setTimeout> | null = null;
const logContainer = ref<HTMLElement | null>(null);

const getLevelColor = (level: string) => {
    if (!level) return 'text-gray-400';
    switch (level.toLowerCase()) {
        case 'info':
            return 'text-green-400';
        case 'warn':
            return 'text-yellow-400';
        case 'debu':
            return 'text-blue-400';
        default:
            return 'text-red-400';
    }
};

const scrollToBottom = () => {
    nextTick(() => {
        if (logContainer.value) {
            logContainer.value.scrollTop = logContainer.value.scrollHeight;
        }
    });
};

async function fetchLogs() {
    if (!props.open) return;

    try {
        const res = await terminalLogs({ since: since.value, limit: limit.value });
        if (res.data.rows && res.data.rows.length > 0) {
            logs.value.push(...res.data.rows);
            since.value = res.data.next;
            scrollToBottom();
            nextTick(fetchLogs);
        } else {
            since.value = res.data.next;
            if (timeoutId) clearTimeout(timeoutId);
            timeoutId = setTimeout(fetchLogs, 1000);
        }
    } catch (error) {
        console.error("Failed to fetch terminal logs:", error);
        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(fetchLogs, 3000);
    }
}

const handleOpenChange = (newOpenState: boolean) => {
    props.onOpenChange(newOpenState);
};

watch(() => props.open, (newVal) => {
    if (newVal) {
        logs.value = [];
        since.value = null;
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutId = null;
        }
        fetchLogs();
    } else {
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutId = null;
        }
    }
});

onUnmounted(() => {
    if (timeoutId) {
        clearTimeout(timeoutId);
    }
});

</script>
