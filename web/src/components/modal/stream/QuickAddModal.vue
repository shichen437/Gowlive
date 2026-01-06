<template>
    <Dialog :open="open" @update:open="onOpenChange">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ $t('stream.quickAdd.title') }}</DialogTitle>
                <DialogDescription>
                    {{ $t('stream.quickAdd.desc') }}
                </DialogDescription>
            </DialogHeader>
            <form @submit="onSubmit">
                <div class="grid gap-4 py-4">
                    <FormField v-slot="{ componentField, errorMessage }" name="roomUrl">
                        <FormItem class="grid grid-cols-4 items-center gap-4">
                            <FormLabel class="text-right">
                                URL
                            </FormLabel>
                            <div class="col-span-3">
                                <FormControl>
                                    <Input id="url" type="text" v-bind="componentField"
                                        :class="{ 'border-red-500': errorMessage }" @keyup.enter="onSubmit" />
                                </FormControl>
                                <FormMessage />
                            </div>
                        </FormItem>
                    </FormField>
                </div>
                <DialogFooter>
                    <Button type="submit" :disabled="isSubmitting">
                        {{ isSubmitting ? $t('common.operation.saving') : $t('common.operation.confirm') }}
                    </Button>
                </DialogFooter>
            </form>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import { useForm } from 'vee-validate';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
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
import { Input } from '@/components/ui/input';
import { quickAdd } from '@/api/stream/live_manage';
import { toast } from 'vue-sonner';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
    open: boolean;
    onOpenChange: (open: boolean) => void;
}>();

const { t } = useI18n();
const isSubmitting = ref(false);

const formSchema = toTypedSchema(
    z.object({
        roomUrl: z.string().url({ message: t('stream.rooms.valid.validUrl') }),
    })
);

const { handleSubmit, resetForm } = useForm({
    validationSchema: formSchema,
    initialValues: {
        roomUrl: '',
    }
});

const onSubmit = handleSubmit(async (values) => {
    if (isSubmitting.value) return;

    isSubmitting.value = true;

    try {
        const response = await quickAdd(values.roomUrl);
        if (response.code === 0) {
            toast.success(t('stream.quickAdd.success'));
            props.onOpenChange(false);
            resetForm();
        } else {
            toast.error(response.msg || t('stream.quickAdd.failed'));
        }
    } catch (error: any) {
        console.log(error)
    } finally {
        isSubmitting.value = false;
    }
});
</script>
