<template>
    <div class="space-y-4">
        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">文件名</TableHead>
                        <TableHead class="text-center">进度</TableHead>
                        <TableHead class="text-center">文件状态</TableHead>
                        <TableHead class="text-center">耗时</TableHead>
                        <TableHead class="text-center">结束时间</TableHead>
                        <TableHead class="text-center">操作</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-if="!list || list.length === 0">
                        <TableCell :colspan="6" class="text-center h-24">暂无数据</TableCell>
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
                        <TableCell class="text-center">{{ item.duration }}秒</TableCell>
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
        :onConfirm="performDelete" title="确认删除" description="你确定要删除这个任务吗？此操作无法撤销。" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { listTasks, deleteTask } from '@/api/media/file_check';
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


const list = ref<FileCheckTask[]>([]);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);
const executing = ref(false);
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
        toast.error("获取列表失败");
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

async function performDelete() {
    if (!itemToDelete.value) return;
    try {
        const res: any = await deleteTask(itemToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || '删除失败');
            return;
        }
        toast.success("删除成功")
        getList();
    } catch (error) {
        console.error(error);
    } finally {
        showConfirmModal.value = false;
        itemToDelete.value = null;
    }
};

const handlePageChange = (page: number) => {
    pageNum.value = page;
    getList();
};

function getFileStatusDisplay(status: number): ({ text: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; }) {
    switch (status) {
        case 0: return { text: "待检测", variant: "default" };
        case 1: return { text: "正常", variant: "outline" };
        case 2: return { text: "疑似损坏", variant: "destructive" };
        case 3: return { text: "文件不存在", variant: "default" };
        default: return { text: "未知", variant: "outline" };
    }
};

function getProgressDisplay(progress: number): ({ text: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; }) {
    switch (progress) {
        case 0: return { text: "等待开始", variant: "default" };
        case 1: return { text: "快检中...", variant: "secondary" };
        case 2: return { text: "快检异常", variant: "destructive" };
        case 3: return { text: "全量检测中...", variant: "secondary" };
        case 4: return { text: "全量检测异常", variant: "destructive" };
        case 5: return { text: "检测完成", variant: "outline" };
        default: return { text: "未知", variant: "secondary" };
    }
};

</script>
