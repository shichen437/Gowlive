<template>
  <Dialog :open="isOpen" @update:open="handleClose">
    <DialogContent class="sm:max-w-[800px]">
      <DialogHeader>
        <DialogTitle>{{ t('stream.rooms.linkExtract.title') }}</DialogTitle>
      </DialogHeader>
      <div class="grid grid-cols-2 gap-4 py-4">
        <div class="flex flex-col gap-2">
          <label
            class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
            >{{ t('stream.rooms.linkExtract.rawContent') }}</label
          >
          <Textarea
            v-model="rawText"
            :placeholder="t('stream.rooms.linkExtract.placeholder')"
            class="h-72 resize-none"
            @input="handleExtract"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label
            class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
            >{{ t('stream.rooms.linkExtract.result') }}</label
          >
          <Textarea
            :value="extractedLinks"
            readonly
            class="h-72 resize-none bg-muted/40"
          />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="handlePasteFromClipboard">
          <ClipboardPaste class="w-4 h-4 mr-2" />
          {{ t('stream.rooms.linkExtract.paste') }}
        </Button>
        <Button
          variant="outline"
          @click="handleCopyResult"
          :disabled="!extractedLinks"
        >
          <ClipboardCopy class="w-4 h-4 mr-2" />
          {{ t('stream.rooms.linkExtract.copy') }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose } from 'vue';
import { useI18n } from 'vue-i18n';
import { toast } from 'vue-sonner';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Textarea } from '@/components/ui/textarea';
import { ClipboardCopy, ClipboardPaste } from 'lucide-vue-next';

const { t } = useI18n();
const isOpen = ref(false);
const rawText = ref('');
const extractedLinks = ref('');

const handleExtract = () => {
  const urlRegex = /https:\/\/[^\s"'<>\)\]{}，。；！？、]+/g;
  const matches = rawText.value.match(urlRegex);
  if (matches) {
    const uniqueLinks = [...new Set(matches)];
    extractedLinks.value = uniqueLinks.join('\n');
  } else {
    extractedLinks.value = '';
  }
};

const handlePasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText();
    rawText.value = text;
    handleExtract();
  } catch {
    toast.error(t('stream.rooms.linkExtract.pasteFailed'));
  }
};

const handleCopyResult = async () => {
  if (!extractedLinks.value) return;
  try {
    await navigator.clipboard.writeText(extractedLinks.value);
    toast.success(t('stream.rooms.linkExtract.copySuccess'));
  } catch {
    toast.error(t('stream.rooms.linkExtract.copyFailed'));
  }
};

const openModal = () => {
  rawText.value = '';
  extractedLinks.value = '';
  isOpen.value = true;
};

const handleClose = () => {
  isOpen.value = false;
};

defineExpose({
  openModal,
});
</script>
