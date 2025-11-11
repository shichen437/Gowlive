<template>
    <div class="flex justify-center p-4 md:p-6">
        <div class="w-full max-w-4xl">
            <div class="mb-8 text-center">
                <h1 class="text-3xl font-bold tracking-tight">常见问题</h1>
                <p class="text-muted-foreground mt-2">这里是一些比较常见的问题。</p>
            </div>
            <Accordion type="single" class="w-full" collapsible :default-value="faqs[0].value">
                <AccordionItem v-for="faq in faqs" :key="faq.value" :value="faq.value">
                    <AccordionTrigger class="text-md">{{
                        faq.question
                        }}</AccordionTrigger>
                    <AccordionContent>
                        <p class="leading-relaxed">
                            {{ faq.answer }}
                        </p>
                    </AccordionContent>
                </AccordionItem>
            </Accordion>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion";

const faqs = [
    {
        value: "item-0",
        question: "关于监控类型",
        answer:
            "智能监控会通过直播状态、直播历史和监控指标动态调整轮询时间。如对实时性要求较高，建议开启实时监控。如直播时间固定，建议开启定时监控。",
    },
    {
        value: "item-1",
        question: "关于监控数量",
        answer: "目前未进行压测，暂无单平台最大同时监控数的有效数据。",
    },
    {
        value: "item-2",
        question: "关于监控指标",
        answer: "监控指标会统计五分钟内的直播请求数据，出现异常指标时会进行全局告警。",
    },
    {
        value: "item-3",
        question: "直播链接格式",
        answer:
            "暂时只支持 web 端直播间链接，例如：https://live.platform.com/xxxxxx",
    },
    {
        value: "item-4",
        question: "主页链接格式",
        answer:
            "暂时只支持 web 端主播主页链接，例如：https://www.platform.com/user/xxxxxx",
    },
    {
        value: "item-5",
        question: "关于视频播放",
        answer: "Flv 格式由于原生不支持，拖拽进度条功能受限。",
    }
];
</script>
