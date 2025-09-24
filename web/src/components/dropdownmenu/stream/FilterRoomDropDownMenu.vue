<template>
    <DropdownMenu :open="isOpen" @update:open="isOpen = $event">
        <DropdownMenuTrigger as-child>
            <Button variant="outline" size="icon">
                <ListFilter class="w-4 h-4" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-64 p-4 space-y-4">
            <div class="space-y-2">
                <Label for="anchor">主播名称</Label>
                <Input id="anchor" v-model="localFilter.anchor" />
            </div>
            <div class="space-y-2">
                <Label for="roomName">房间名称</Label>
                <Input id="roomName" v-model="localFilter.roomName" />
            </div>
            <div class="space-y-2">
                <Label for="platform">平台</Label>
                <Select v-model="localFilter.platform">
                    <SelectTrigger class="w-full">
                        <SelectValue placeholder="选择平台" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="item in platformOptions" :key="item.value" :value="String(item.value)">
                            {{ item.label }}
                        </SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <div class="flex justify-end space-x-2">
                <Button variant="outline" @click="handleReset">重置</Button>
                <Button @click="handleApply">确定</Button>
            </div>
        </DropdownMenuContent>
    </DropdownMenu>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { ListFilter } from 'lucide-vue-next'
import { useDict } from '@/utils/useDict'

interface Filter {
    anchor: string
    roomName: string
    platform: string
}

const props = defineProps<{
    filter: Filter
}>()

const emit = defineEmits(['update:filter'])

const { options: platformOptions } = useDict('live_platform')

const isOpen = ref(false)
const localFilter = ref<Filter>({ ...props.filter })

watch(() => props.filter, (newFilter) => {
    localFilter.value = { ...newFilter }
})

const handleReset = () => {
    localFilter.value = { anchor: '', roomName: '', platform: '' }
    emit('update:filter', localFilter.value)
    isOpen.value = false
}

const handleApply = () => {
    emit('update:filter', localFilter.value)
    isOpen.value = false
}
</script>