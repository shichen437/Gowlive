<template>
    <Dialog v-model:open="open">
        <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
                <DialogTitle>{{ isEdit ? t('common.fields.edit') : t('common.fields.add') }}{{
                    t('system.channel.fields.channel.title') }}</DialogTitle>
            </DialogHeader>
            <form @submit="onSubmit" class="space-y-4">
                <div class="grid grid-cols-2 gap-4">
                    <FormField name="name" v-slot="{ componentField }">
                        <FormItem>
                            <FormLabel>{{ t('system.channel.fields.name') }}</FormLabel>
                            <FormControl>
                                <Input type="text" :placeholder="t('system.channel.placeholder.name')"
                                    v-bind="componentField" />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>
                    <FormField name="status" v-slot="{ value, handleChange }">
                        <FormItem>
                            <FormLabel>{{ t('system.channel.fields.status') }}</FormLabel>
                            <FormControl>
                                <div class="flex gap-2">
                                    <Button type="button" class="flex-1" :variant="value === 1 ? 'default' : 'outline'"
                                        @click="handleChange(1)">
                                        {{ t('system.channel.fields.statusActive') }}
                                    </Button>
                                    <Button type="button" class="flex-1" :variant="value === 0 ? 'default' : 'outline'"
                                        @click="handleChange(0)">
                                        {{ t('system.channel.fields.statusDisabled') }}
                                    </Button>
                                </div>
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    </FormField>
                </div>

                <Tabs :model-value="values.type" @update:model-value="(val) => handleTypeChange(String(val))"
                    default-value="email" class="w-full">
                    <TabsList class="grid w-full grid-cols-5" :class="{ 'pointer-events-none opacity-50': isEdit }">
                        <TabsTrigger value="email"> {{ t('system.channel.fields.channel.email') }} </TabsTrigger>
                        <TabsTrigger value="gotify"> Gotify </TabsTrigger>
                        <TabsTrigger value="lark"> {{ t('system.channel.fields.channel.lark') }} </TabsTrigger>
                        <TabsTrigger value="dingTalk"> {{ t('system.channel.fields.channel.dingTalk') }} </TabsTrigger>
                        <TabsTrigger value="weCom"> {{ t('system.channel.fields.channel.weCom') }} </TabsTrigger>
                    </TabsList>
                    <TabsContent value="email" class="space-y-4 mt-4 border-0 p-0">
                        <div class="grid grid-cols-2 gap-4">
                            <FormField name="email.server" v-slot="{ componentField }">
                                <FormItem>
                                    <FormLabel>{{ t('system.channel.fields.smtpServer') }}</FormLabel>
                                    <FormControl>
                                        <Input type="text" placeholder="smtp.example.com" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                            <FormField name="email.port" v-slot="{ componentField }">
                                <FormItem>
                                    <FormLabel>{{ t('system.channel.fields.port') }}</FormLabel>
                                    <FormControl>
                                        <Input type="number" placeholder="465" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                        </div>
                        <FormField name="email.sender" v-slot="{ componentField }">
                            <FormItem>
                                <FormLabel>{{ t('system.channel.fields.sender') }}</FormLabel>
                                <FormControl>
                                    <Input type="email" placeholder="sender@example.com" v-bind="componentField" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                        <FormField name="email.authCode" v-slot="{ componentField }">
                            <FormItem>
                                <FormLabel>{{ t('system.channel.fields.authCode') }}</FormLabel>
                                <FormControl>
                                    <Input type="password" :placeholder="t('system.channel.placeholder.authCode')"
                                        v-bind="componentField" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                        <FormField name="email.receiver" v-slot="{ componentField }">
                            <FormItem>
                                <FormLabel>{{ t('system.channel.fields.receiver') }}</FormLabel>
                                <FormControl>
                                    <Input type="email" placeholder="receiver@example.com" v-bind="componentField" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                    </TabsContent>
                    <TabsContent value="gotify" class="space-y-4 mt-4 border-0 p-0">
                        <FormField name="url" v-slot="{ componentField }">
                            <FormItem>
                                <FormLabel>Gotify URL</FormLabel>
                                <FormControl>
                                    <Input type="text" placeholder="https://gotify.example.com"
                                        v-bind="componentField" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        </FormField>
                    </TabsContent>
                    <TabsContent value="lark">
                        <WebhookForm platform="lark" :messageTypeOptions="larkOptions" :showSign="true" />
                    </TabsContent>

                    <TabsContent value="dingTalk">
                        <WebhookForm platform="dingTalk" :messageTypeOptions="basicOptions" :showSign="true" />
                    </TabsContent>

                    <TabsContent value="weCom">
                        <WebhookForm platform="weCom" :messageTypeOptions="basicOptions" :showSign="false" />
                    </TabsContent>
                </Tabs>

                <FormField name="remark" v-slot="{ componentField }">
                    <FormItem>
                        <FormLabel>{{ t('common.fields.remark') }}</FormLabel>
                        <FormControl>
                            <Input :placeholder="t('system.channel.placeholder.remark')" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>

                <DialogFooter>
                    <Button type="button" variant="outline" @click="open = false">{{ t('common.operation.cancel')
                        }}</Button>
                    <Button type="submit">{{ t('common.operation.save') }}</Button>
                </DialogFooter>
            </form>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { watch, computed } from "vue";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import {
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    addPushChannel,
    updatePushChannel,
    getPushChannel,
} from "@/api/system/push_channel";
import type { PushChannel } from "@/types/system";
import { toast } from "vue-sonner";
import WebhookForm from './webhookForm.vue'
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const larkOptions = [
    { value: 0, label: '文本' },
    { value: 1, label: '富文本' },
    { value: 2, label: '卡片' },
]
const basicOptions = [
    { value: 0, label: '文本' },
    { value: 1, label: '富文本' },
]

const getInitialValues = () => ({
    status: 1,
    type: 'email',
    remark: '',
    url: '',
    email: {
        sender: '',
        receiver: '',
        server: '',
        port: undefined,
        authCode: '',
    },
    webhook: {
        webhookUrl: '',
        messageType: 0,
        sign: '',
    },
})

const props = defineProps<{
    modelValue: boolean;
    channel: PushChannel | null;
}>();

const emit = defineEmits(["update:modelValue", "success"]);

const open = computed({
    get: () => props.modelValue,
    set: (value) => emit("update:modelValue", value),
});

const isEdit = computed(() => !!props.channel);

const formSchema = toTypedSchema(
    z
        .object({
            name: z.string().min(1, t('system.channel.valid.nameNotEmpty')),
            type: z.string().min(1, t('system.channel.valid.typeNotEmpty')),
            status: z.coerce.number(),
            remark: z.string().optional(),
            url: z.string().optional(),
            email: z
                .object({
                    sender: z.string().optional(),
                    receiver: z.string().optional(),
                    server: z.string().optional(),
                    port: z.coerce.number().optional(),
                    authCode: z.string().optional(),
                })
                .optional(),
            webhook: z
                .object({
                    webhookUrl: z.string().optional(),
                    messageType: z.coerce.number().optional(),
                    sign: z.string().optional(),
                })
                .optional(),
        })
);

const { handleSubmit, resetForm, values, setFieldValue, setValues, setFieldError } = useForm({
    validationSchema: formSchema,
});

watch(open, async (isOpen) => {
    if (isOpen) {
        if (isEdit.value && props.channel) {
            try {
                const res: any = await getPushChannel(props.channel.id);
                if (res.code !== 0) {
                    toast.error(res.msg || t('system.channel.toast.detailErr'))
                }
                setValues({
                    ...res.data,
                    email: res.data.email || {},
                    webhook: res.data.webhook || {},
                });
            } catch (error) {
                console.error("Failed to fetch channel details:", error);
            }
        } else {
            resetForm({ values: getInitialValues() })
        }
    } else {
        resetForm({ values: getInitialValues() })
    }
});

const onSubmit = handleSubmit(async (formValues) => {
    try {
        if (isEdit.value && props.channel) {
            await updatePushChannel({ ...formValues, id: props.channel.id });
        } else {
            await addPushChannel(formValues);
        }
        emit("success");
        open.value = false;
    } catch (error) {
        console.error("Failed to save channel:", error);
    }
});

const handleTypeChange = (newType: string) => {
    setFieldValue('type', newType);
    setValues({
        ...values,
        url: '',
        email: getInitialValues().email,
        webhook: getInitialValues().webhook,
    })
    setFieldError('url', undefined)
};
</script>
