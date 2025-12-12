<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child>
            <Button variant="outline" size="icon">
                <Download class="w-4 h-4" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
            <DropdownMenuItem @click="handleExport(0)">
                {{ t('stream.rooms.fields.export') }} TXT
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleExport(1)">
                {{ t('stream.rooms.fields.export') }} Excel
            </DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>

<script setup lang="ts">
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Download } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const emit = defineEmits(['export']);

const handleExport = (exportType: number) => {
    emit('export', exportType);
};
</script>
