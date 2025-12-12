<template>
  <div class="flex flex-col h-full">
    <div>
      <h3 class="text-xl font-bold tracking-tight mb-4">
        {{ greeting }}
      </h3>
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <router-link to="/stream/index">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">{{ t('project.dashboard.cardTitle.total') }}</CardTitle>
              <Home class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ overviewData?.liveRoomCount ?? 0 }}
              </div>
            </CardContent>
          </router-link>
        </Card>
        <Card>
          <router-link to="/stream/index">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">{{ t('project.dashboard.cardTitle.recording') }}</CardTitle>
              <MonitorDot class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ overviewData?.recordingRoomCount ?? 0 }}
              </div>
            </CardContent>
          </router-link>
        </Card>
        <Card>
          <router-link to="/stream/history">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">{{ t('project.dashboard.cardTitle.duration') }}</CardTitle>
              <Clock class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ overviewData?.recordTimeCount ?? 0 }}
              </div>
            </CardContent>
          </router-link>
        </Card>
        <Card>
          <router-link to="/system/notify">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">{{ t('project.dashboard.cardTitle.unread') }}</CardTitle>
              <Bell class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ overviewData?.unreadMessageCount ?? 0 }}
              </div>
            </CardContent>
          </router-link>
        </Card>
      </div>
      <div class="grid gap-2 mt-4">
        <Card>
          <CardContent class="grid gap-4 md:grid-cols-3">
            <div>
              <div class="flex items-center space-x-2 mb-2">
                <Cpu class="h-4 w-4 text-muted-foreground" />
                <span class="font-semibold">CPU</span>
              </div>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.cpu.model') }}:</span> <span>{{ monitorInfo.cpu.modelName }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.cpu.cores') }}:</span> <span>{{ monitorInfo.cpu.cores }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.cpu.frequency') }}:</span> <span>{{ monitorInfo.cpu.mhz }} MHz</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.cpu.usage') }}:</span>
                <span>{{ monitorInfo.cpu.percent.toFixed(2) }}%</span>
              </p>
            </div>
            <div>
              <div class="flex items-center space-x-2 mb-2">
                <MemoryStick class="h-4 w-4 text-muted-foreground" />
                <span class="font-semibold">{{ t('project.dashboard.memory') }}</span>
              </div>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.total') }}:</span>
                <span>{{ formatBytes(monitorInfo.mem.total) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.used') }}:</span>
                <span>{{ formatBytes(monitorInfo.mem.used) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.available') }}:</span>
                <span>{{ formatBytes(monitorInfo.mem.available) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.usage') }}:</span>
                <span>{{ monitorInfo.mem.usedPercent.toFixed(2) }}%</span>
              </p>
            </div>
            <div>
              <div class="flex items-center space-x-2 mb-2">
                <HardDrive class="h-4 w-4 text-muted-foreground" />
                <span class="font-semibold">{{ t('project.dashboard.disk') }} ({{ monitorInfo.disk.fstype != "" ? monitorInfo.disk.fstype : "N/A" }})</span>
              </div>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.total') }}:</span>
                <span>{{ formatBytes(monitorInfo.disk.total) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.used') }}:</span>
                <span>{{ formatBytes(monitorInfo.disk.used) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.available') }}:</span>
                <span>{{ formatBytes(monitorInfo.disk.free) }}</span>
              </p>
              <p class="text-sm text-muted-foreground flex justify-between">
                <span>{{ t('project.dashboard.common.usage') }}:</span>
                <span :class="monitorInfo.disk.usedPercent > 90
                  ? 'text-rose-600'
                  : monitorInfo.disk.usedPercent >= 70
                    ? 'text-amber-600'
                    : ''
                  ">{{ monitorInfo.disk.usedPercent.toFixed(2) }}%</span>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Home,
  Clock,
  MonitorDot,
  Bell,
  Cpu,
  MemoryStick,
  HardDrive,
} from "lucide-vue-next";
import { onMounted, onUnmounted, ref, computed } from "vue";
import { overview as getOverview } from "@/api/system/overview";
import type { Overview, MonitorInfo } from "@/types/overview";
import { useUserStore } from "@/store/user";
import { getGreetingTime } from "@/utils/greeting";
import { formatBytes } from "@/utils/convert";
import { createSSEConnection } from "@/lib/sse";
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const overviewData = ref<Overview>();
const monitorInfo = ref<MonitorInfo>({
  cpu: { cpu: 0, cores: 0, modelName: "", mhz: 0, percent: 1 },
  mem: { total: 0, used: 0, available: 0, usedPercent: 1 },
  disk: { path: "", fstype: "", total: 0, free: 0, used: 0, usedPercent: 1 },
});

const userStore = useUserStore();

let sseClient: any = null;

const greeting = computed(() => {
  const time = getGreetingTime();
  const greetingKey = `project.dashboard.greeting.${time}`;

  if (time === 'night') {
    return t(greetingKey);
  }
  return t(greetingKey, { nickname: userStore.userInfo?.nickname ?? "" });
});

onMounted(async () => {
  const res: any = await getOverview();
  if (res.code === 0) {
    overviewData.value = res.data.data;
  }
  if (!userStore.userInfo) {
    userStore.getUserInfo();
  }

  sseClient = createSSEConnection({
    channel: "monitor",
    onMessage: (msg) => {
      if (msg.event === "monitor") {
        monitorInfo.value = msg.data;
      }
    },
    onError: (error) => {
      console.error("Monitor SSE error:", error);
    },
  });
});

onUnmounted(() => {
  if (sseClient) {
    sseClient.disconnect();
  }
});
</script>