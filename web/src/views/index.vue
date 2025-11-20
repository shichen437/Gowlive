<template>
    <SidebarProvider>
        <div class="flex h-screen w-full">
            <Sidebar>
                <SidebarHeader class="h-18 p-1 border-b items-center">
                    <div class="flex items-center p-4">
                        <img src="/logo.png" alt="logo" class="h-6 w-12" />
                        <h1 class="text-xl font-semibold">Gowlive</h1>
                    </div>
                </SidebarHeader>
                <SidebarContent class="flex-grow p-2">
                    <SidebarMenu>
                        <!-- 固定的主页菜单项 -->
                        <SidebarMenuItem>
                            <SidebarMenuButton :as-child="true" class="px-4 text-base">
                                <router-link to="/">
                                    <ChartNoAxesGantt class="h-5 w-5 shrink-0" />
                                    <span>概览</span>
                                </router-link>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                        <!-- 可配置的动态菜单项 -->
                        <SidebarNav :menu="configurableMenu" />
                    </SidebarMenu>
                </SidebarContent>
                <SidebarFooter class="p-4 border-t dark:border-neutral-800">
                    <SidebarMenu>
                        <SidebarMenuItem @click="showConfirmModal = true" class="w-full cursor-pointer">
                            <SidebarMenuButton
                                class="w-full justify-center text-red-600 dark:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20">
                                <LogOut class="h-8 w-8" />
                                <span>退出登录</span>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarFooter>
            </Sidebar>

            <div class="flex flex-1 flex-col">
                <Topbar />
                <main class="flex-grow p-4">
                    <div class="h-full rounded-lg p-4">
                        <router-view />
                    </div>
                </main>
            </div>
        </div>
    </SidebarProvider>
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="handleLogout" title="确认退出" description="你确定要退出登录吗？" />
</template>

<script setup lang="ts">
import { configurableMenu } from '@/lib/consts/NavItems';
import { useRouter } from 'vue-router';
import {
    SidebarProvider,
    Sidebar,
    SidebarContent,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
    SidebarFooter,
    SidebarMenuButton,
} from '@/components/ui/sidebar';
import SidebarNav from '@/components/layout/SidebarNav.vue';
import Topbar from '@/components/layout/Topbar.vue';
import { useUserStore } from '@/store/user';
import {
    ChartNoAxesGantt,
    LogOut,
} from 'lucide-vue-next';
import { ref } from 'vue';
import ConfirmModal from '@/components/modal/ConfirmModal.vue';

const userStore = useUserStore();
const router = useRouter();
const showConfirmModal = ref(false);

async function handleLogout() {
    await userStore.logout();
    router.push('/login');
}
</script>
