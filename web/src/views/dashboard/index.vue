<template>
    <div class="flex flex-col h-full">
        <div>
            <h3 class="text-xl font-bold tracking-tight mb-4">
                {{ greeting }}
            </h3>
            <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">房间总数</CardTitle>
                        <Home class="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">{{ overviewData?.liveRoomCount ?? 0 }}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">正在录制</CardTitle>
                        <MonitorDot class="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">{{ overviewData?.recordingRoomCount ?? 0 }}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">录制时长</CardTitle>
                        <Clock class="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">{{ overviewData?.recordTimeCount ?? 0 }}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">未读通知</CardTitle>
                        <Bell class="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">{{ overviewData?.unreadMessageCount ?? 0 }}</div>
                    </CardContent>
                </Card>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from '@/components/ui/card';
import { Home, Clock, MonitorDot, Bell } from 'lucide-vue-next';
import { onMounted, ref, computed } from 'vue';
import { overview as getOverview } from '@/api/system/overview';
import type { Overview } from '@/types/overview';
import { useUserStore } from '@/store/user';
import { getGreeting } from '@/utils/greeting';

const overviewData = ref<Overview>();

const userStore = useUserStore();

const greeting = computed(() => {
    return getGreeting(userStore.userInfo?.nickname ?? '');
});

onMounted(async () => {
    const res: any = await getOverview();
    if (res.code === 0) {
        overviewData.value = res.data.data;
    }
    if (!userStore.userInfo) {
        userStore.getUserInfo();
    }
});
</script>