<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="outline" size="icon">
        <Filter class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" class="w-64 p-4 space-y-4">
      <div class="space-y-2">
        <Label for="platform">平台</Label>
        <Select v-model="localFilter.platform">
          <SelectTrigger class="w-full">
            <SelectValue placeholder="选择平台" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="item in platformOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="space-y-2">
        <Label for="nickname">主播名称</Label>
        <Input id="nickname" v-model="localFilter.nickname" placeholder="输入主播名称" />
      </div>
      <div class="flex justify-end space-x-2">
        <Button variant="ghost" @click="resetFilter">重置</Button>
        <Button @click="applyFilter">确定</Button>
      </div>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Filter } from 'lucide-vue-next';
import { useDict } from '@/utils/useDict';

const props = defineProps<{
  filter: {
    platform: string;
    nickname: string;
  };
}>();

const emit = defineEmits(['update:filter']);

const { options: platformOptions } = useDict('live_platform');

const localFilter = ref({ ...props.filter });

watch(() => props.filter, (newFilter) => {
  localFilter.value = { ...newFilter };
}, { deep: true });

const applyFilter = () => {
  emit('update:filter', { ...localFilter.value });
};

const resetFilter = () => {
  localFilter.value = {
    platform: '',
    nickname: '',
  };
  emit('update:filter', { ...localFilter.value });
};
</script>
