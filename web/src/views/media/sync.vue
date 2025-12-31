<template>
    <div class="space-y-4">
        <div class="flex items-center justify-between">
            <Button variant="destructive" @click="openConfirmModal(null)">
                <Trash2 class="w-4 h-4 mr-2" />
                {{ t('media.sync.clearAll') }}
            </Button>
        </div>
        <div class="border rounded-lg">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('media.common.fields.filename') }}</TableHead>
                        <TableHead class="text-center">{{ t('media.sync.fields.syncPath') }}</TableHead>
                        <TableHead class="w-40 text-center">{{ t('media.common.fields.duration') }}</TableHead>
                        <TableHead class="w-40 text-center">{{ t('common.fields.status') }}</TableHead>
                        <TableHead class="w-48 text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-if="!tasks || tasks.length === 0">
                        <TableCell :colspan="6" class="text-center h-24">{{ t('common.noData') }}</TableCell>
                    </TableRow>
                    <TableRow v-for="task in tasks" :key="task.id">
                        <TableCell class="truncate max-w-xs text-center">{{ task.filename }}</TableCell>
                        <TableCell class="truncate max-w-xs text-center">{{ task.syncPath }}</TableCell>
                        <TableCell class="text-center">{{ task.duration }}</TableCell>
                        <TableCell class="text-center">
                            <Badge :variant="getStatusVariant(task.status)">
                                {{ getStatusText(task.status) }}
                            </Badge>
                        </TableCell>
                        <TableCell class="text-center">
                            <template v-if="task.status === 3">
                                <Button variant="ghost" size="sm" @click="handleResync(task.id)">
                                    <RefreshCw class="w-4 h-4 mr-2" />
                                </Button>
                            </template>
                            <Button variant="ghost" size="sm" class="text-destructive hover:text-destructive"
                                @click="openConfirmModal(task.id)">
                                <Trash2 class="w-4 h-4 mr-2" />
                            </Button>
                        </TableCell>
                    </TableRow>

                </TableBody>
            </Table>
        </div>
    </div>
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="handleConfirm" :title="confirmTitle" :description="confirmDescription" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { listTasks, deleteTask, resyncTask, deleteAll } from '@/api/media/file_sync';
import type { FileSyncTask } from '@/types/media';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Trash2, RefreshCw, Archive } from 'lucide-vue-next';
import ConfirmModal from '@/components/modal/ConfirmModal.vue';
import { toast } from 'vue-sonner';

const { t } = useI18n();

const tasks = ref<FileSyncTask[]>([]);
const showConfirmModal = ref(false);
const confirmTitle = ref('');
const confirmDescription = ref('');
const currentTaskId = ref<number | null>(null);

const getTasks = async () => {
    try {
        const res: any = await listTasks({});
        if (res.code === 0) {
            tasks.value = res.data.rows;
        } else {
            toast.error(res.msg || t('media.sync.toast.loadFailed'));
        }
    } catch (error) {
        toast.error(t('media.sync.toast.loadFailed'));
    }
};

onMounted(() => {
    getTasks();
});

const getStatusText = (status: number) => {
    switch (status) {
        case 1:
            return t('media.sync.fields.status.syncing');
        case 2:
            return t('media.sync.fields.status.success');
        case 3:
            return t('media.sync.fields.status.failed');
        case 4:
            return t('media.sync.fields.status.notFound');
        default:
            return t('media.sync.fields.status.pending');
    }
};

const getStatusVariant = (status: number): 'default' | 'secondary' | 'destructive' | 'outline' => {
    switch (status) {
        case 1:
            return 'default';
        case 2:
            return 'secondary';
        case 3:
            return 'destructive';
        case 4:
            return 'destructive';
        default:
            return 'outline';
    }
};

const handleResync = async (id: number) => {
    try {
        const res: any = await resyncTask(id);
        if (res.code === 0) {
            toast.success(t('media.sync.toast.resyncSuccess'));
            getTasks();
        } else {
            toast.error(res.msg || t('media.sync.toast.resyncFailed'));
        }
    } catch (error) {
        toast.error(t('media.sync.toast.resyncFailed'));
    }
};


const openConfirmModal = (id: number | null) => {
    currentTaskId.value = id;
    if (id === null) {
        confirmTitle.value = t('media.sync.confirm.clearAllTitle');
        confirmDescription.value = t('media.sync.confirm.clearAllDesc');
    } else {
        confirmTitle.value = t('common.operation.deleteConfirm');
        confirmDescription.value = t('media.sync.confirm.deleteDesc');
    }
    showConfirmModal.value = true;
};

const handleConfirm = () => {
    if (currentTaskId.value === null) {
        handleDeleteAll();
    } else {
        handleDelete(currentTaskId.value);
    }
};


const handleDelete = async (id: number) => {
    try {
        const res: any = await deleteTask(id);
        if (res.code === 0) {
            toast.success(t('common.toast.deleteSuccess'));
            getTasks();
        } else {
            toast.error(res.msg || t('common.toast.deleteFailed'));
        }
    } catch (error) {
        toast.error(t('common.toast.deleteFailed'));
    } finally {
        showConfirmModal.value = false;
    }
};

const handleDeleteAll = async () => {
    try {
        const res: any = await deleteAll({ status: 3 });
        if (res.code === 0) {
            toast.success(t('common.toast.deleteSuccess'));
            getTasks();
        } else {
            toast.error(res.msg || t('common.toast.deleteFailed'));
        }
    } catch (error) {
        toast.error(t('common.toast.deleteFailed'));
    } finally {
        showConfirmModal.value = false;
    }
};
</script>
