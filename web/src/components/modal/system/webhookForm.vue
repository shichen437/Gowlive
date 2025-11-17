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
            return '选择飞书消息类型'
        case 'dingTalk':
            return '选择钉钉消息类型'
        case 'weCom':
            return '选择企业微信消息类型'
        default:
            return '选择消息类型'
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
                <FormLabel>消息类型</FormLabel>
                <Select v-bind="componentField">
                    <FormControl>
                        <SelectTrigger class="w-full">
                            <SelectValue :placeholder="messageTypePlaceholder" />
                        </SelectTrigger>
                    </FormControl>
                    <SelectContent class="w-[--radix-select-trigger-width]">
                        <SelectItem v-for="opt in messageTypeOptions" :key="opt.value" :value="opt.value">
                            {{ opt.label }}
                        </SelectItem>
                    </SelectContent>
                </Select>
                <FormMessage />
            </FormItem>
        </FormField>

        <FormField v-if="showSign" name="webhook.sign" v-slot="{ componentField }">
            <FormItem>
                <FormLabel>签名</FormLabel>
                <FormControl>
                    <Input type="password" v-bind="componentField" />
                </FormControl>
                <FormMessage />
            </FormItem>
        </FormField>
    </div>
</template>
