<template>
    <div class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Left Column -->
            <div class="lg:col-span-1">
                <Card>
                    <CardHeader class="flex flex-col items-center gap-4 text-center">
                        <Avatar class="h-24 w-24 text-3xl">
                            <AvatarFallback>{{ userInfo?.nickname?.charAt(0)?.toUpperCase() }}</AvatarFallback>
                        </Avatar>
                        <div>
                            <h2 class="text-xl font-semibold">{{ userInfo?.nickname }}</h2>
                            <p class="text-sm text-muted-foreground">{{ userInfo?.username }}</p>
                        </div>
                    </CardHeader>
                    <CardContent>
                        <ul class="space-y-3">
                            <li class="flex items-center justify-between">
                                <span class="text-sm font-medium text-muted-foreground">性别</span>
                                <div>
                                    <Mars v-if="userInfo?.sex === 1" class="h-5 w-5 text-blue-500" />
                                    <Venus v-else class="h-5 w-5 text-pink-500" />
                                </div>
                            </li>
                            <li class="flex items-center justify-between">
                                <span class="text-sm font-medium text-muted-foreground">状态</span>
                                <Badge :variant="userInfo?.status === 1 ? 'default' : 'destructive'">
                                    {{ formattedStatus }}
                                </Badge>
                            </li>
                        </ul>
                    </CardContent>
                    <CardFooter class="grid grid-cols-2 gap-4">
                        <Button @click="showProfileModal = true">修改信息</Button>
                        <Button variant="outline" @click="showPwdModal = true">修改密码</Button>
                    </CardFooter>
                </Card>
            </div>
            <!-- Right Column -->
            <div class="lg:col-span-2">
                <!-- Content for the right column can be added here in the future -->
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
    CardContent,
    CardFooter,
    CardHeader,
} from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Mars, Venus } from 'lucide-vue-next';
import UpdateProfileModal from '@/components/modal/admin/UpdateProfileModal.vue';
import UpdatePwdModal from '@/components/modal/admin/UpdatePwdModal.vue';

const userStore = useUserStore();
const userInfo: Ref<UserInfo | null> = computed(() => userStore.userInfo);

const showProfileModal = ref(false);
const showPwdModal = ref(false);

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