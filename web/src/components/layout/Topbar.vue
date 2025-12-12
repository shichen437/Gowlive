<template>
    <header class="flex h-18 items-center justify-between border-b bg-background px-4 md:px-6">
        <div class="flex items-center gap-4">
            <SidebarTrigger class="md:hidden" />
            <Breadcrumb>
                <BreadcrumbList class="text-lg font-medium">
                    <template v-for="(item, index) in breadcrumbs" :key="item.path">
                        <BreadcrumbItem>
                            <BreadcrumbPage v-if="index === breadcrumbs.length - 1">
                                {{ t(item.meta.title as string) }}
                            </BreadcrumbPage>
                            <BreadcrumbLink as="span" v-else-if="
                                item.children && item.children.length > 0
                            " class="cursor-default">
                                {{ t(item.meta.title as string) }}
                            </BreadcrumbLink>
                            <BreadcrumbLink v-else :as-child="true">
                                <router-link :to="item.path">{{
                                    t(item.meta.title as string)
                                    }}</router-link>
                            </BreadcrumbLink>
                        </BreadcrumbItem>
                        <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
                    </template>
                </BreadcrumbList>
            </Breadcrumb>
        </div>
        <div class="flex items-center gap-4">
            <TooltipProvider v-if="isUnhealthy">
                <Tooltip>
                    <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon">
                            <Frown class="h-5 w-5 text-destructive" />
                        </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                        <div class="grid gap-2">
                            <p>
                                <span class="mr-2">{{ t('project.topbar.errRate') }}:</span>
                                <span>{{
                                    healthInfo.errorPercent.toFixed(2)
                                }}%</span>
                            </p>
                            <p>
                                <span class="mr-2">{{ t('project.topbar.diskUsage') }}:</span>
                                <span>{{
                                    healthInfo.diskUsage.toFixed(2)
                                }}%</span>
                            </p>
                        </div>
                    </TooltipContent>
                </Tooltip>
            </TooltipProvider>
            <Button @click="toggleFullscreen" variant="ghost" size="icon">
                <component :is="isFullscreen ? Minimize : Maximize" class="h-5 w-5" />
            </Button>
            <DropdownMenu>
                <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon">
                        <Languages class="h-5 w-5" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuItem v-for="l in locales" :key="l.locale" @click="changeLocale(l.locale)">
                        {{ l.name }}
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
            <Button @click="cycle()" variant="ghost" size="icon">
                <Sun class="h-5 w-5" v-if="mode === 'light'" />
                <Moon class="h-5 w-5" v-else-if="mode === 'dark'" />
                <SunMoon class="h-5 w-5" v-else />
                <span class="sr-only">Toggle theme</span>
            </Button>
            <HoverCard v-if="Object.keys(features).length > 0">
                <HoverCardTrigger as-child>
                    <Button variant="ghost" size="icon">
                        <Megaphone class="h-5 w-5" />
                    </Button>
                </HoverCardTrigger>
                <HoverCardContent class="w-80">
                    <div class="flex justify-between space-x-4">
                        <div class="space-y-1">
                            <h4 class="text-lg font-semibold">
                                v{{ features.version }}
                            </h4>
                            <div v-for="(category, key) in features.info" :key="key" class="text-sm">
                                <p class="font-medium pt-2" v-if="category.items.length > 0">
                                    {{ category.name }}
                                </p>
                                <ul class="list-disc list-inside">
                                    <li v-for="item in category.items" :key="item.desc">
                                        {{ item.desc }}
                                    </li>
                                </ul>
                            </div>
                            <div class="flex pt-2">
                                <Calendar class="mr-2 h-4 w-4 opacity-70" />
                                <span class="text-xs text-muted-foreground">
                                    {{ t('project.topbar.publish') }} {{ features.date }}
                                </span>
                            </div>
                        </div>
                    </div>
                </HoverCardContent>
            </HoverCard>
            <Button variant="ghost" class="relative h-8 w-8 rounded-full">
                <router-link to="/user/index">
                    <Avatar class="h-8 w-8">
                        <AvatarFallback>{{
                            userInfo?.nickname?.charAt(0)?.toUpperCase()
                            }}</AvatarFallback>
                    </Avatar>
                </router-link>
            </Button>
        </div>
    </header>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, reactive, type Ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/store/user";
import { useTheme } from "@/composables/useTheme";
import { createSSEConnection } from "@/lib/sse";
import { useI18n } from 'vue-i18n';
import { toast } from 'vue-sonner';
import { SidebarTrigger } from "@/components/ui/sidebar";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";
import type { UserInfo } from "@/types/user";
import {
    Maximize,
    Minimize,
    Moon,
    Sun,
    SunMoon,
    Megaphone,
    Calendar,
    Frown,
    Languages,
} from "lucide-vue-next";
import features from "@/lib/others/feature.json";

const { t, locale } = useI18n();

const locales = [
    { locale: 'en', name: 'English' },
    { locale: 'zh-CN', name: '简体中文' },
    { locale: 'zh-TW', name: '繁體中文' },
]

const changeLocale = (lang: string) => {
    locale.value = lang;
    localStorage.setItem('locale', lang);
};

const { mode, cycle } = useTheme();

const userStore = useUserStore();
const userInfo: Ref<UserInfo | null> = computed(() => userStore.userInfo);
const route = useRoute();
const router = useRouter()

const breadcrumbs = computed(() => {
    return route.matched.filter((item) => item.meta && item.meta.title);
});

const isFullscreen = ref(false);

const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
    } else if (document.exitFullscreen) {
        document.exitFullscreen();
    }
};

const onFullscreenChange = () => {
    isFullscreen.value = !!document.fullscreenElement;
};

let sseClient: any = null;
const healthInfo = reactive({
    errorPercent: 0,
    diskUsage: 0,
});
const isUnhealthy = computed(
    () => healthInfo.errorPercent > 1 || healthInfo.diskUsage >= 85,
);

onMounted(async () => {
    if (!userStore.userInfo) {
        await userStore.getUserInfo();
    }
    document.addEventListener("fullscreenchange", onFullscreenChange);

    sseClient = createSSEConnection({
        channel: "global",
        onMessage: (msg) => {
            if (msg.event === "global") {
                if (msg.data.health) {
                    healthInfo.errorPercent = msg.data.health.errorPercent;
                    healthInfo.diskUsage = msg.data.health.diskUsage;
                }
                if (msg.data.notify) {
                    if (msg.data.notify.level === 'info') {
                        toast.info(msg.data.notify.title, { position: 'top-center', description: msg.data.notify.content });
                    } else {
                        toast.error(msg.data.notify.title, {
                            position: 'top-center', duration: 7000, action: {
                                label: t('project.topbar.detail'),
                                onClick: () => {
                                    router.push('/system/notify')
                                }
                            }
                        });
                    }
                }
            }
        },
        onError: (error) => {
            console.error("Health SSE error:", error);
        },
    });
});

onUnmounted(() => {
    document.removeEventListener("fullscreenchange", onFullscreenChange);
    if (sseClient) {
        sseClient.disconnect();
    }
});
</script>
