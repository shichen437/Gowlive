<template>
    <DropdownMenu :open="isOpen" @update:open="isOpen = $event">
        <DropdownMenuTrigger as-child>
            <Button variant="outline" size="icon">
                <ListFilter class="w-4 h-4" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-64 p-4 space-y-4">
            <div class="space-y-2">
                <Label for="type">{{ t('system.logs.fields.type') }}</Label>
                <Select v-model="localFilter.type">
                    <SelectTrigger class="w-full">
                        <SelectValue :placeholder="t('system.logs.placeholders.type')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="item in typeOptions" :key="item.value" :value="String(item.value)">
                            {{ item.label }}
                        </SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <div class="space-y-2">
                <Label for="status">{{ t('system.logs.fields.status') }}</Label>
                <Select v-model="localFilter.status">
                    <SelectTrigger class="w-full">
                        <SelectValue :placeholder="t('system.logs.placeholders.status')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="item in statusOptions" :key="item.value" :value="String(item.value)">
                            {{ item.label }}
                        </SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <div class="flex justify-end space-x-2">
                <Button variant="outline" @click="handleReset">{{ t('common.operation.reset') }}</Button>
                <Button @click="handleApply">{{ t('common.operation.confirm') }}</Button>
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
import { Label } from '@/components/ui/label'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { ListFilter } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface Filter {
    type: string
    status: string
}

const props = defineProps<{
    filter: Filter
}>()

const emit = defineEmits(['update:filter'])

const typeOptions = [
    { value: 1, label: '用户' },
    { value: 2, label: '直播' },
]

const statusOptions = [
    { value: 1, label: '成功' },
    { value: 0, label: '错误' },
]

const isOpen = ref(false)
const localFilter = ref<Filter>({ ...props.filter })

watch(() => props.filter, (newFilter) => {
    localFilter.value = { ...newFilter }
})

const handleReset = () => {
    localFilter.value = { type: '', status: '' }
    emit('update:filter', localFilter.value)
    isOpen.value = false
}

const handleApply = () => {
    emit('update:filter', localFilter.value)
    isOpen.value = false
}
</script>
