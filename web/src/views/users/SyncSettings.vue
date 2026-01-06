<template>
    <div class="space-y-4">
        <h2 class="text-2xl font-bold tracking-tight">{{ t('user.syncSettings.title') }}</h2>
        <p class="text-muted-foreground">{{ t('user.syncSettings.desc') }}</p>
        <Card>
            <CardContent class="space-y-4">
                <div class="grid grid-cols-2 gap-4">
                    <div class="flex items-center justify-between">
                        <Label for="data-sync-enable" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.syncSettings.enable.title') }}</span>
                        </Label>
                        <Switch id="data-sync-enable" :checked="dataSyncEnable" v-model="dataSyncEnable"
                            @update:checked="updateSetting('sk_data_sync_enable', $event)"
                            :disabled="openlist.status === 0" />
                    </div>
                    <div class="flex items-center justify-between">
                        <Label for="data-sync-failed-retry" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.syncSettings.failedRetry') }}</span>
                        </Label>
                        <Switch id="data-sync-failed-retry" :checked="dataSyncFailedRetry" v-model="dataSyncFailedRetry"
                            @update:checked="updateSetting('data-sync-failed-retry', $event)"
                            :disabled="openlist.status === 0" />
                    </div>
                </div>
                <div class="grid grid-cols-2 gap-4">
                    <div class="flex items-center justify-between">
                        <Label for="data-sync-auto-delete" class="flex flex-col space-y-1">
                            <span class="text-md">{{ t('user.syncSettings.deleteOnFinish') }}</span>
                        </Label>
                        <Switch id="data-sync-auto-delete" :checked="dataSyncAutoDelete" v-model="dataSyncAutoDelete"
                            @update:checked="updateSetting('sk_data_sync_auto_delete', $event)"
                            :disabled="openlist.status === 0" />
                    </div>
                    <div class="flex items-center justify-between">

                    </div>
                </div>
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <div class="grid gap-2">
                    <div class="flex items-center justify-between">
                        <Label class="flex flex-col space-y-1">
                            <span class="text-md">Openlist</span>
                        </Label>
                        <span :class="openlist.color">{{ openlist.text }}</span>
                    </div>
                </div>
            </CardContent>
        </Card>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
    Card,
    CardContent,
} from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { getSettings, updateSettings } from '@/api/system/settings';
import { openlistStatus } from '@/api/third/openlist';
import { toast } from "vue-sonner";

const { t } = useI18n();

const dataSyncEnable = ref(false);
const dataSyncFailedRetry = ref(false);
const dataSyncAutoDelete = ref(false);
const openlist = ref({ status: 0, text: "", color: "" })

async function fetchSetting(key: string): Promise<Record<string, number>> {
    const res: any = await getSettings({ key });
    if (res.code === 0) {
        return res.data.data;
    }
    return {};
};

async function fetchOpenlistStatus() {
    const statusMap: Record<number, { text: string, color: string }> = {
        0: { text: t('user.syncSettings.noConfig'), color: "text-gray-500" },
        200: { text: t('user.syncSettings.connnected'), color: "text-green-500" },
        400: { text: t('user.syncSettings.connectFailed'), color: "text-red-500" },
    };

    try {
        const res: any = await openlistStatus();
        const status = res.data?.status;
        const statusInfo = statusMap[status];

        if (res.code === 0 && statusInfo) {
            openlist.value = { status, ...statusInfo };
        } else {
            // If res.code is not 0, or status is not in statusMap, it's a connection failure
            openlist.value = { status: 400, text: t('user.syncSettings.connectionFailed'), color: "text-red-500" };
        }
    } catch (error) {
        console.error('Error fetching Openlist status:', error);
        openlist.value = { status: 400, text: t('user.syncSettings.connectionFailed'), color: "text-red-500" };
    }
}
let loadingSettings = true
onMounted(async () => {
    try {
        const result: Record<string, number> = await fetchSetting('sk_data_sync_enable,sk_data_sync_failed_retry,sk_data_sync_auto_delete');
        dataSyncEnable.value = result['sk_data_sync_enable'] == 1;
        dataSyncFailedRetry.value = result['sk_data_sync_failed_retry'] == 1;
        dataSyncAutoDelete.value = result['sk_data_sync_auto_delete'] == 1;

        bindToggleSetting(dataSyncEnable, 'sk_data_sync_enable');
        bindToggleSetting(dataSyncFailedRetry, 'sk_data_sync_failed_retry');
        bindToggleSetting(dataSyncAutoDelete, 'sk_data_sync_auto_delete');

        const defaultStatusInfo = { text: t('user.syncSettings.notEnabled'), color: "text-gray-500" };
        openlist.value = { status: 0, ...defaultStatusInfo };
        fetchOpenlistStatus();
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
        toast.error(t('user.syncSettings.toast.updateErr'));
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
