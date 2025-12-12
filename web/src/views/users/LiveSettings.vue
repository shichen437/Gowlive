<template>
    <div class="space-y-4">
        <h2 class="text-2xl font-bold tracking-tight">{{ t('user.liveSettings.title') }}</h2>
        <p class="text-muted-foreground">{{ t('user.liveSettings.desc') }}</p>
        <Card>
            <CardContent>
                <div class="grid gap-2">
                    <div class="flex items-center justify-between">
                        <Label for="archive-strategy" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.liveSettings.archive.title') }}</span>
                        </Label>
                        <Select id="archive-strategy" v-model="archiveStrategy"
                            @update:model-value="updateSetting('sk_archive_strategy', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue :placeholder="t('user.liveSettings.archive.placeholder')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectItem :value=0>{{ t('user.liveSettings.archive.values.monthly') }}</SelectItem>
                                    <SelectItem :value=1>{{ t('user.liveSettings.archive.values.daily') }}</SelectItem>
                                    <SelectItem :value=2>{{ t('user.liveSettings.archive.values.monthlyAndDaily') }}</SelectItem>
                                    <SelectItem :value=3>{{ t('user.liveSettings.archive.values.none') }}</SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                    <div class="flex items-center justify-between">
                        <Label for="filename-template" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.liveSettings.filenameTemplate.title') }}</span>
                        </Label>
                        <Select id="filename-template" v-model="filenameTemplate"
                            @update:model-value="updateSetting('sk_filename_template', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue :placeholder="t('user.liveSettings.filenameTemplate.placeholder')" />
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
                </div>
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <div class="grid gap-2">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-1">
                            <Label for="disk-protection" class="flex flex-col space-y-1">
                                <span class="text-md">{{ t('user.liveSettings.diskProtected.title') }}</span>
                            </Label>
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger as-child>
                                        <Info class="w-4 h-4 ml-1" />
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>{{ t('user.liveSettings.diskProtected.tooltip') }}</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                        </div>
                        <Select id="disk-protection" v-model="diskProtection"
                            @update:model-value="updateSetting('sk_disk_protection', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue :placeholder="t('user.liveSettings.diskProtected.placeholder')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectItem :value=0>禁用</SelectItem>
                                    <SelectItem :value=5>5GB</SelectItem>
                                    <SelectItem :value=10>10GB</SelectItem>
                                    <SelectItem :value=20>20GB</SelectItem>
                                    <SelectItem :value=50>50GB</SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-1">
                            <Label for="auto-clean-little-file" class="flex flex-col space-y-1">
                                <span class="text-md">{{ t('user.liveSettings.autoCleanLittleFile.title') }}</span>
                            </Label>
                        </div>
                        <Select id="auto-clean-little-file" v-model="autoCleanLittleFile"
                            @update:model-value="updateSetting('sk_auto_clean_little_file', $event)">
                            <SelectTrigger class="w-[330px]">
                                <SelectValue :placeholder="t('user.liveSettings.autoCleanLittleFile.placeholder')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectItem :value=0>禁用</SelectItem>
                                    <SelectItem :value=20>20MB</SelectItem>
                                    <SelectItem :value=50>50MB</SelectItem>
                                    <SelectItem :value=100>100MB</SelectItem>
                                    <SelectItem :value=200>200MB</SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </div>
                </div>
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <div class="grid gap-2">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-1">
                            <Label for="fixed-resolution" class="flex flex-col space-y-1">
                                <span class="text-md">{{ t('user.liveSettings.fixedResolution.title') }}</span>
                            </Label>
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger as-child>
                                        <Badge variant="secondary">BETA</Badge>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>{{ t('user.liveSettings.fixedResolution.tooltip') }}</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                        </div>
                        <Switch id="fixed-resolution" :checked="fixedResolution" v-model="fixedResolution"
                            @update:checked="updateSetting('sk_fixed_resolution', $event)" />
                    </div>
                </div>
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <div class="grid gap-2">
                    <div class="flex items-center justify-between">
                        <Label for="live-end-notify" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.liveSettings.liveCloseNotice') }}</span>
                        </Label>
                        <Switch id="live-end-notify" :checked="liveEndNotify" v-model="liveEndNotify"
                            @update:checked="updateSetting('sk_live_end_notify', $event)" />
                    </div>
                </div>
            </CardContent>
        </Card>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import {
    Card,
    CardContent,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
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
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger
} from '@/components/ui/tooltip';
import { Info } from 'lucide-vue-next';
import { getSettings, updateSettings } from '@/api/system/settings';
import { toast } from "vue-sonner";

const { t } = useI18n();

const liveEndNotify = ref(false);
const fixedResolution = ref(false);
const filenameTemplate = ref<number>();
const archiveStrategy = ref<number>();
const diskProtection = ref<number>();
const autoCleanLittleFile = ref<number>();

async function fetchSetting(key: string): Promise<Record<string, number>> {
    const res: any = await getSettings({ key });
    if (res.code === 0) {
        return res.data.data;
    }
    return {};
};
let loadingSettings = true
onMounted(async () => {
    try {
        const result: Record<string, number> = await fetchSetting('sk_live_end_notify,sk_filename_template,sk_archive_strategy,sk_disk_protection,sk_auto_clean_little_file,sk_fixed_resolution');
        liveEndNotify.value = result['sk_live_end_notify'] == 1;
        fixedResolution.value = result['sk_fixed_resolution'] == 1;
        filenameTemplate.value = result['sk_filename_template'] || 0;
        archiveStrategy.value = result['sk_archive_strategy'] || 0;
        diskProtection.value = result['sk_disk_protection'] || 0;
        autoCleanLittleFile.value = result['sk_auto_clean_little_file'] || 0;

        bindToggleSetting(liveEndNotify, 'sk_live_end_notify');
        bindToggleSetting(fixedResolution, 'sk_fixed_resolution');
        loadingSettings = false
    } catch (error) {
        console.error('Error fetching settings:', error);
    }
});

const updateSetting = async (key: string, value: any) => {
    try {
        await updateSettings({ key, value });
    } catch (error) {
        console.error(`Error updating setting ${key}:`, error);
        toast.error(t('user.liveSettings.toast.updateErr'));
    }
};

function bindToggleSetting(refVar: Ref<boolean>, key: string) {
    watch(refVar, async (val) => {
        if (loadingSettings) return
        const oldVal = !val;
        try {
            await updateSettings({ key, value: val ? 1 : 0 });
        } catch (e) {
            console.error(e);
            refVar.value = oldVal;
        }
    });
}
</script>
