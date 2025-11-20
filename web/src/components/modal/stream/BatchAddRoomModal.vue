<template>
    <Dialog :open="isOpen" @update:open="handleClose">
        <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
                <DialogTitle>批量添加房间</DialogTitle>
            </DialogHeader>
            <form @submit="onSubmit">
                <div class="grid gap-4 py-4">
                    <FormField v-slot="{ componentField, errorMessage }" name="roomUrls">
                        <FormItem>
                            <FormLabel>直播链接</FormLabel>
                            <FormControl>
                                <Textarea placeholder="请输入直播链接，支持英文逗号或换行分隔，最多30个" v-bind="componentField"
                                    :class="{ 'border-red-500': errorMessage }" class="h-32 max-h-64 min-h-24" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>

                    <div class="grid grid-cols-2 gap-4">
                        <FormField v-slot="{ componentField, errorMessage }" name="interval">
                            <FormItem>
                                <FormLabel>间隔时间 (秒)</FormLabel>
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
                                <FormLabel>录制格式</FormLabel>
                                <Select v-bind="componentField">
                                    <FormControl>
                                        <SelectTrigger class="w-full" :class="{
                                            'border-red-500': errorMessage,
                                        }">
                                            <SelectValue placeholder="选择录制格式" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="w-[--radix-select-trigger-width]">
                                        <SelectItem value="flv">FLV</SelectItem>
                                        <SelectItem value="mkv">MKV</SelectItem>
                                        <SelectItem value="ts">TS</SelectItem>
                                        <SelectItem value="mp4">MP4</SelectItem>
                                        <SelectItem value="mp3">MP3(仅音频)</SelectItem>
                                    </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                    </div>

                    <FormField v-slot="{ componentField, errorMessage }" name="remark">
                        <FormItem>
                            <FormLabel>备注</FormLabel>
                            <FormControl>
                                <Input type="text" placeholder="请输入备注" v-bind="componentField"
                                    :class="{ 'border-red-500': errorMessage }" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>
                </div>
                <DialogFooter>
                    <Button type="submit" :disabled="isSubmitting">
                        {{ isSubmitting ? "保存中..." : "保存" }}
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

const isOpen = ref(false);
const emit = defineEmits(["refresh"]);
const isSubmitting = ref(false);

const formSchema = toTypedSchema(
    z.object({
        roomUrls: z.string().min(1, { message: "请输入直播链接" }),
        interval: z.coerce
            .number()
            .min(30, { message: "间隔时间最小为30秒" })
            .max(600, { message: "间隔时间最大为600秒" }),
        format: z.enum(["flv", "mp4", "mp3", "mkv", "ts"]),
        remark: z
            .string()
            .max(45, { message: "备注最长为45个字符" })
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
        toast.error("请输入有效的直播链接");
        return;
    }
    if (urls.length > 30) {
        toast.error("最多支持30个直播链接");
        return;
    }

    const urlRegex =
        /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/;
    for (const url of urls) {
        if (!urlRegex.test(url)) {
            toast.error(`链接格式无效: ${url}`);
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
            toast.error(res.data.msg || "添加失败");
            return;
        }
        toast.success("添加成功,请稍后查看");

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
