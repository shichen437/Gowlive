<template>
    <Dialog :open="open" @update:open="onOpenChange">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>修改密码</DialogTitle>
                <DialogDescription>请输入你的旧密码和新密码。修改成功后将会退出登录。</DialogDescription>
            </DialogHeader>
            <form @submit="onSubmit" class="space-y-4">
                <FormField v-slot="{ componentField }" name="oldPwd">
                    <FormItem>
                        <FormLabel>旧密码</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" placeholder="请输入旧密码" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="newPwd">
                    <FormItem>
                        <FormLabel>新密码</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" placeholder="请输入新密码" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="confirmPwd">
                    <FormItem>
                        <FormLabel>确认新密码</FormLabel>
                        <FormControl>
                            <Input type="password" v-bind="componentField" placeholder="请再次输入新密码" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
            </form>
            <DialogFooter>
                <Button type="submit" @click="onSubmit">保存</Button>
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

const props = defineProps<{
    open: boolean;
}>();

const emit = defineEmits(['update:open']);

const userStore = useUserStore();

const formSchema = toTypedSchema(z.object({
    oldPwd: z.string().min(1, '旧密码不能为空。'),
    newPwd: z.string().min(6, '新密码至少需要 6 个字符。'),
    confirmPwd: z.string().min(1, '请再次输入新密码。'),
}).refine(data => data.newPwd === data.confirmPwd, {
    message: '两次输入的密码不一致。',
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
