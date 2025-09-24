<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center mb-4">
      <Button @click="handleAdd">
        <Plus class="w-4 h-4 mr-2" />
        添加Cookie
      </Button>
    </div>

    <div class="border rounded-md">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="text-center">平台</TableHead>
            <TableHead class="text-center">Cookie</TableHead>
            <TableHead class="text-center">备注</TableHead>
            <TableHead class="text-center w-[120px]">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="!cookies || cookies.length === 0">
            <TableCell :colspan="4" class="text-center h-24">暂无数据</TableCell>
          </TableRow>
          <TableRow v-for="cookie in cookies" :key="cookie.id">
            <TableCell class="text-center">{{ getPlatformLabel(cookie.platform) }}</TableCell>
            <TableCell class="max-w-xs truncate text-center">
              <div class="truncate">{{ cookie.cookie }}</div>
            </TableCell>
            <TableCell class="max-w-24 text-center truncate">{{ cookie.remark }}</TableCell>
            <TableCell class="text-center space-x-2">
              <Button variant="ghost" size="icon" @click="handleEdit(cookie)">
                <Pencil class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                @click="openConfirmModal(cookie.id)">
                <Trash2 class="w-4 h-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <CookieModal v-if="isModalOpen" :cookie="selectedCookie" @close="isModalOpen = false" @success="handleSuccess" />
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
      :onConfirm="handleDelete" title="确认删除" description="你确定要删除这个Cookie吗？此操作无法撤销。" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import CookieModal from '@/components/modal/stream/CookieModal.vue';
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import { allCookies, deleteCookie } from '@/api/stream/live_cookie';
import type { LiveCookie } from '@/types/stream';
import { toast } from 'vue-sonner';
import { Plus, Pencil, Trash2 } from 'lucide-vue-next';
import { useDict } from '@/utils/useDict';

const { getLabel: getPlatformLabel } = useDict("live_platform");

const cookies = ref<LiveCookie[]>([]);
const cookieToDelete = ref<number | null>(null);
const isModalOpen = ref(false);
const selectedCookie = ref<LiveCookie | null>(null);
const showConfirmModal = ref(false);

const fetchCookies = async () => {
  try {
    const res = await allCookies({});
    cookies.value = res.data.rows;
  } catch (error) {
    console.error('Failed to fetch cookies:', error);
    toast.error('获取Cookie列表失败');
  }
};

onMounted(fetchCookies);

const handleAdd = () => {
  selectedCookie.value = null;
  isModalOpen.value = true;
};

const handleEdit = (cookie: LiveCookie) => {
  selectedCookie.value = { ...cookie };
  isModalOpen.value = true;
};

const openConfirmModal = (id: number) => {
  cookieToDelete.value = id;
  showConfirmModal.value = true;
};

async function handleDelete() {
  if (!cookieToDelete.value) {
    return;
  }
  try {
    const res: any = await deleteCookie(cookieToDelete.value);
    if (res.code !== 0) {
      toast.error(res.msg || '删除失败');
      return;
    }
    toast.success('删除成功');
    fetchCookies();
  } catch (error) {
    console.error('Failed to delete cookie:', error);
  } finally {
    showConfirmModal.value = false;
    cookieToDelete.value = null;
  }
};

const handleSuccess = () => {
  isModalOpen.value = false;
  fetchCookies();
};
</script>
