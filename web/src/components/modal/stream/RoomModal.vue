<template>
    <Dialog :open="isOpen" @update:open="handleClose">
        <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
                <DialogTitle>{{ dialogTitle }}</DialogTitle>
            </DialogHeader>
            <form @submit="onSubmit">
                <div class="grid gap-4 py-4">
                    <FormField v-slot="{ componentField, errorMessage }" name="roomUrl">
                        <FormItem>
                            <FormLabel>{{ t('stream.rooms.fields.roomUrl') }}</FormLabel>
                            <FormControl>
                                <Input type="text" :placeholder="t('stream.rooms.placeholder.roomUrl')"
                                    v-bind="componentField" :disabled="isEditMode"
                                    :class="{ 'border-red-500': errorMessage }" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>

                    <div class="grid grid-cols-2 gap-4">
                        <FormField v-slot="{ componentField, errorMessage }" name="interval">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.interval') }}</FormLabel>
                                <FormControl>
                                    <Input type="number" placeholder="30-600" v-bind="componentField"
                                        :class="{ 'border-red-500': errorMessage }" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>

                        <FormField v-slot="{ componentField, errorMessage }" name="format">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.format') }}</FormLabel>
                                <Select v-bind="componentField">
                                    <FormControl>
                                        <SelectTrigger class="w-full" :class="{ 'border-red-500': errorMessage }">
                                            <SelectValue :placeholder="t('stream.rooms.placeholder.format')" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="w-[--radix-select-trigger-width]">
                                        <SelectItem value="flv">FLV</SelectItem>
                                        <SelectItem value="mkv">MKV</SelectItem>
                                        <SelectItem value="ts">TS</SelectItem>
                                        <SelectItem value="mp4">MP4</SelectItem>
                                        <SelectItem value="mp3">{{ t('stream.rooms.fields.onlyAudio') }}</SelectItem>
                                    </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <FormField v-slot="{ componentField, errorMessage }" name="quality">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.quality.title') }}
                                    <TooltipProvider>
                                        <Tooltip>
                                            <TooltipTrigger as-child>
                                                <Info class="w-4 h-4 ml-1" />
                                            </TooltipTrigger>
                                            <TooltipContent>
                                                <p>{{ t('stream.rooms.fields.quality.tooltip') }}</p>
                                            </TooltipContent>
                                        </Tooltip>
                                    </TooltipProvider>
                                </FormLabel>
                                <Select v-bind="componentField">
                                    <FormControl>
                                        <SelectTrigger class="w-full" :class="{ 'border-red-500': errorMessage }">
                                            <SelectValue :placeholder="t('stream.rooms.placeholder.quality')" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="w-[--radix-select-trigger-width]">
                                        <SelectItem :value="0">{{ t('stream.rooms.fields.quality.original') }}
                                        </SelectItem>
                                        <SelectItem :value="1">{{ t('stream.rooms.fields.quality.super') }}</SelectItem>
                                        <SelectItem :value="2">{{ t('stream.rooms.fields.quality.high') }}</SelectItem>
                                        <SelectItem :value="3">{{ t('stream.rooms.fields.quality.medium') }}
                                        </SelectItem>
                                        <SelectItem :value="4">{{ t('stream.rooms.fields.quality.low') }}</SelectItem>
                                    </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                        <FormField v-slot="{ componentField, errorMessage }" name="segmentTime">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.segment') }}</FormLabel>
                                <Select v-bind="componentField">
                                    <FormControl>
                                        <SelectTrigger class="w-full" :class="{ 'border-red-500': errorMessage }">
                                            <SelectValue :placeholder="t('stream.rooms.placeholder.segment')" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="w-[--radix-select-trigger-width]">
                                        <SelectItem :value="0">{{ t('stream.rooms.fields.segmentNone') }}</SelectItem>
                                        <SelectItem :value="900">15 {{ t('common.fields.minute') }}</SelectItem>
                                        <SelectItem :value="1800">30 {{ t('common.fields.minute') }}</SelectItem>
                                        <SelectItem :value="3600">1 {{ t('common.fields.hour') }}</SelectItem>
                                        <SelectItem :value="7200">2 {{ t('common.fields.hour') }}</SelectItem>
                                        <SelectItem :value="14400">4 {{ t('common.fields.hour') }}</SelectItem>
                                    </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                    </div>

                    <FormField v-slot="{ field, errorMessage }" name="monitorType">
                        <FormItem>
                            <FormLabel class="flex items-center">
                                {{ t('stream.rooms.fields.monitorType.title') }}
                                <TooltipProvider>
                                    <Tooltip>
                                        <TooltipTrigger as-child>
                                            <Info class="w-4 h-4 ml-1" />
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>{{ t('stream.rooms.fields.monitorType.realTimeTip') }}</p>
                                            <p>{{ t('stream.rooms.fields.monitorType.cronTip') }}</p>
                                            <p>{{ t('stream.rooms.fields.monitorType.intelligentTip') }}</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </TooltipProvider>
                            </FormLabel>
                            <FormControl>
                                <RadioGroup class="grid grid-cols-4 gap-2" :model-value="field.value"
                                    @update:model-value="field.onChange">
                                    <div v-for="item in monitorTypeOptions" :key="item.value">
                                        <RadioGroupItem :id="`monitor-type-${item.value}`" :value="item.value"
                                            class="sr-only" />
                                        <label :for="`monitor-type-${item.value}`" :class="cn(
                                            buttonVariants({ variant: field.value === item.value ? 'default' : 'outline' }),
                                            'w-full cursor-pointer flex items-center justify-center h-9',
                                            errorMessage ? 'border-red-500' : ''
                                        )">
                                            {{ item.label }}
                                        </label>
                                    </div>
                                </RadioGroup>
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>

                    <template v-if="values.monitorType === 2">
                        <div class="grid grid-cols-2 gap-4">
                            <FormField v-slot="{ componentField, errorMessage }" name="monitorStartAt">
                                <FormItem>
                                    <FormLabel>{{ t('stream.rooms.fields.monitorStartTime') }}</FormLabel>
                                    <FormControl>
                                        <VueDatePicker v-bind="componentField" time-picker auto-apply model-type="HH:mm"
                                            :format="'HH:mm'" text-input placeholder="HH:mm"
                                            :class="{ 'border-red-500': errorMessage }" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                            <FormField v-slot="{ componentField, errorMessage }" name="monitorStopAt">
                                <FormItem>
                                    <FormLabel>{{ t('stream.rooms.fields.monitorEndTime') }}</FormLabel>
                                    <FormControl>
                                        <VueDatePicker v-bind="componentField" time-picker auto-apply model-type="HH:mm"
                                            :format="'HH:mm'" text-input placeholder="HH:mm"
                                            :class="{ 'border-red-500': errorMessage }" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                        </div>
                    </template>

                    <FormField v-slot="{ componentField, errorMessage }" name="remark">
                        <FormItem>
                            <FormLabel>{{ t('common.fields.remark') }}</FormLabel>
                            <FormControl>
                                <Input type="text" :placeholder="t('stream.rooms.placeholder.remark')"
                                    v-bind="componentField" :class="{ 'border-red-500': errorMessage }" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>
                </div>
                <DialogFooter>
                    <Button type="submit" :disabled="isSubmitting">
                        {{ isSubmitting ? t('common.operation.saving') : t('common.operation.save') }}
                    </Button>
                </DialogFooter>
            </form>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose, defineEmits, computed } from 'vue';
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import { useForm } from 'vee-validate';
import { addRoom, updateRoom, roomDetail } from '@/api/stream/live_manage';
import { toast } from 'vue-sonner';
import { cn } from '@/lib/utils';
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
import type { LiveManage } from '@/types/stream';
import { Button, buttonVariants } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Info } from 'lucide-vue-next';
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger
} from '@/components/ui/tooltip';
import { Input } from '@/components/ui/input';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const isOpen = ref(false);
const emit = defineEmits(['refresh']);
const roomId = ref<number | null>(null);
const isSubmitting = ref(false);

const isEditMode = computed(() => !!roomId.value);
const dialogTitle = computed(() => isEditMode.value ? t('stream.rooms.edit.title') : t('stream.rooms.add.title'));

const monitorTypeOptions = computed(() => [
    { value: 0, label: t('stream.rooms.fields.monitorType.stop') },
    { value: 1, label: t('stream.rooms.fields.monitorType.realtime') },
    { value: 2, label: t('stream.rooms.fields.monitorType.cron') },
    { value: 3, label: t('stream.rooms.fields.monitorType.intelligent') }
]);

const formSchema = toTypedSchema(
    z.object({
        roomUrl: z.string().url({ message: t('stream.rooms.valid.validUrl') }),
        interval: z.coerce.number().min(30, { message: t('stream.rooms.valid.intervalMin') }).max(600, { message: t('stream.rooms.valid.intervalMax') }),
        format: z.enum(['flv', 'mp4', 'mkv', 'ts', 'mp3']),
        quality: z.coerce.number(),
        segmentTime: z.coerce.number(),
        monitorType: z.coerce.number(),
        monitorStartAt: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, { message: t('stream.rooms.valid.timeFormat') }).optional().nullable(),
        monitorStopAt: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, { message: t('stream.rooms.valid.timeFormat') }).optional().nullable(),
        remark: z.string().max(45, { message: t('stream.rooms.valid.remarkLength') }).optional().nullable(),
    }).refine(data => {
        if (data.monitorType === 2) {
            return !!data.monitorStartAt && !!data.monitorStopAt;
        }
        return true;
    }, {
        message: t('stream.rooms.valid.cronTime'),
        path: ['monitorStartAt'],
    })
);

const { handleSubmit, values, setValues, resetForm } = useForm({
    validationSchema: formSchema,
    initialValues: {
        roomUrl: '',
        interval: 30,
        format: 'flv',
        monitorType: 0,
        quality: 0,
        segmentTime: 0,
        remark: '',
        monitorStartAt: null,
        monitorStopAt: null,
    }
});

const onSubmit = handleSubmit(async (formValues) => {
    if (isSubmitting.value) return;

    isSubmitting.value = true;

    try {
        const payload = {
            ...formValues,
            monitorType: Number(formValues.monitorType)
        };

        if (isEditMode.value && roomId.value) {
            const res: any = await updateRoom({ id: roomId.value, ...payload });
            if (res.code !== 0) {
                toast.error(res.data.msg || t('common.toast.updateFailed'));
                return;
            }
            toast.success(t('common.toast.updateSuccess'));
        } else {
            const res: any = await addRoom(payload);
            if (res.code !== 0) {
                toast.error(res.data.msg || t('common.toast.addFailed'));
                return;
            }
            toast.success(t('common.toast.addSuccess'));
        }

        isOpen.value = false;
        emit('refresh');
    } catch (error) {
        console.error('Failed to save room:', error);
    } finally {
        isSubmitting.value = false;
    }
});

const openModal = async (id?: number) => {
    resetForm();
    if (id) {
        roomId.value = id;
        try {
            const response = await roomDetail(id);
            const roomData = response.data.data as LiveManage;

            const formattedData = {
                ...roomData,
                quality: Number(roomData.quality),
                segmentTime: Number(roomData.segmentTime),
                monitorType: Number(roomData.monitorType),
                monitorStartAt: roomData.monitorStartAt || null,
                monitorStopAt: roomData.monitorStopAt || null,
                remark: roomData.remark || ''
            };

            setValues(formattedData);
        } catch (error) {
            console.error('Failed to fetch room details:', error);
            toast.error(t('stream.rooms.toast.detailErr'));
            return;
        }
    } else {
        roomId.value = null;
    }
    isOpen.value = true;
};

const handleClose = () => {
    isOpen.value = false;
};

defineExpose({
    openModal,
});
</script>
