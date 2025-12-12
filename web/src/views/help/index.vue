<template>
    <div class="flex justify-center p-4 md:p-6">
        <div class="w-full max-w-4xl">
            <div class="mb-8 text-center">
                <h1 class="text-3xl font-bold tracking-tight">{{ t('project.help.title') }}</h1>
                <p class="text-muted-foreground mt-2">{{ t('project.help.desc') }}</p>
            </div>
            <Accordion v-if="faqs.length > 0" type="single" class="w-full" collapsible :default-value="faqs[0].value">
                <AccordionItem v-for="faq in faqs" :key="faq.value" :value="faq.value">
                    <AccordionTrigger class="text-md">{{ faq.question }}</AccordionTrigger>
                    <AccordionContent>
                        <p class="leading-relaxed whitespace-pre-line">
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
import { useI18n } from 'vue-i18n';
import { computed } from 'vue';

const { t, locale, messages } = useI18n();

interface Faq {
    value: string;
    question: string;
    answer: string;
}

const faqs = computed<Faq[]>(() => {
    try {
        const faqsData = (messages.value[locale.value] as any).project.help.faqs;
        return Array.isArray(faqsData) ? faqsData : [];
    } catch (e) {
        console.error("Could not load FAQs:", e);
        return [];
    }
});
</script>
