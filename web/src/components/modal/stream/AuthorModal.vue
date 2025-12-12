<template>
  <Dialog :open="isOpen" @update:open="handleOpenChange">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{{ t('stream.anchor.add.button') }}</DialogTitle>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid grid-cols-4 items-center gap-4">
          <Label for="url" class="text-right">
            {{ t('stream.anchor.fields.homeUrl') }}
          </Label>
          <Input id="url" v-model="url" class="col-span-3" />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleSubmit">
          {{ t('common.operation.save') }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
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

const { t } = useI18n();
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
      toast.error(res.msg || t('common.toast.addFailed'));
      return;
    }
    toast.success(t('common.toast.addSuccess'));
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
