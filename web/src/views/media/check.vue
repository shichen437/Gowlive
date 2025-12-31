<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <Button variant="destructive" @click="openConfirmDeleteAllModal">
                <Trash2 class="w-4 h-4 mr-2" />
                {{ t('media.check.clear.button') }}
            </Button>
        </div>
        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('media.common.fields.filename') }}</TableHead>
                        <TableHead class="text-center">{{ t('media.check.fields.progress.title') }}</TableHead>
                        <TableHead class="text-center">{{ t('media.check.fields.status.title') }}</TableHead>
                        <TableHead class="text-center">{{ t('media.common.fields.duration') }}</TableHead>
                        <TableHead class="text-center">{{ t('media.check.fields.finishTime') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-if="!list || list.length === 0">
                        <TableCell :colspan="6" class="text-center h-24">{{ t('common.noData') }}</TableCell>
                    </TableRow>
                    <TableRow v-for="item in list" :key="item.id">
                        <TableCell class="text-center">{{ item.filename }}</TableCell>
                        <TableCell class="text-center">
                            <Badge :variant="getProgressDisplay(item.progress).variant">
                                {{ getProgressDisplay(item.progress).text }}
                            </Badge>
                        </TableCell>
                        <TableCell class="text-center">
                            <Badge :variant="getFileStatusDisplay(item.fileStatus).variant">
                                {{ getFileStatusDisplay(item.fileStatus).text }}
                            </Badge>
                        </TableCell>
                        <TableCell class="text-center">{{ item.duration }}{{ t('common.fields.second') }}</TableCell>
                        <TableCell class="text-center">{{ item.updatedAt }}</TableCell>
                        <TableCell class="text-center space-x-2">
                            <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                @click="openConfirmModal(item.id)">
                                <Trash2 class="w-4 h-4" />
                            </Button>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </div>
        <Pagination v-if="total > 0" v-slot="{ page: currentPage }" :total="total" :items-per-page="pageSize"
            :sibling-count="1" show-edges :page="pageNum" @update:page="handlePageChange">
            <PaginationContent v-slot="{ items }">
                <PaginationPrevious />

                <template v-for="(item, index) in items">
                    <PaginationItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                        <Button class="w-10 h-10 p-0" :variant="item.value === currentPage ? 'secondary' : 'outline'"
                            :disabled="item.value === currentPage">
                            {{ item.value }}
                        </Button>
                    </PaginationItem>
                    <PaginationEllipsis v-else :key="item.type" :index="index" />
                </template>

                <PaginationNext />
            </PaginationContent>
        </Pagination>
    </div>

    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="performDelete" :title="t('common.operation.deleteConfirm')"
        :description="t('media.check.deleteConfirmDesc')" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => (showConfirmDeleteAllModal = open)"
        :onConfirm="handleDeleteAllTasks" :title="t('media.check.clear.title')"
        :description="t('media.check.clear.desc')" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { listTasks, deleteTask, deleteAllTask } from '@/api/media/file_check';
import type { FileCheckTask } from '@/types/media';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationItem,
    PaginationNext,
    PaginationPrevious,
} from '@/components/ui/pagination';
import { Button } from '@/components/ui/button';
import { toast } from "vue-sonner";
import { Trash2 } from 'lucide-vue-next';
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import { Badge } from '@/components/ui/badge';

const { t } = useI18n();
const list = ref<FileCheckTask[]>([]);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);
const executing = ref(false);
const showConfirmDeleteAllModal = ref(false);
let timer: any = null;

const getList = async () => {
    if (timer) {
        clearTimeout(timer);
        timer = null;
    }
    try {
        const res = await listTasks({ pageNum: pageNum.value, pageSize: pageSize.value });
        list.value = res.data.rows;
        total.value = res.data.total;
        executing.value = res.data.executing;

        if (executing.value) {
            timer = setTimeout(() => {
                getList();
            }, 2000);
        }
    } catch (error) {
        console.error(error);
        toast.error(t('media.check.toast.listErr'));
    }
};

onMounted(() => {
    getList();
});

onUnmounted(() => {
    if (timer) {
        clearTimeout(timer);
    }
});

const showConfirmModal = ref(false);
const itemToDelete = ref<number | null>(null);

const openConfirmModal = (id: number) => {
    itemToDelete.value = id;
    showConfirmModal.value = true;
};

const openConfirmDeleteAllModal = () => {
    showConfirmDeleteAllModal.value = true;
};

async function performDelete() {
    if (!itemToDelete.value) return;
    try {
        const res: any = await deleteTask(itemToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        toast.success(t('common.toast.deleteSuccess'))
        getList();
    } catch (error) {
        console.error(error);
    } finally {
        showConfirmModal.value = false;
        itemToDelete.value = null;
    }
};

async function handleDeleteAllTasks() {
    try {
        const res: any = await deleteAllTask();
        if (res.code !== 0) {
            toast.error(res.msg || t('media.check.toast.clearErr'));
            return;
        }
        getList();
        toast.success(t('media.check.clear.success'));
    } catch (error) {
        console.error("Failed to delete all tasks:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}

const handlePageChange = (page: number) => {
    pageNum.value = page;
    getList();
};

function getFileStatusDisplay(status: number): ({ text: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; }) {
    switch (status) {
        case 0: return { text: t('media.check.fields.status.pending'), variant: "default" };
        case 1: return { text: t('media.check.fields.status.normal'), variant: "outline" };
        case 2: return { text: t('media.check.fields.status.bad'), variant: "destructive" };
        case 3: return { text: t('media.check.fields.status.notExist'), variant: "default" };
        default: return { text: t('media.check.fields.status.unknown'), variant: "outline" };
    }
};

function getProgressDisplay(progress: number): ({ text: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; }) {
    switch (progress) {
        case 0: return { text: t('media.check.fields.progress.pending'), variant: "default" };
        case 1: return { text: t('media.check.fields.progress.qcing'), variant: "secondary" };
        case 2: return { text: t('media.check.fields.progress.qcerr'), variant: "destructive" };
        case 3: return { text: t('media.check.fields.progress.acing'), variant: "secondary" };
        case 4: return { text: t('media.check.fields.progress.acerr'), variant: "destructive" };
        case 5: return { text: t('media.check.fields.progress.success'), variant: "outline" };
        default: return { text: t('media.check.fields.status.unknown'), variant: "secondary" };
    }
};

</script>
