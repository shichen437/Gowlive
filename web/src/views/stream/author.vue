<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <Button @click="openAddAnchorModal">
                <Plus class="w-4 h-4 mr-2" />
                添加主播
            </Button>
            <div class="flex items-center space-x-2">
                <FilterAuthorDropDownMenu v-model:filter="filter" @update:filter="handleFilterChange" />
            </div>
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">平台</TableHead>
                        <TableHead class="text-center">主播名称</TableHead>
                        <TableHead class="text-center">签名</TableHead>
                        <TableHead class="text-center">关注数</TableHead>
                        <TableHead class="text-center">粉丝数</TableHead>
                        <TableHead class="text-center">点赞数</TableHead>
                        <TableHead class="text-center">视频数</TableHead>
                        <TableHead class="text-center">操作</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="anchors.length > 0">
                        <TableRow v-for="anchor in anchors" :key="anchor.id">
                            <TableCell class="text-center">
                                <Badge variant="outline">{{ getPlatformLabel(anchor.platform) }}</Badge>
                            </TableCell>
                            <TableCell class="text-center">{{ anchor.anchorName }}</TableCell>
                            <TableCell class="text-center">
                                <TooltipProvider>
                                    <Tooltip>
                                        <TooltipTrigger as-child>
                                            <div class="truncate max-w-xs mx-auto">{{ anchor.signature }}</div>
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>{{ anchor.signature }}</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </TooltipProvider>
                            </TableCell>
                            <TableCell class="text-center">{{ formatBigNumber(anchor.followingCount) }}</TableCell>
                            <TableCell class="text-center">{{ formatBigNumber(anchor.followerCount) }}</TableCell>
                            <TableCell class="text-center">{{ formatBigNumber(anchor.likeCount) }}</TableCell>
                            <TableCell class="text-center">{{ anchor.videoCount }}</TableCell>
                            <TableCell class="text-center space-x-2">
                                <Button variant="ghost" size="icon" @click="openStatModal(anchor.id)">
                                    <BarChart class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                    @click="openConfirmModal(anchor.id)">
                                    <Trash2 class="w-4 h-4" />
                                </Button>
                            </TableCell>
                        </TableRow>
                    </template>
                    <template v-else>
                        <TableRow>
                            <TableCell :colspan="9" class="h-24 text-center">
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
    <AuthorModal ref="authorModal" @refresh="getAnchors" />
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="handleDeleteAnchor" title="确认删除" description="你确定要删除该主播吗？此操作无法撤销。" />
    <AnchorStatInfoModal ref="statModal" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { anchorList, deleteAnchor, getAnchorStatInfo } from "@/api/stream/anchor_info";
import type { AnchorInfo } from "@/types/stream";
import AuthorModal from "@/components/modal/stream/AuthorModal.vue";
import AnchorStatInfoModal from "@/components/modal/stream/AnchorStatInfo.vue";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import FilterAuthorDropDownMenu from "@/components/dropdownmenu/stream/FilterAuthorDropDownMenu.vue";
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
import {
    TooltipProvider,
    Tooltip,
    TooltipTrigger,
    TooltipContent,
} from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { BarChart, Plus, Trash2 } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";
import { useDict } from "@/utils/useDict";
import { formatBigNumber } from "@/utils/convert";

const showConfirmModal = ref(false);
const anchors = ref<AnchorInfo[]>([]);
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);
const authorModal = ref<InstanceType<typeof AuthorModal> | null>(null);
const statModal = ref<InstanceType<typeof AnchorStatInfoModal> | null>(null);
const anchorToDelete = ref<number | null>(null);
const filter = ref({
    platform: '',
    nickname: '',
});

const { getLabel: getPlatformLabel } = useDict("live_platform");

const getAnchors = async () => {
    try {
        const params = {
            pageNum: pageNum.value,
            pageSize: pageSize.value,
            ...filter.value,
        };
        const response = await anchorList(params);
        anchors.value = response.data.rows || [];
        total.value = response.data.total || 0;
    } catch (error) {
        console.error("Failed to fetch anchors:", error);
    }
};

onMounted(async () => {
    getAnchors();
});

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getAnchors();
}

const handleFilterChange = (newFilter: any) => {
    filter.value = newFilter;
    getAnchors();
};

const openAddAnchorModal = () => {
    authorModal.value?.openModal();
};

const openStatModal = async (id: number) => {
    try {
        const response: any = await getAnchorStatInfo(id);
        if (response.code === 0) {
            statModal.value?.openModal(response.data.data);
        } else {
            toast.error(response.msg || "获取统计信息失败");
        }
    } catch (error) {
        console.error("Failed to fetch anchor stats:", error);
    }
};

const openConfirmModal = (id: number) => {
    anchorToDelete.value = id;
    showConfirmModal.value = true;
};

async function handleDeleteAnchor() {
    if (!anchorToDelete.value) {
        return;
    }
    try {
        const res: any = await deleteAnchor(anchorToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || "删除失败");
            return;
        }
        getAnchors();
        toast.success("删除成功");
    } catch (error) {
        console.error("Failed to delete anchor:", error);
    } finally {
        showConfirmModal.value = false;
        anchorToDelete.value = null;
    }
}
</script>
