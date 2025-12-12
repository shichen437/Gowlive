<script setup lang="ts">
import { computed } from 'vue'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface MessageTypeOption {
    value: number
    label: string
}
type Platform = 'lark' | 'dingTalk' | 'weCom'
const props = defineProps<{
    platform: Platform
    messageTypeOptions: MessageTypeOption[]
    showSign?: boolean
}>()

const messageTypePlaceholder = computed(() => {
    switch (props.platform) {
        case 'lark':
            return t('system.channel.placeholder.larkMsgType')
        case 'dingTalk':
            return t('system.channel.placeholder.dingMsgType')
        case 'weCom':
            return t('system.channel.placeholder.weComMsgType')
        default:
            return t('system.channel.placeholder.msgType')
    }
})
</script>

<template>
    <div class="space-y-4 mt-4 border-0 p-0">
        <FormField name="webhook.webhookUrl" v-slot="{ componentField }">
            <FormItem>
                <FormLabel>Webhook</FormLabel>
                <FormControl>
                    <Input type="text" v-bind="componentField" />
                </FormControl>
                <FormMessage />
            </FormItem>
        </FormField>

        <FormField name="webhook.messageType" v-slot="{ componentField }">
            <FormItem>
                <FormLabel>{{ t('system.channel.fields.messageType') }}</FormLabel>
                <Select v-bind="componentField">
                    <FormControl>
                        <SelectTrigger class="w-full">
                            <SelectValue :placeholder="messageTypePlaceholder" />
                        </SelectTrigger>
                    </FormControl>
                    <SelectContent class="w-[--radix-select-trigger-width]">
                        <SelectItem v-for="opt in messageTypeOptions" :key="opt.value" :value="opt.value">
                            {{ opt.label === '文本' ? t('system.channel.fields.text') : opt.label === '富文本' ? t('system.channel.fields.richtext') : opt.label === '卡片' ? t('system.channel.fields.card') : opt.label }}
                        </SelectItem>
                    </SelectContent>
                </Select>
                <FormMessage />
            </FormItem>
        </FormField>

        <FormField v-if="showSign" name="webhook.sign" v-slot="{ componentField }">
            <FormItem>
                <FormLabel>{{ t('system.channel.fields.sign') }}</FormLabel>
                <FormControl>
                    <Input type="password" v-bind="componentField" />
                </FormControl>
                <FormMessage />
            </FormItem>
        </FormField>
    </div>
</template>
