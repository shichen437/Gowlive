<template>
    <Dialog :open="open" @update:open="onOpenChange">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ t('user.profile.password.title') }}</DialogTitle>
                <DialogDescription>{{ t('user.profile.password.desc') }}</DialogDescription>
            </DialogHeader>
            <form @submit="onSubmit" class="space-y-4">
                <FormField v-slot="{ componentField }" name="oldPwd">
                    <FormItem>
                        <FormLabel>{{ t('user.profile.password.fields.oldPwd') }}</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" :placeholder="t('user.profile.password.placeholder.oldPwd')" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="newPwd">
                    <FormItem>
                        <FormLabel>{{ t('user.profile.password.fields.newPwd') }}</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" :placeholder="t('user.profile.password.placeholder.newPwd')" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="confirmPwd">
                    <FormItem>
                        <FormLabel>{{ t('user.profile.password.fields.confirmPwd') }}</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" :placeholder="t('user.profile.password.placeholder.confirmPwd')" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
            </form>
            <DialogFooter>
                <Button type="submit" @click="onSubmit">{{ t('common.operation.save') }}</Button>
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
import { putPassword } from '@/api/admin/user';
import { useUserStore } from '@/store/user';
import { toast } from 'vue-sonner';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps<{
    open: boolean;
}>();

const emit = defineEmits(['update:open']);

const userStore = useUserStore();

const formSchema = toTypedSchema(z.object({
    oldPwd: z.string().min(1, t('user.profile.password.valid.oldBlank')),
    newPwd: z.string().min(6, t('user.profile.password.valid.newPwdLength')),
    confirmPwd: z.string().min(1, t('user.profile.password.valid.confirmPwd')),
}).refine(data => data.newPwd === data.confirmPwd, {
    message: t('user.profile.password.valid.notSame'),
    path: ['confirmPwd'],
}));

const { handleSubmit, resetForm } = useForm({
    validationSchema: formSchema,
    initialValues: {
        oldPwd: '',
        newPwd: '',
        confirmPwd: '',
    }
});

watch(() => props.open, (isOpen) => {
    if (isOpen) {
        resetForm();
    }
});

const onOpenChange = (open: boolean) => {
    emit('update:open', open);
};

const onSubmit = handleSubmit(async (values) => {
    try {
        const res: any = await putPassword(values.oldPwd, values.newPwd);
        if (res.code === 0) {
            onOpenChange(false);
            userStore.logout().then(() => {
                window.location.href = '/login';
            });
        } else {
            toast.error(res.msg);
        }
    } catch (error) {
        // Error handled by interceptor
    }
});
</script>
