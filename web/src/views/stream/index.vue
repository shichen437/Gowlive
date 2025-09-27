<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <Button @click="openAddRoomModal">
        <Plus class="w-4 h-4 mr-2" />
        添加房间
      </Button>
      <div class="flex items-center space-x-2">
        <SortRoomDropDownMenu v-model:sort="sort" @update:sort="handleSortChange" />
        <FilterRoomDropDownMenu v-model:filter="filter" @update:filter="handleFilterChange" />
      </div>
    </div>

    <div class="border rounded-md">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="text-center">主播名称</TableHead>
            <TableHead class="text-center">房间名称</TableHead>
            <TableHead class="text-center">平台</TableHead>
            <TableHead class="text-center">监控状态</TableHead>
            <TableHead class="text-center">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="rooms.length > 0">
            <TableRow v-for="room in rooms" :key="room.id">
              <TableCell class="text-center">{{ room.anchor }}</TableCell>
              <TableCell class="text-center">{{ room.roomName }}</TableCell>
              <TableCell class="text-center">
                <Badge variant="outline">{{ getPlatformLabel(room.platform) }}</Badge>
              </TableCell>
              <TableCell class="text-center" :class="getStatusColor(room.status, room.isRecording)">{{
                getStatusText(room.status, room.isRecording) }}</TableCell>
              <TableCell class="text-center space-x-2">
                <Button variant="ghost" size="icon" @click="openEditRoomModal(room.liveId)">
                  <Pencil class="w-4 h-4" />
                </Button>
                <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                  @click="openConfirmModal(room.liveId)">
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
  <RoomModal ref="roomModal" @refresh="getRooms" />
  <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
    :onConfirm="handleDeleteRoom" title="确认删除" description="你确定要删除这个房间吗？此操作无法撤销。" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { roomList, deleteRoom } from "@/api/stream/live_manage";
import type { RoomInfo } from "@/types/stream";
import RoomModal from "@/components/modal/stream/RoomModal.vue";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import SortRoomDropDownMenu from "@/components/dropdownmenu/stream/SortRoomDropDownMenu.vue";
import FilterRoomDropDownMenu from "@/components/dropdownmenu/stream/FilterRoomDropDownMenu.vue";
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
import { Plus, Pencil, Trash2 } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";
import { useDict } from "@/utils/useDict";

const showConfirmModal = ref(false);
const rooms = ref<RoomInfo[]>([]);
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);
const roomModal = ref<InstanceType<typeof RoomModal> | null>(null);
const roomToDelete = ref<number | null>(null);
const sort = ref('');
const filter = ref({
  anchor: '',
  roomName: '',
  platform: '',
});

const { getLabel: getPlatformLabel } = useDict("live_platform");

const getRooms = async () => {
  try {
    const params = {
      pageNum: pageNum.value,
      pageSize: pageSize.value,
      sort: sort.value,
      ...filter.value,
    };
    const response = await roomList(params);
    rooms.value = response.data.rows || [];
    total.value = response.data.total || 0;
  } catch (error) {
    console.error("Failed to fetch rooms:", error);
  }
};

onMounted(async () => {
  getRooms();
});

const getStatusColor = (status: number, isRecording: boolean) => {
  if (isRecording) {
    return "text-red-500";
  }
  switch (status) {
    case 1:
      return "text-cyan-600";
    case 2:
      return "text-green-600";
    default:
      return "";
  }
};

const getStatusText = (status: number, isRecording: boolean) => {
  if (isRecording) {
    return "录制中";
  }
  switch (status) {
    case 1:
      return "实时监控";
    case 2:
      return "定时监控";
    default:
      return "未监控";
  }
};

const handlePageChange = (newPage: number) => {
  pageNum.value = newPage;
  getRooms();
}

const handleSortChange = (newSort: string) => {
  sort.value = newSort;
  getRooms();
};

const handleFilterChange = (newFilter: any) => {
  filter.value = newFilter;
  getRooms();
};

const openAddRoomModal = () => {
  roomModal.value?.openModal();
};

const openEditRoomModal = (id: number) => {
  roomModal.value?.openModal(id);
};

const openConfirmModal = (id: number) => {
  roomToDelete.value = id;
  showConfirmModal.value = true;
};

async function handleDeleteRoom() {
  if (!roomToDelete.value) {
    return;
  }
  try {
    const res: any = await deleteRoom(roomToDelete.value);
    if (res.code !== 0) {
      toast.error(res.msg || "删除失败");
      return;
    }
    getRooms();
    toast.success("删除成功");
  } catch (error) {
    console.error("Failed to delete room:", error);
  } finally {
    showConfirmModal.value = false;
    roomToDelete.value = null;
  }
}
</script>
