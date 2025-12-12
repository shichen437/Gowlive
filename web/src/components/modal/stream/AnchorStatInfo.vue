<template>
    <Dialog :open="isOpen" @update:open="isOpen = $event">
        <DialogContent class="sm:max-w-[750px]">
            <DialogHeader>
                <DialogTitle>历史统计</DialogTitle>
            </DialogHeader>
            <div v-if="statInfo" class="grid gap-4 py-4">
                <div class="grid grid-cols-3 gap-2">
                    <div class="p-4 border rounded-md">
                        <h3 class="text-sm font-medium text-muted-foreground">周点赞</h3>
                        <p class="text-2xl font-bold">{{ statInfo.weekLikeNumIncr }}</p>
                    </div>
                    <div class="p-4 border rounded-md">
                        <h3 class="text-sm font-medium text-muted-foreground">周涨粉</h3>
                        <p class="text-2xl font-bold">{{ statInfo.weekFollowersIncr }}</p>
                    </div>
                    <div class="p-4 border rounded-md">
                        <h3 class="text-sm font-medium text-muted-foreground">月涨粉</h3>
                        <p class="text-2xl font-bold">{{ statInfo.monthFollowersIncr }}</p>
                    </div>
                </div>
                <div v-if="chartData.length > 0">
                    <h3 class="text-s mb-2">粉丝趋势</h3>
                    <div class="border rounded-md p-4">
                        <LineChart :data="chartData" :categories="[t('stream.anchor.fields.followers'), t('stream.anchor.fields.likes')]" index="date" :show-legend="true"
                            :margin="{ left: 60, right: 20, top: 20, bottom: 50 }" :curve-type="CurveType.Linear" />
                    </div>
                </div>
            </div>
            <DialogFooter>
                <Button @click="closeModal">{{ t('common.operation.cancel') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { AnchorStatInfo, AnchorStatData } from '@/types/stream';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { LineChart } from '@/components/ui/chart-line';
import { CurveType } from '@unovis/ts';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface ChartDataItem {
    date: string;
    [key: string]: number | string;
}

const isOpen = ref(false);
const statInfo = ref<AnchorStatInfo | null>(null);

const chartData = computed<ChartDataItem[]>(() => {
    if (!statInfo.value?.historyData || statInfo.value.historyData.length === 0) {
        return [];
    }

    const dataList = [...statInfo.value.historyData];
    let displayData: AnchorStatData[] = [];
    displayData = dataList;

    return displayData.map((item, _) => ({
        date: item.recordDate,
        [t('stream.anchor.fields.followers')]: item.followers,
        [t('stream.anchor.fields.likes')]: item.likeCount
    }));
});

const openModal = (data: AnchorStatInfo) => {
    if (data) {
        statInfo.value = data;
    } else {
        statInfo.value = {
            weekFollowersIncr: 0,
            weekLikeNumIncr: 0,
            monthFollowersIncr: 0,
            historyData: []
        }
    }
    isOpen.value = true;
};

const closeModal = () => {
    isOpen.value = false;
    statInfo.value = null;
};

defineExpose({
    openModal,
});
</script>
