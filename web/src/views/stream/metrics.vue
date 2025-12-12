<template>
  <div class="space-y-4">
    <div v-if="!metricsData || metricsData.size === 0" class="grid gap-4 mt-4">
      <Card>
        <CardContent>
          <Empty class="border-none p-0 md:p-0">
            <EmptyMedia variant="icon">
              <SquareActivity class="size-8" />
            </EmptyMedia>
            <EmptyHeader>
              <EmptyTitle>{{ t('common.noData') }}</EmptyTitle>
              <EmptyDescription> {{ t('stream.metrics.emptyDesc') }} </EmptyDescription>
            </EmptyHeader>
          </Empty>
        </CardContent>
      </Card>
    </div>
    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card
        v-for="[platform, metric] in Array.from(metricsData.entries())"
        :key="platform"
      >
        <CardHeader
          class="flex flex-row items-center justify-between space-y-0 pb-2"
        >
          <CardTitle class="text-sm font-medium">{{
            getPlatformLabel(platform)
          }}</CardTitle>
          <CardDescription class="text-sm text-muted-foreground">
            {{ t('stream.metrics.every5Min') }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 gap-2 text-sm">
            <div class="text-muted-foreground">{{ t('stream.metrics.totalReqCount') }}</div>
            <div class="font-semibold text-right">
              {{ metric.totalRequests }}
            </div>
            <div class="text-muted-foreground">{{ t('stream.metrics.totalFailedCount') }}</div>
            <div class="font-semibold text-right">{{ metric.totalErrors }}</div>
            <div class="text-muted-foreground">{{ t('stream.metrics.totalFailedRate') }}</div>
            <div
              class="font-semibold text-right"
              :class="
                metric.totalPercent > 50
                  ? 'text-rose-600'
                  : metric.totalPercent >= 10
                  ? 'text-amber-600'
                  : ''
              "
            >
              {{ metric.totalPercent }}%
            </div>
            <div class="text-muted-foreground">{{ t('stream.metrics.mainReqCount') }}</div>
            <div class="font-semibold text-right">
              {{ metric.mainRequests }}
            </div>
            <div class="text-muted-foreground">{{ t('stream.metrics.mainFailedCount') }}</div>
            <div class="font-semibold text-right">{{ metric.mainErrors }}</div>
            <div class="text-muted-foreground">{{ t('stream.metrics.mainFailedRate') }}</div>
            <div
              class="font-semibold text-right"
              :class="
                metric.mainPercent > 50
                  ? 'text-rose-600'
                  : metric.mainPercent >= 10
                  ? 'text-amber-600'
                  : ''
              "
            >
              {{ metric.mainPercent }}%
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import { SquareActivity } from "lucide-vue-next";
import { onMounted, onUnmounted, ref } from "vue";
import type { MetricsData } from "@/types/overview";
import { createSSEConnection } from "@/lib/sse";
import { useDict } from "@/utils/useDict";
import CardDescription from "@/components/ui/card/CardDescription.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const metricsData = ref<Map<string, MetricsData>>(new Map());
const { getLabel: getPlatformLabel } = useDict("live_platform");

let metricsSseClient: any = null;

onMounted(async () => {
  metricsSseClient = createSSEConnection({
    channel: "metric",
    onMessage: (msg) => {
      if (msg.event === "metric") {
        metricsData.value = new Map(Object.entries(msg.data.data));
      }
    },
    onError: (error) => {
      console.error("Metric SSE error:", error);
    },
  });
});

onUnmounted(() => {
  if (metricsSseClient) {
    metricsSseClient.disconnect();
  }
});
</script>
