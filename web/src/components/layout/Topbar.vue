<template>
    <header class="flex h-18 items-center justify-between border-b bg-background px-4 md:px-6">
        <div class="flex items-center gap-4">
            <SidebarTrigger class="md:hidden" />
            <Breadcrumb>
                <BreadcrumbList class="text-lg font-medium">
                    <template v-for="(item, index) in breadcrumbs" :key="item.path">
                        <BreadcrumbItem>
                            <BreadcrumbPage v-if="index === breadcrumbs.length - 1">
                                {{ item.meta.title }}
                            </BreadcrumbPage>
                            <BreadcrumbLink as="span" v-else-if="item.children && item.children.length > 0"
                                class="cursor-default">
                                {{ item.meta.title }}
                            </BreadcrumbLink>
                            <BreadcrumbLink v-else :as-child="true">
                                <router-link :to="item.path">{{ item.meta.title }}</router-link>
                            </BreadcrumbLink>
                        </BreadcrumbItem>
                        <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
                    </template>
                </BreadcrumbList>
            </Breadcrumb>
        </div>
        <div class="flex items-center gap-4">
            <Button variant="ghost" size="icon">
                <Rss class="h-5 w-5" />
            </Button>
            <Button @click="toggleFullscreen" variant="ghost" size="icon">
                <component :is="isFullscreen ? Minimize : Maximize" class="h-5 w-5" />
            </Button>
            <Button @click="cycle()" variant="ghost" size="icon">
                <Sun class="h-5 w-5" v-if="mode === 'light'" />
                <Moon class="h-5 w-5" v-else-if="mode === 'dark'" />
                <SunMoon class="h-5 w-5" v-else />
                <span class="sr-only">Toggle theme</span>
            </Button>
            <Button variant="ghost" class="relative h-8 w-8 rounded-full">
                <router-link to="/user/index">
                    <Avatar class="h-8 w-8">
                        <AvatarFallback>{{ userInfo?.nickname?.charAt(0)?.toUpperCase() }}</AvatarFallback>
                    </Avatar>
                </router-link>
            </Button>

        </div>
    </header>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useTheme } from '@/composables/useTheme';
import { SidebarTrigger } from '@/components/ui/sidebar';
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import type { UserInfo } from '@/types/user';
import { Maximize, Minimize, Moon, Sun, SunMoon, Rss } from 'lucide-vue-next';

const { mode, cycle } = useTheme();

const userStore = useUserStore();
const userInfo: Ref<UserInfo | null> = computed(() => userStore.userInfo);
const route = useRoute();

const breadcrumbs = computed(() => {
    return route.matched.filter(item => item.meta && item.meta.title);
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

onMounted(async () => {
    if (!userStore.userInfo) {
        await userStore.getUserInfo();
    }
    document.addEventListener('fullscreenchange', onFullscreenChange);
});

onUnmounted(() => {
    document.removeEventListener('fullscreenchange', onFullscreenChange);
});
</script>
