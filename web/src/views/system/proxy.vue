<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center mb-4">
            <Button @click="handleAdd">
                <Plus class="w-4 h-4 mr-2" />
                {{ t('system.proxy.add.title') }}
            </Button>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('common.fields.platform') }}</TableHead>
                        <TableHead class="text-center">{{ t('system.proxy.fields.proxy') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.remark') }}</TableHead>
                        <TableHead class="text-center w-[120px]">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-if="!proxies || proxies.length === 0">
                        <TableCell :colspan="4" class="text-center h-24">{{ t('common.noData') }}</TableCell>
                    </TableRow>
                    <TableRow v-for="proxy in proxies" :key="proxy.id">
                        <TableCell class="text-center">{{ getPlatformLabel(proxy.platform) }}</TableCell>
                        <TableCell class="max-w-xs truncate text-center">
                            <div class="truncate">{{ proxy.proxy }}</div>
                        </TableCell>
                        <TableCell class="max-w-24 text-center truncate">{{ proxy.remark }}</TableCell>
                        <TableCell class="text-center space-x-2">
                            <Button variant="ghost" size="icon" @click="handleEdit(proxy)">
                                <Pencil class="w-4 h-4" />
                            </Button>
                            <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                @click="openConfirmModal(proxy.id)">
                                <Trash2 class="w-4 h-4" />
                            </Button>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </div>

        <ProxyModal v-if="isModalOpen" :proxy="selectedProxy" @close="isModalOpen = false" @success="handleSuccess" />
        <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
            :onConfirm="handleDelete" :title="t('common.operation.deleteConfirm')"
            :description="t('system.proxy.deleteDesc')" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { Button } from '@/components/ui/button';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import ProxyModal from '@/components/modal/system/ProxyModal.vue';
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import { allProxies, deleteProxy } from '@/api/system/sys_proxy';
import type { SysProxy } from '@/types/system';
import { toast } from 'vue-sonner';
import { Plus, Pencil, Trash2 } from 'lucide-vue-next';
import { useDict } from '@/utils/useDict';


const { t } = useI18n();
const { getLabel: getPlatformLabel } = useDict("live_platform");

const proxies = ref<SysProxy[]>([]);
const proxyToDelete = ref<number | null>(null);
const isModalOpen = ref(false);
const selectedProxy = ref<SysProxy | null>(null);
const showConfirmModal = ref(false);

const fetchProxies = async () => {
    try {
        const res = await allProxies({});
        proxies.value = res.data.rows;
    } catch (error) {
        console.error('Failed to fetch proxies:', error);
        toast.error(t('system.proxy.toast.listErr'));
    }
};

onMounted(fetchProxies);

const handleAdd = () => {
    selectedProxy.value = null;
    isModalOpen.value = true;
};

const handleEdit = (proxy: SysProxy) => {
    selectedProxy.value = { ...proxy };
    isModalOpen.value = true;
};

const openConfirmModal = (id: number) => {
    proxyToDelete.value = id;
    showConfirmModal.value = true;
};

async function handleDelete() {
    if (!proxyToDelete.value) {
        return;
    }
    try {
        const res: any = await deleteProxy(proxyToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        toast.success(t('common.toast.deleteSuccess'));
        fetchProxies();
    } catch (error) {
        console.error('Failed to delete proxy:', error);
    } finally {
        showConfirmModal.value = false;
        proxyToDelete.value = null;
    }
};

const handleSuccess = () => {
    isModalOpen.value = false;
    fetchProxies();
};
</script>
