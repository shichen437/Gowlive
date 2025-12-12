<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <Button variant="destructive" @click="openConfirmDeleteAllModal">
                <Trash2 class="w-4 h-4 mr-2" />
                {{ t('stream.history.clear.button') }}
            </Button>
        </div>
        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('stream.common.fields.anchorName') }}</TableHead>
                        <TableHead class="text-center">{{ t('stream.history.fields.startTime') }}</TableHead>
                        <TableHead class="text-center">{{ t('stream.history.fields.endTime') }}</TableHead>
                        <TableHead class="text-center">{{ t('stream.history.fields.duration') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-if="!list || list.length === 0">
                        <TableCell :colspan="5" class="text-center h-24">{{ t('common.noData') }}</TableCell>
                    </TableRow>
                    <TableRow v-for="item in list" :key="item.id">
                        <TableCell class="text-center">{{ item.anchor }}</TableCell>
                        <TableCell class="text-center">{{ item.startedAt }}</TableCell>
                        <TableCell class="text-center">{{ item.endedAt }}</TableCell>
                        <TableCell class="text-center">{{ item.duration }}</TableCell>
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
        :onConfirm="performDelete" :title="t('common.operation.deleteConfirm')" :description="t('stream.history.deleteDesc')" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => (showConfirmDeleteAllModal = open)"
        :onConfirm="handleDeleteAllHistory" :title="t('stream.history.clear.title')" :description="t('stream.history.clear.desc')" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { historyList, deleteHistory, deleteAllHistory } from '@/api/stream/live_history';
import type { LiveHistory } from '@/types/stream';
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

const { t } = useI18n();
const list = ref<LiveHistory[]>([]);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);
const showConfirmDeleteAllModal = ref(false);


const getList = async () => {
    try {
        const res = await historyList({ pageNum: pageNum.value, pageSize: pageSize.value });
        list.value = res.data.rows;
        total.value = res.data.total;
    } catch (error) {
        console.error(error);
        toast.error(t('stream.history.toast.listErr'));
    }
};

onMounted(() => {
    getList();
});

const showConfirmModal = ref(false);
const itemToDelete = ref<number | null>(null);

const openConfirmDeleteAllModal = () => {
    showConfirmDeleteAllModal.value = true;
};

const openConfirmModal = (id: number) => {
    itemToDelete.value = id;
    showConfirmModal.value = true;
};

async function performDelete() {
    if (!itemToDelete.value) return;
    try {
        const res: any = await deleteHistory(itemToDelete.value);
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

async function handleDeleteAllHistory() {
    try {
        const res: any = await deleteAllHistory();
        if (res.code !== 0) {
            toast.error(res.msg || t('stream.history.clear.failed'));
            return;
        }
        getList();
        toast.success(t('stream.history.clear.success'));
    } catch (error) {
        console.error("Failed to delete all history:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}

const handlePageChange = (page: number) => {
    pageNum.value = page;
    getList();
};

</script>
