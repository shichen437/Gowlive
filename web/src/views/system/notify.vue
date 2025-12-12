<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <div>
                <Button variant="secondary" @click="openConfirmMarkAllModal" class="mr-2">
                    <CheckCheck class="w-4 h-4 mr-2" />
                    {{ t('system.notify.allRead.button') }}
                </Button>
                <Button variant="destructive" @click="openConfirmDeleteAllModal">
                    <Trash2 class="w-4 h-4 mr-2" />
                    {{ t('system.notify.allDelete.button') }}
                </Button>
            </div>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('system.notify.fields.level') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.title') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.content') }}</TableHead>
                        <TableHead class="text-center">{{ t('system.notify.fields.time') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="notifys.length > 0">
                        <TableRow v-for="notify in notifys" :key="notify.id"
                            :class="notify.status === 0 ? 'font-bold' : ''">
                            <TableCell class="text-center">
                                <Badge :variant="getLevelVariant(notify.level)">{{
                                    notify.level === "info"
                                        ? t('common.fields.info')
                                        : notify.level === "warning"
                                            ? t('common.fields.warning')
                                            : notify.level
                                }}</Badge>
                            </TableCell>
                            <TableCell class="text-center">{{
                                notify.title
                            }}</TableCell>
                            <TableCell class="text-center">
                                <TooltipProvider>
                                    <Tooltip>
                                        <TooltipTrigger as-child>
                                            <div class="truncate max-w-xs mx-auto">
                                                {{ notify.content }}
                                            </div>
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>{{ notify.content }}</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </TooltipProvider>
                            </TableCell>
                            <TableCell class="text-center">{{
                                notify.createdAt
                            }}</TableCell>
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
    <ConfirmModal :open="showConfirmDeleteModal" :onOpenChange="(open: any) => (showConfirmDeleteModal = open)"
        :onConfirm="handleDeleteNotify" :title="t('common.operation.deleteConfirm')"
        :description="t('system.notify.delete.desc')" />
    <ConfirmModal :open="showConfirmDeleteAllModal" :onOpenChange="(open: any) => (showConfirmDeleteAllModal = open)"
        :onConfirm="handleDeleteAllNotifys" :title="t('system.notify.allDelete.title')"
        :description="t('system.notify.allDelete.desc')" />
    <ConfirmModal :open="showConfirmMarkAllModal" :onOpenChange="(open: any) => (showConfirmMarkAllModal = open)"
        :onConfirm="handleMarkAllRead" :title="t('system.notify.allRead.title')"
        :description="t('system.notify.allRead.desc')" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import {
    listNotify,
    markNotify,
    markAllNotify,
    deleteNotify,
    deleteAllNotify,
} from "@/api/system/sys_notify";
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
} from "@/components/ui/pagination";
import {
    TooltipProvider,
    Tooltip,
    TooltipTrigger,
    TooltipContent,
} from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { Trash2, BadgeCheck, CheckCheck } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";

const { t } = useI18n();
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
            return "destructive";
        default:
            return "outline";
    }
};

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getNotifys();
};

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
            toast.error(res.msg || t('system.notify.toast.markFailed'));
            return;
        }
        getNotifys();
        toast.success(t('system.notify.toast.markSuccess'));
    } catch (error) {
        console.error("Failed to mark notify:", error);
    }
}

async function handleMarkAllRead() {
    try {
        const res: any = await markAllNotify();
        if (res.code !== 0) {
            toast.error(res.msg || t('system.notify.toast.markAllFailed'));
            return;
        }
        getNotifys();
        toast.success(t('system.notify.toast.markAllSuccess'));
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
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        getNotifys();
        toast.success(t('common.toast.deleteSuccess'));
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
            toast.error(res.msg || t('system.notify.toast.deleteAllFailed'));
            return;
        }
        getNotifys();
        toast.success(t('system.notify.toast.deleteAllSuccess'));
    } catch (error) {
        console.error("Failed to delete all notifys:", error);
    } finally {
        showConfirmDeleteAllModal.value = false;
    }
}
</script>
