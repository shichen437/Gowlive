<template>
    <Dialog :open="open" @update:open="onOpenChange">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ t('user.profile.edit.title') }}</DialogTitle>
                <DialogDescription>{{ t('user.profile.edit.desc') }}</DialogDescription>
            </DialogHeader>
            <form class="space-y-4">
                <FormField v-slot="{ componentField }" name="nickname">
                    <FormItem class="grid grid-cols-4 items-center gap-4">
                        <FormLabel class="text-right">{{ t('user.profile.fileds.nickname') }}</FormLabel>
                        <FormControl class="col-span-3">
                            <Input v-bind="componentField" />
                        </FormControl>
                        <FormMessage class="col-start-2 col-span-3" />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="sex">
                    <FormItem class="grid grid-cols-4 items-center gap-4">
                        <FormLabel class="text-right">{{ t('user.profile.fileds.gender') }}</FormLabel>
                        <Select v-bind="componentField">
                            <FormControl class="col-span-3">
                                <SelectTrigger class="w-full">
                                    <SelectValue :placeholder="t('user.profile.edit.placeholder.sex')" />
                                </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                                <SelectItem :value="1">{{ t('user.profile.fileds.male') }}</SelectItem>
                                <SelectItem :value="0">{{ t('user.profile.fileds.female') }}</SelectItem>
                            </SelectContent>
                        </Select>
                        <FormMessage class="col-start-2 col-span-3" />
                    </FormItem>
                </FormField>
            </form>
            <DialogFooter>
                <Button type="button" @click="onSubmit">{{ t('common.operation.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import { useForm } from 'vee-validate';
import { watch } from 'vue';
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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { putProfile } from '@/api/admin/user';
import type { UserInfo } from '@/types/user';
import { toast } from 'vue-sonner';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps<{
    open: boolean;
    userInfo: UserInfo | null;
}>();

const emit = defineEmits(['update:open', 'success']);

const formSchema = toTypedSchema(z.object({
    nickname: z.string().min(2, t('user.profile.edit.valid.nickname')),
    sex: z.coerce.number().int().nonnegative().max(1),
}));

const { handleSubmit, setValues } = useForm({
    validationSchema: formSchema,
});

watch([() => props.open, () => props.userInfo], ([isOpen, userInfo]) => {
    if (isOpen && userInfo) {
        setValues({
            nickname: userInfo.nickname,
            sex: userInfo.sex,
        });
    }
}, { immediate: true });

const onOpenChange = (open: boolean) => {
    emit('update:open', open);
};

const onSubmit = handleSubmit((values) => {
    const promise = putProfile(values.nickname, values.sex);
    toast.promise(promise, {
        loading: t('common.operation.saving'),
        success: () => {
            emit('success');
            onOpenChange(false);
            return t('common.toast.updateSuccess');
        },
        error: t('common.toast.updateFailed'),
    });
});
</script>
