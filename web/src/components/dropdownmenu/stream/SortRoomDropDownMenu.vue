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
                <span>添加时间</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('id:asc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'id:asc' }">
                <CalendarArrowUp class="w-4 h-4" />
                <span>添加时间</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('status:desc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'status:desc' }">
                <ArrowDownWideNarrow class="w-4 h-4" />
                <span>监控状态</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSortChange('status:asc')"
                :class="{ 'text-accent-foreground bg-accent': sort === 'status:asc' }">
                <ArrowUpNarrowWide class="w-4 h-4" />
                <span>监控状态</span>
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

defineProps<{
    sort: string
}>()

const emit = defineEmits(['update:sort'])

const handleSortChange = (value: string) => {
    emit('update:sort', value)
}
</script>