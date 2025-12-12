<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child>
            <Button variant="outline" size="icon">
                <ArrowUpDown class="w-4 h-4" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
            <DropdownMenuItem @click="handleSortChange('')"
                :class="{ 'text-accent-foreground bg-accent': sort === '' }">
                <CalendarArrowDown class="w-4 h-4" />
                <span>{{ t('stream.rooms.fields.createTime') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('id:asc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'id:asc' }">
                <CalendarArrowUp class="w-4 h-4" />
                <span>{{ t('stream.rooms.fields.createTime') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('status:desc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'status:desc' }">
                <ArrowDownWideNarrow class="w-4 h-4" />
                <span>{{ t('stream.rooms.fields.monitorStatus.title') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('status:asc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'status:asc' }">
                <ArrowUpNarrowWide class="w-4 h-4" />
                <span>{{ t('stream.rooms.fields.monitorStatus.title') }}</span>
            </DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>

<script setup lang="ts">
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { ArrowUpDown, CalendarArrowDown, CalendarArrowUp, ArrowDownWideNarrow, ArrowUpNarrowWide } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

defineProps<{
    sort: string
}>()

const emit = defineEmits(['update:sort'])

const handleSortChange = (value: string) => {
    emit('update:sort', value)
}
</script>
