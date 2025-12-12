<template>
    <Dialog :open="isOpen" @update:open="handleClose">
        <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
                <DialogTitle>{{ t('stream.rooms.batch.title') }}</DialogTitle>
            </DialogHeader>
            <form @submit="onSubmit">
                <div class="grid gap-4 py-4">
                    <FormField v-slot="{ componentField, errorMessage }" name="roomUrls">
                        <FormItem>
                            <FormLabel>{{ t('stream.rooms.fields.roomUrl') }}</FormLabel>
                            <FormControl>
                                <Textarea :placeholder="t('stream.rooms.placeholder.batchUrls')" v-bind="componentField"
                                    :class="{ 'border-red-500': errorMessage }" class="h-32 max-h-64 min-h-24" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>

                    <div class="grid grid-cols-2 gap-4">
                        <FormField v-slot="{ componentField, errorMessage }" name="interval">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.interval') }}</FormLabel>
                                <FormControl>
                                    <Input type="number" placeholder="30-600" v-bind="componentField" :class="{
                                        'border-red-500': errorMessage,
                                    }" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>

                        <FormField v-slot="{ componentField, errorMessage }" name="format">
                            <FormItem>
                                <FormLabel>{{ t('stream.rooms.fields.format') }}</FormLabel>
                                <Select v-bind="componentField">
                                    <FormControl>
                                        <SelectTrigger class="w-full" :class="{
                                            'border-red-500': errorMessage,
                                        }">
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
import { ref, defineExpose, defineEmits } from "vue";
import { toTypedSchema } from "@vee-validate/zod";
import * as z from "zod";
import { useForm } from "vee-validate";
import { addBatchRoom } from "@/api/stream/live_manage";
import { toast } from "vue-sonner";
import { useI18n } from 'vue-i18n';

import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

const { t } = useI18n();
const isOpen = ref(false);
const emit = defineEmits(["refresh"]);
const isSubmitting = ref(false);

const formSchema = toTypedSchema(
    z.object({
        roomUrls: z.string().min(1, { message: t('stream.rooms.placeholder.roomUrl') }),
        interval: z.coerce
            .number()
            .min(30, { message: t('stream.rooms.valid.intervalMin') })
            .max(600, { message: t('stream.rooms.valid.intervalMax') }),
        format: z.enum(["flv", "mp4", "mp3", "mkv", "ts"]),
        remark: z
            .string()
            .max(45, { message: t('stream.rooms.valid.remarkLength') })
            .optional()
            .nullable(),
    }),
);

const { handleSubmit, resetForm } = useForm({
    validationSchema: formSchema,
    initialValues: {
        roomUrls: "",
        interval: 30,
        format: "flv",
        remark: "",
    },
});

const onSubmit = handleSubmit(async (formValues) => {
    if (isSubmitting.value) return;

    const urls = formValues.roomUrls
        .split(/[\n,]/)
        .map((url) => url.trim())
        .filter((url) => url);
    if (urls.length === 0) {
        toast.error(t('stream.rooms.valid.validUrl'));
        return;
    }
    if (urls.length > 30) {
        toast.error(t('stream.rooms.valid.batchUrlsNum'));
        return;
    }

    const urlRegex =
        /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/;
    for (const url of urls) {
        if (!urlRegex.test(url)) {
            toast.error(`${t('stream.rooms.valid.batchUrlInvalid')}${url}`);
            return;
        }
    }

    isSubmitting.value = true;

    try {
        const payload = {
            ...formValues,
            roomUrls: urls,
        };

        const res: any = await addBatchRoom(payload);
        if (res.code !== 0) {
            toast.error(res.data.msg || t('common.toast.addFailed'));
            return;
        }
        toast.success(t('stream.rooms.toast.batchSuccess'));

        isOpen.value = false;
        emit("refresh");
    } catch (error) {
        console.error("Failed to save batch rooms:", error);
    } finally {
        isSubmitting.value = false;
    }
});

const openModal = () => {
    resetForm();
    isOpen.value = true;
};

const handleClose = () => {
    isOpen.value = false;
};

defineExpose({
    openModal,
});
</script>
