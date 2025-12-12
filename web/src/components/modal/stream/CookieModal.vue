<template>
    <Dialog :open="true" @update:open="handleClose">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ isEditMode ? t('stream.cookie.edit.title') : t('stream.cookie.add.title') }}
                </DialogTitle>
                <DialogDescription>
                    {{ isEditMode ? t('stream.cookie.edit.desc') : t('stream.cookie.add.desc') }}
                </DialogDescription>
            </DialogHeader>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="platform" class="text-right">{{ t('common.fields.platform') }}</Label>
                    <Select v-model="form.platform">
                        <SelectTrigger class="col-span-3 w-full">
                            <SelectValue :placeholder="t('stream.cookie.placeholder.platform')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem v-for="item in platformDict" :key="item.dictValue" :value="item.dictValue">
                                {{ item.dictLabel }}
                            </SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="cookie" class="text-right">Cookie</Label>
                    <Input id="cookie" v-model="form.cookie" class="col-span-3" />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="remark" class="text-right">{{ t('common.fields.remark') }}</Label>
                    <Input id="remark" v-model="form.remark" class="col-span-3" />
                </div>
            </div>
            <DialogFooter>
                <Button type="button" variant="secondary" @click="handleClose">{{ t('common.operation.cancel')
                    }}</Button>
                <Button type="submit" @click="handleSubmit">{{ t('common.operation.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { addCookie, updateCookie } from '@/api/stream/live_cookie';
import type { LiveCookie } from '@/types/stream';
import { toast } from 'vue-sonner';
import { useDict } from '@/utils/useDict';

const { t } = useI18n();

const props = defineProps<{
    cookie: LiveCookie | null;
}>();

const emits = defineEmits(['close', 'success']);

const { dict: platformDict } = useDict("live_platform");

const form = ref<Partial<LiveCookie>>({
    platform: '',
    cookie: '',
    remark: '',
});

const isEditMode = ref(false);

watch(
    () => props.cookie,
    (newVal) => {
        if (newVal) {
            form.value = { ...newVal };
            isEditMode.value = true;
        } else {
            form.value = {
                platform: '',
                cookie: '',
                remark: '',
            };
            isEditMode.value = false;
        }
    },
    { immediate: true }
);

const handleSubmit = async () => {
    try {
        if (isEditMode.value) {
            const res: any = await updateCookie(form.value);
            if (res.code !== 0) {
                toast.error(res.msg || t('common.toast.updateFailed'));
                return;
            }
            toast.success(t('common.toast.updateSuccess'));
        } else {
            const res: any = await addCookie(form.value);
            if (res.code !== 0) {
                toast.error(res.msg || t('common.toast.addFailed'));
                return;
            }
            toast.success(t('common.toast.addSuccess'));
        }
        emits('success');
    } catch (error) {
        console.error('Failed to save cookie:', error);
    }
};

const handleClose = () => {
    emits('close');
};
</script>
