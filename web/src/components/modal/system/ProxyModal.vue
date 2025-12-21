<template>
    <Dialog :open="true" @update:open="handleClose">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ isEditMode ? t('system.proxy.edit.title') : t('system.proxy.add.title') }}
                </DialogTitle>
                <DialogDescription>
                    {{ isEditMode ? t('system.proxy.edit.desc') : t('system.proxy.add.desc') }}
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
                    <Label for="proxy" class="text-right">{{ t('system.proxy.fields.proxy') }}</Label>
                    <Textarea id="proxy" v-model="proxyText" class="col-span-3"
                        :placeholder="t('system.proxy.placeholder.proxy')" />
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
import { Textarea } from '@/components/ui/textarea';
import { addProxy, updateProxy } from '@/api/system/sys_proxy';
import type { SysProxy } from '@/types/system';
import { toast } from 'vue-sonner';
import { useDict } from '@/utils/useDict';


const { t } = useI18n();
const { dict: platformDict } = useDict("live_platform");

const props = defineProps<{
    proxy: SysProxy | null;
}>();

const emits = defineEmits(['close', 'success']);

const form = ref<Partial<SysProxy>>({
    platform: '',
    proxy: [],
    remark: '',
});

const proxyText = ref('');

const isEditMode = ref(false);

watch(
    () => props.proxy,
    (newVal) => {
        if (newVal) {
            form.value = { ...newVal };
            proxyText.value = newVal.proxy ? newVal.proxy : '';
            isEditMode.value = true;
        } else {
            form.value = {
                platform: '',
                proxy: [],
                remark: '',
            };
            proxyText.value = '';
            isEditMode.value = false;
        }
    },
    { immediate: true }
);

const handleSubmit = async () => {
    try {
        form.value.proxy = proxyText.value
            .split('\n')
            .flatMap(line => line.split(','))
            .map(s => s.trim())
            .filter(s => s.length > 0);

        if (isEditMode.value) {
            const res: any = await updateProxy(form.value);
            if (res.code !== 0) {
                toast.error(res.msg || t('common.toast.updateFailed'));
                return;
            }
            toast.success(t('common.toast.updateSuccess'));
        } else {
            const res: any = await addProxy(form.value);
            if (res.code !== 0) {
                toast.error(res.msg || t('common.toast.addFailed'));
                return;
            }
            toast.success(t('common.toast.addSuccess'));
        }
        emits('success');
    } catch (error) {
        console.error('Failed to save proxy:', error);
    }
};

const handleClose = () => {
    emits('close');
};
</script>
