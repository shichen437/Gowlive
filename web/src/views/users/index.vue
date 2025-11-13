<template>
    <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- Left Navigation -->
        <div class="md:col-span-1">
            <nav class="space-y-1 py-4">
                <Button variant="ghost" class="w-full justify-start"
                    :class="{ 'bg-accent text-accent-foreground': activeView === 'account' }"
                    @click="activeView = 'account'">
                    账户信息
                </Button>
                <Button variant="ghost" class="w-full justify-start"
                    :class="{ 'bg-accent text-accent-foreground': activeView === 'live-settings' }"
                    @click="activeView = 'live-settings'">
                    直播设置
                </Button>
            </nav>
        </div>

        <!-- Right Content -->
        <div class="md:col-span-3">
            <div v-if="activeView === 'account'">
                <Card>
                    <!-- 右侧内容区卡片：互换左右列 -->
                    <CardHeader class="grid grid-cols-1 md:grid-cols-4 gap-4 items-center">
                        <!-- 原右列（头像/按钮）改为左列 -->
                        <div class="md:col-span-1 flex flex-col items-center justify-center space-y-4">
                            <Avatar class="h-24 w-24 text-3xl">
                                <AvatarFallback>{{ userInfo?.nickname?.charAt(0)?.toUpperCase() }}</AvatarFallback>
                            </Avatar>
                            <div class="flex gap-4 justify-center">
                                <Button @click="showProfileModal = true">修改信息</Button>
                                <Button variant="outline" @click="showPwdModal = true">修改密码</Button>
                            </div>
                        </div>

                        <!-- 原左列（信息列表）改为右列 -->
                        <div class="md:col-span-3 space-y-4">
                            <div class="flex items-center">
                                <span class="w-20 text-sm text-muted-foreground">用户名</span>
                                <div class="ml-auto text-right">
                                    <p class="text-muted-foreground text-right">{{ userInfo?.username }}</p>
                                </div>
                            </div>
                            <div class="flex items-center">
                                <span class="w-20 text-sm text-muted-foreground">昵称</span>
                                <p class="ml-auto">{{ userInfo?.nickname }}</p>
                            </div>
                            <div class="flex items-center">
                                <span class="w-20 text-sm text-muted-foreground">性别</span>
                                <div class="ml-auto text-right">
                                    <Mars v-if="userInfo?.sex === 1" class="h-5 w-5 text-blue-500 inline-block" />
                                    <Venus v-else class="h-5 w-5 text-pink-500 inline-block" />
                                </div>
                            </div>
                            <div class="flex items-center">
                                <span class="w-20 text-sm text-muted-foreground">状态</span>
                                <Badge class="ml-auto" :variant="userInfo?.status === 1 ? 'default' : 'destructive'">
                                    {{ formattedStatus }}
                                </Badge>
                            </div>
                        </div>
                    </CardHeader>
                </Card>
            </div>
            <div v-if="activeView === 'live-settings'">
                <LiveSettings />
            </div>
        </div>
    </div>
    <UpdateProfileModal v-model:open="showProfileModal" :user-info="userInfo" @success="onProfileUpdateSuccess" />
    <UpdatePwdModal v-model:open="showPwdModal" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref, type Ref } from 'vue';
import { useUserStore } from '@/store/user';
import type { UserInfo } from '@/types/user';
import {
    Card,
    CardHeader,
} from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Mars, Venus } from 'lucide-vue-next';
import UpdateProfileModal from '@/components/modal/admin/UpdateProfileModal.vue';
import UpdatePwdModal from '@/components/modal/admin/UpdatePwdModal.vue';
import LiveSettings from './LiveSettings.vue';

const userStore = useUserStore();
const userInfo: Ref<UserInfo | null> = computed(() => userStore.userInfo);

const showProfileModal = ref(false);
const showPwdModal = ref(false);
const activeView = ref('account');

const formattedStatus = computed(() => {
    if (userInfo.value?.status === 1) return '活跃';
    return '禁用';
});

const onProfileUpdateSuccess = async () => {
    await userStore.getUserInfo(true);
};

onMounted(async () => {
    if (!userStore.userInfo) {
        await userStore.getUserInfo();
    }
});
</script>
