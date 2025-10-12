<template>
  <Dialog :open="isOpen" @update:open="handleOpenChange">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>添加主播</DialogTitle>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid grid-cols-4 items-center gap-4">
          <Label for="url" class="text-right">
            主页链接
          </Label>
          <Input id="url" v-model="url" class="col-span-3" />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleSubmit">
          保存
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { addAnchor } from '@/api/stream/anchor_info';
import { toast } from 'vue-sonner';

const isOpen = ref(false);
const url = ref('');

const emit = defineEmits(['refresh']);

const openModal = () => {
  isOpen.value = true;
  url.value = '';
};

const handleOpenChange = (open: boolean) => {
  isOpen.value = open;
};

const handleSubmit = async () => {
  try {
    const res: any = await addAnchor({ url: url.value });
    if (res.code !== 0) {
      toast.error(res.msg || '添加失败');
      return;
    }
    toast.success('添加成功');
    isOpen.value = false;
    emit('refresh');
  } catch (error) {
    console.error('Failed to add anchor:', error);
  }
};

defineExpose({
  openModal,
});
</script>
