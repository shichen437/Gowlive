<template>
  <div class="space-y-4">
    <div class="border rounded-md">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="text-center">主播名称</TableHead>
            <TableHead class="text-center">直播开始时间</TableHead>
            <TableHead class="text-center">直播结束时间</TableHead>
            <TableHead class="text-center">直播时长</TableHead>
            <TableHead class="text-center">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="!list || list.length === 0">
            <TableCell :colspan="5" class="text-center h-24">暂无数据</TableCell>
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
    :onConfirm="performDelete" title="确认删除" description="你确定要删除这条历史记录吗？此操作无法撤销。" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { historyList, deleteHistory } from '@/api/stream/live_history';
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

const list = ref<LiveHistory[]>([]);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);


const getList = async () => {
  try {
    const res = await historyList({ pageNum: pageNum.value, pageSize: pageSize.value });
    list.value = res.data.rows;
    total.value = res.data.total;
  } catch (error) {
    console.error(error);
    toast.error("获取列表失败");
  }
};

onMounted(() => {
  getList();
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
    const res: any = await deleteHistory(itemToDelete.value);
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

</script>