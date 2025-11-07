<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <Button variant="destructive" @click="openConfirmDeleteAllModal">
                <Trash2 class="w-4 h-4 mr-2" />
                清空日志
            </Button>
            <div class="flex items-center space-x-2">
                <SortLogsDropDownMenu v-model:sort="sort" @update:sort="handleSortChange" />
                <FilterLogsDropDownMenu v-model:filter="filter" @update:filter="handleFilterChange" />
            </div>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">类型</TableHead>
                        <TableHead class="text-center">状态</TableHead>
                        <TableHead class="text-center">内容</TableHead>
                        <TableHead class="text-center">时间</TableHead>
                        <TableHead class="text-center">操作</TableHead>
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
                                暂无数据
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
        :onConfirm="handleDeleteLog" title="确认删除" description="你确定要删除这条日志吗？此操作无法撤销。" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => (showConfirmDeleteAllModal = open)"
        :onConfirm="handleDeleteAllLogs" title="确认清空" description="你确定要清空所有日志吗？此操作无法撤销。" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { listLogs, deleteLogs, deleteAllLogs } from "@/api/system/sys_logs";
import type { SysLogs } from "@/types/system";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
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
import { Trash2 } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";

const showConfirmModal = ref(false);
const showConfirmDeleteAllModal = ref(false);
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
            return "用户";
        case 2:
            return "直播";
        case 3:
            return "推送";
        default:
            return "未知";
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
            return "错误";
        case 1:
            return "成功";
        default:
            return "未知";
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
            toast.error(res.msg || "删除失败");
            return;
        }
        getLogs();
        toast.success("删除成功");
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
            toast.error(res.msg || "清空失败");
            return;
        }
        getLogs();
        toast.success("清空成功");
    } catch (error) {
        console.error("Failed to delete all logs:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}
</script>
