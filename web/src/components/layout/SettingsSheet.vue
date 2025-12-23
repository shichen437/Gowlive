<template>
    <Sheet>
        <SheetTrigger as-child>
            <Button variant="ghost" size="icon">
                <Settings2 class="h-5 w-5" />
            </Button>
        </SheetTrigger>
        <SheetContent>
            <SheetHeader>
                <SheetTitle>{{ t('project.topbar.sheet.title') }}</SheetTitle>
                <SheetDescription>
                    {{ t('project.topbar.sheet.desc') }}
                </SheetDescription>
            </SheetHeader>
            <div class="grid flex-1 auto-rows-min gap-6 px-4">
                <div class="grid gap-3">
                    <Label class="font-medium">{{ t('project.topbar.sheet.fields.appearance') }}</Label>
                    <RadioGroup v-model="mode" class="grid grid-cols-3 gap-2">
                        <div>
                            <RadioGroupItem value="light" id="light" class="peer sr-only" />
                            <Label for="light"
                                class="flex h-12 w-full cursor-pointer items-center justify-center rounded-md border-2 border-muted bg-popover text-center hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                                <Sun class="h-5 w-5" />
                            </Label>
                        </div>
                        <div>
                            <RadioGroupItem value="dark" id="dark" class="peer sr-only" />
                            <Label for="dark"
                                class="flex h-12 w-full cursor-pointer items-center justify-center rounded-md border-2 border-muted bg-popover text-center hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                                <Moon class="h-5 w-5" />
                            </Label>
                        </div>
                        <div>
                            <RadioGroupItem value="auto" id="auto" class="peer sr-only" />
                            <Label for="auto"
                                class="flex h-12 w-full cursor-pointer items-center justify-center rounded-md border-2 border-muted bg-popover text-center hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                                <SunMoon class="h-5 w-5" />
                            </Label>
                        </div>
                    </RadioGroup>
                </div>
                <div class="grid gap-3">
                    <Label class="font-medium">{{ t('project.topbar.sheet.fields.theme') }}</Label>
                    <RadioGroup :model-value="theme" @update:model-value="setTheme" class="grid grid-cols-3 gap-2">
                        <div v-for="ct in themes" :key="ct.value">
                            <RadioGroupItem :value="ct.value" :id="ct.value" class="peer sr-only" />
                            <Label :for="ct.value"
                                class="flex h-12 w-full cursor-pointer items-center justify-center rounded-md border-2 border-muted bg-popover text-center hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                                {{ t(ct.name) }}
                            </Label>
                        </div>
                    </RadioGroup>
                </div>
                <div class="grid gap-3">
                    <Label class="font-medium">{{ t('project.topbar.sheet.fields.language') }}</Label>
                    <RadioGroup :model-value="locale" @update:model-value="changeLocale" class="grid grid-cols-3 gap-2">
                        <div v-for="l in locales" :key="l.locale">
                            <RadioGroupItem :value="l.locale" :id="l.locale" class="peer sr-only" />
                            <Label :for="l.locale"
                                class="flex h-12 w-full cursor-pointer items-center justify-center rounded-md border-2 border-muted bg-popover text-center hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                                {{ l.name }}
                            </Label>
                        </div>
                    </RadioGroup>
                </div>
            </div>
        </SheetContent>
    </Sheet>
</template>

<script setup lang="ts">
import { useTheme } from "@/composables/useTheme";
import { useI18n } from 'vue-i18n';
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Settings2, Sun, Moon, SunMoon } from "lucide-vue-next";
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Label } from '@/components/ui/label';

const { t, locale } = useI18n();

const locales = [
    { locale: 'en', name: 'English' },
    { locale: 'zh-CN', name: '简体中文' },
    { locale: 'zh-TW', name: '繁體中文' },
]

const changeLocale = (lang: string) => {
    if (typeof lang !== 'string') return;
    locale.value = lang;
    localStorage.setItem('locale', lang);
};

const { mode, theme, setTheme, themes } = useTheme();

</script>
