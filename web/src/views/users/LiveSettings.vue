<template>
    <div class="space-y-6">
        <Card>
            <CardHeader>
                <CardTitle>直播设置</CardTitle>
                <CardDescription>在这里调整您的直播相关设置。</CardDescription>
            </CardHeader>
            <CardContent>
                <div class="grid gap-4">
                    <div class="flex items-center justify-between">
                        <Label for="live-end-notify" class="flex flex-col space-y-1">
                            <span class="text-md">下播通知</span>
                        </Label>
                        <Switch id="live-end-notify" :checked="liveEndNotify" v-model="liveEndNotify"
                            @update:checked="updateSetting('sk_live_end_notify', $event)" />
                    </div>
                    <div class="flex items-center justify-between">
                        <Label for="filename-template" class="flex flex-col space-y-1">
                            <span class="text-md">文件名称模板</span>
                        </Label>
                        <Select id="filename-template" v-model="filenameTemplate"
                            @update:model-value="updateSetting('sk_filename_template', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue placeholder="选择文件名称模板" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectItem :value=0>[2025-01-01 00:00:00][主播名称][房间名称]</SelectItem>
                                    <SelectItem :value=1>[主播名称][房间名称][2025-01-01 00:00:00]</SelectItem>
                                    <SelectItem :value=2>2025-01-01 00:00:00_主播名称_房间名称</SelectItem>
                                    <SelectItem :value=3>主播名称_房间名称_2025-01-01 00:00:00</SelectItem>
                                    <SelectItem :value=4>[2025-01-01 00:00:00][主播名称]</SelectItem>
                                    <SelectItem :value=5>[主播名称][2025-01-01 00:00:00]</SelectItem>
                                    <SelectItem :value=6>2025-01-01 00:00:00_主播名称</SelectItem>
                                    <SelectItem :value=7>主播名称_2025-01-01 00:00:00</SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                    <div class="flex items-center justify-between">
                        <Label for="archive-strategy" class="flex flex-col space-y-1">
                            <span class="text-md">归档策略</span>
                        </Label>
                        <Select id="archive-strategy" v-model="archiveStrategy"
                            @update:model-value="updateSetting('sk_archive_strategy', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue placeholder="选择归档策略" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectItem :value=0>月</SelectItem>
                                    <SelectItem :value=1>天</SelectItem>
                                    <SelectItem :value=2>月+天</SelectItem>
                                    <SelectItem :value=3>不归档</SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                </div>
            </CardContent>
        </Card>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { getSettings, updateSettings } from '@/api/system/settings';
import { toast } from "vue-sonner";

const liveEndNotify = ref(false);
const filenameTemplate = ref<number>();
const archiveStrategy = ref<number>();

async function fetchSetting(key: string): Promise<Record<string, number>> {
    const res: any = await getSettings({ key });
    if (res.code === 0) {
        return res.data.data;
    }
    return {};
};

onMounted(async () => {
    try {
        const result: Record<string, number> = await fetchSetting('sk_live_end_notify,sk_filename_template,sk_archive_strategy');
        liveEndNotify.value = result['sk_live_end_notify'] == 1;
        filenameTemplate.value = result['sk_filename_template'] || 0;
        archiveStrategy.value = result['sk_archive_strategy'] || 0;
    } catch (error) {
        console.error('Error fetching settings:', error);
    }
});

const updateSetting = async (key: string, value: any) => {
    try {
        await updateSettings({ key, value });
    } catch (error) {
        console.error(`Error updating setting ${key}:`, error);
        toast.error("更新设置失败");
    }
};

watch(liveEndNotify, async (val) => {
    try {
        await updateSettings({ key: 'sk_live_end_notify', value: val ? 1 : 0 });
    } catch (e) {
        console.error(e);
        liveEndNotify.value = !val;
    }
});
</script>
