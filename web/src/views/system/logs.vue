<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <div class="flex items-center space-x-2">
                <Button variant="destructive" @click="openConfirmDeleteAllModal">
                    <Trash2 class="w-4 h-4 mr-2" />
                    {{ t('system.logs.clear.button') }}
                </Button>
                <Button @click="showTerminalLogModal = true">
                    <Terminal class="w-4 h-4 mr-2" />
                    {{ t('system.logs.terminal.button') }}
                </Button>
            </div>
            <div class="flex items-center space-x-2">
                <SortLogsDropDownMenu v-model:sort="sort" @update:sort="handleSortChange" />
                <FilterLogsDropDownMenu v-model:filter="filter" @update:filter="handleFilterChange" />
            </div>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('system.logs.fields.type') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.status') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.content') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.time') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="logs.length > 0">
                        <TableRow v-for="log in logs" :key="log.id">
                            <TableCell class="text-center">
                                <Badge :variant="getTypeVariant(log.type)">{{
                                    getTypeLabel(log.type)
                                }}</Badge>
                            </TableCell>
                            <TableCell class="text-center">
                                <Badge :variant="getStatusVariant(log.status)">{{ getStatusLabel(log.status) }}</Badge>
                            </TableCell>
                            <TableCell class="text-center">
                                <TooltipProvider>
                                    <Tooltip>
                                        <TooltipTrigger as-child>
                                            <div class="truncate max-w-xs mx-auto">{{ log.content }}</div>
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>{{ log.content }}</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </TooltipProvider>
                            </TableCell>
                            <TableCell class="text-center">{{
                                log.createdAt
                            }}</TableCell>
                            <TableCell class="text-center">
                                <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                    @click="openConfirmModal(log.id)">
                                    <Trash2 class="w-4 h-4" />
                                </Button>
                            </TableCell>
                        </TableRow>
                    </template>
                    <template v-else>
                        <TableRow>
                            <TableCell :colspan="5" class="h-24 text-center">
                                {{ t('common.noData') }}
                            </TableCell>
                        </TableRow>
                    </template>
                </TableBody>
            </Table>
        </div>

        <Pagination v-if="total > 0" v-slot="{ page: currentPage }" :total="total" :items-per-page="pageSize"
            :sibling-count="1" show-edges :page="pageNum" @update:page="handlePageChange">
            <PaginationContent v-slot="{ items }">
                <PaginationPrevious />
                <template v-for="(item, index) in items">
                    <PaginationItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                        <Button class="w-10 h-10 p-0" :variant="item.value === currentPage
                            ? 'secondary'
                            : 'outline'
                            " :disabled="item.value === currentPage">
                            {{ item.value }}
                        </Button>
                    </PaginationItem>
                    <PaginationEllipsis v-else :key="item.type" :index="index" />
                </template>
                <PaginationNext />
            </PaginationContent>
        </Pagination>
    </div>
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => (showConfirmModal = open)"
        :onConfirm="handleDeleteLog" :title="t('common.operation.deleteConfirm')"
        :description="t('system.logs.delete.desc')" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => (showConfirmDeleteAllModal = open)"
        :onConfirm="handleDeleteAllLogs" :title="t('system.logs.clear.title')"
        :description="t('system.logs.clear.desc')" />
    <TerminalLogModal :open="showTerminalLogModal" :on-open-change="(open: any) => (showTerminalLogModal = open)" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { listLogs, deleteLogs, deleteAllLogs } from "@/api/system/sys_logs";
import type { SysLogs } from "@/types/system";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import TerminalLogModal from "@/components/modal/system/TerminalLogModal.vue";
import SortLogsDropDownMenu from "@/components/dropdownmenu/system/SortLogsDropDownMenu.vue";
import FilterLogsDropDownMenu from "@/components/dropdownmenu/system/FilterLogsDropDownMenu.vue";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationNext,
    PaginationPrevious,
    PaginationItem,
} from "@/components/ui/pagination";
import {
    TooltipProvider,
    Tooltip,
    TooltipTrigger,
    TooltipContent,
} from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { Trash2, Terminal } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";

const { t } = useI18n();
const showConfirmModal = ref(false);
const showConfirmDeleteAllModal = ref(false);
const showTerminalLogModal = ref(false);
const logs = ref<SysLogs[]>([]);
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);
const logToDelete = ref<number | null>(null);
const sort = ref("id:desc");
const filter = ref({
    type: "",
    status: "",
});

const getLogs = async () => {
    try {
        const params = {
            pageNum: pageNum.value,
            pageSize: pageSize.value,
            sort: sort.value,
            ...filter.value,
        };
        const response = await listLogs(params);
        logs.value = response.data.rows || [];
        total.value = response.data.total || 0;
    } catch (error) {
        console.error("Failed to fetch logs:", error);
    }
};

onMounted(async () => {
    getLogs();
});

const getTypeLabel = (type: number) => {
    switch (type) {
        case 1:
            return t('common.fields.user');
        case 2:
            return t('common.fields.live');
        case 3:
            return t('common.fields.push');
        default:
            return t('common.fields.unknown');
    }
};

const getTypeVariant = (type: number) => {
    switch (type) {
        case 1:
            return "secondary";
        case 2:
            return "outline";
        default:
            return "default";
    }
};

const getStatusLabel = (status: number) => {
    switch (status) {
        case 0:
            return t('common.fields.error');
        case 1:
            return t('common.fields.success');
        default:
            return t('common.fields.unknown');
    }
};

const getStatusVariant = (status: number) => {
    switch (status) {
        case 0:
            return "destructive";
        case 1:
            return "default";
        default:
            return "secondary";
    }
};

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getLogs();
};

const handleSortChange = (newSort: string) => {
    sort.value = newSort;
    getLogs();
};

const handleFilterChange = (newFilter: any) => {
    filter.value = newFilter;
    getLogs();
};

const openConfirmModal = (id: number) => {
    logToDelete.value = id;
    showConfirmModal.value = true;
};

const openConfirmDeleteAllModal = () => {
    showConfirmDeleteAllModal.value = true;
};

async function handleDeleteLog() {
    if (!logToDelete.value) {
        return;
    }
    try {
        const res: any = await deleteLogs(logToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        getLogs();
        toast.success(t('common.toast.deleteSuccess'));
    } catch (error) {
        console.error("Failed to delete log:", error);
    } finally {
        showConfirmModal.value = false;
        logToDelete.value = null;
    }
}

async function handleDeleteAllLogs() {
    try {
        const res: any = await deleteAllLogs({});
        if (res.code !== 0) {
            toast.error(res.msg || t('system.logs.toast.clearFailed'));
            return;
        }
        getLogs();
        toast.success(t('system.logs.toast.clearSuccess'));
    } catch (error) {
        console.error("Failed to delete all logs:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}
</script>
