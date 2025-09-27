<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <div>
                <Button variant="secondary" @click="openConfirmMarkAllModal" class="mr-2">
                    <CheckCheck class="w-4 h-4 mr-2" />
                    全部已读
                </Button>
                <Button variant="destructive" @click="openConfirmDeleteAllModal">
                    <Trash2 class="w-4 h-4 mr-2" />
                    全部删除
                </Button>
            </div>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">通知级别</TableHead>
                        <TableHead class="text-center">标题</TableHead>
                        <TableHead class="text-center">内容</TableHead>
                        <TableHead class="text-center">通知时间</TableHead>
                        <TableHead class="text-center">操作</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="notifys.length > 0">
                        <TableRow v-for="notify in notifys" :key="notify.id"
                            :class="notify.status === 0 ? 'font-bold' : ''">
                            <TableCell class="text-center">
                                <Badge :variant="getLevelVariant(notify.level)">{{ notify.level === 'info' ? '信息' :
                                    notify.level === 'warning' ? '预警' : notify.level }}</Badge>
                            </TableCell>
                            <TableCell class="text-center">{{ notify.title }}</TableCell>
                            <TableCell class="text-center">{{ notify.content }}</TableCell>
                            <TableCell class="text-center">{{ notify.createdAt }}</TableCell>
                            <TableCell class="text-center">
                                <Button v-if="notify.status === 0" variant="ghost" size="icon"
                                    class="text-blue-500 hover:text-blue-600" @click="handleMarkRead(notify.id)">
                                    <BadgeCheck class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                    @click="openConfirmDeleteModal(notify.id)">
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
    <ConfirmModal :open="showConfirmDeleteModal" :onOpenChange="(open: any) => showConfirmDeleteModal = open"
        :onConfirm="handleDeleteNotify" title="确认删除" description="你确定要删除这条通知吗？此操作无法撤销。" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => showConfirmDeleteAllModal = open"
        :onConfirm="handleDeleteAllNotifys" title="确认全部删除" description="你确定要删除所有通知吗？此操作无法撤销。" />
    <ConfirmModal :open="showConfirmMarkAllModal" :onOpenChange="(open: any) => showConfirmMarkAllModal = open"
        :onConfirm="handleMarkAllRead" title="确认全部已读" description="你确定要将所有通知标记为已读吗？" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { listNotify, markNotify, markAllNotify, deleteNotify, deleteAllNotify } from "@/api/system/sys_notify";
import type { SysNotify } from "@/types/system";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
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
} from '@/components/ui/pagination';
import { Button } from "@/components/ui/button";
import { Trash2, BadgeCheck, CheckCheck } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";

const showConfirmDeleteModal = ref(false);
const showConfirmDeleteAllModal = ref(false);
const showConfirmMarkAllModal = ref(false);
const notifys = ref<SysNotify[]>([]);
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);
const notifyToDelete = ref<number | null>(null);

const getNotifys = async () => {
    try {
        const params = {
            pageNum: pageNum.value,
            pageSize: pageSize.value,
        };
        const response = await listNotify(params);
        notifys.value = response.data.rows || [];
        total.value = response.data.total || 0;
    } catch (error) {
        console.error("Failed to fetch notifys:", error);
    }
};

onMounted(async () => {
    getNotifys();
});

const getLevelVariant = (level: string) => {
    switch (level) {
        case "info":
            return "default";
        case "warning":
            return "secondary";
        default:
            return "outline";
    }
};

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getNotifys();
}

const openConfirmDeleteModal = (id: number) => {
    notifyToDelete.value = id;
    showConfirmDeleteModal.value = true;
};

const openConfirmDeleteAllModal = () => {
    showConfirmDeleteAllModal.value = true;
};

const openConfirmMarkAllModal = () => {
    showConfirmMarkAllModal.value = true;
};

async function handleMarkRead(id: number) {
    try {
        const res: any = await markNotify(id);
        if (res.code !== 0) {
            toast.error(res.msg || "标记失败");
            return;
        }
        getNotifys();
        toast.success("标记成功");
    } catch (error) {
        console.error("Failed to mark notify:", error);
    }
}

async function handleMarkAllRead() {
    try {
        const res: any = await markAllNotify();
        if (res.code !== 0) {
            toast.error(res.msg || "全部标记失败");
            return;
        }
        getNotifys();
        toast.success("全部标记成功");
    } catch (error) {
        console.error("Failed to mark all notifys:", error);
    } finally {
        showConfirmMarkAllModal.value = false;
    }
}

async function handleDeleteNotify() {
    if (!notifyToDelete.value) {
        return;
    }
    try {
        const res: any = await deleteNotify(notifyToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || "删除失败");
            return;
        }
        getNotifys();
        toast.success("删除成功");
    } catch (error) {
        console.error("Failed to delete notify:", error);
    } finally {
        showConfirmDeleteModal.value = false;
        notifyToDelete.value = null;
    }
}

async function handleDeleteAllNotifys() {
    try {
        const res: any = await deleteAllNotify();
        if (res.code !== 0) {
            toast.error(res.msg || "全部删除失败");
            return;
        }
        getNotifys();
        toast.success("全部删除成功");
    } catch (error) {
        console.error("Failed to delete all notifys:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}
</script>
