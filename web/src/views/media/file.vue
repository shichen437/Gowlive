<template>
  <div class="p-2 space-y-4">
    <Breadcrumb>
      <BreadcrumbList>
        <template v-for="(crumb, index) in breadcrumbs" :key="crumb.path">
          <BreadcrumbItem>
            <BreadcrumbLink v-if="!crumb.isCurrent" @click="navigateToPath(crumb.path)" class="cursor-pointer">
              {{ crumb.label }}
            </BreadcrumbLink>
            <BreadcrumbPage v-else>
              {{ crumb.label }}
            </BreadcrumbPage>
          </BreadcrumbItem>
          <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
        </template>
      </BreadcrumbList>
    </Breadcrumb>

    <div class="border rounded-lg">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[40px]"></TableHead>
            <TableHead>文件名</TableHead>
            <TableHead>文件大小</TableHead>
            <TableHead>上次修改时间</TableHead>
            <TableHead class="text-center">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="file in files" :key="file.filename">
            <TableCell @click="handleItemClick(file)"
              :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">
              <component :is="file.isFolder ? Folder : File" class="h-5 w-5" />
            </TableCell>
            <TableCell class="text-center" @click="handleItemClick(file)"
              :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{ file.filename }}</TableCell>
            <TableCell class="text-center" @click="handleItemClick(file)"
              :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{ file.isFolder ? "-" : formatSize(file.size)
              }}</TableCell>
            <TableCell class="text-center" @click="handleItemClick(file)"
              :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{ formatDate(file.lastModified) }}
            </TableCell>
            <TableCell class="text-center">
              <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                @click="openConfirmModal(file)" :disabled="file.isFolder || file.filename.endsWith('.db')">
                <Trash2 class="w-4 h-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
  <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
    :onConfirm="handleDeleteFile" title="确认删除" description="你确定要删除这个文件吗？此操作无法撤销。" />
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listFiles, deleteFile } from "@/api/media/file_manage";
import type { FileInfo } from "@/types/media";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Folder, File, Trash2 } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";

const route = useRoute();
const router = useRouter();

const files = ref<FileInfo[]>([]);
const currentPath = ref(".");
const showConfirmModal = ref(false);
const fileToDelete = ref<FileInfo | null>(null);

const fetchFiles = async (path: string) => {
  try {
    const res: any = await listFiles({ path: path });
    if (res.code !== 0) {
      toast.error(res.msg || "获取文件列表失败")
      return
    }
    files.value = res.data.rows;
    currentPath.value = path;
  } catch (error) {
    console.error("Failed to fetch files:", error);
  }
};

onMounted(() => {
  const path = route.query.path as string || ".";
  fetchFiles(path);
});

watch(
  () => route.query.path,
  (newPath) => {
    fetchFiles((newPath as string) || ".");
  }
);

const handleItemClick = (file: FileInfo) => {
  if (file.isFolder) {
    const newPath = `${currentPath.value}/${file.filename}`.replace(/^\.\//, '');
    router.push({ query: { path: newPath } });
  }
};

const openConfirmModal = (file: FileInfo) => {
  fileToDelete.value = file;
  showConfirmModal.value = true;
};

async function handleDeleteFile() {
  if (!fileToDelete.value) {
    return;
  }
  try {
    const res: any = await deleteFile({
      path: currentPath.value,
      filename: fileToDelete.value.filename,
    });
    if (res.code !== 0) {
      toast.error(res.msg || "删除失败");
      return;
    }
    fetchFiles(currentPath.value);
    toast.success("删除成功");
  } catch (error) {
    console.error("Failed to delete file:", error);
  } finally {
    showConfirmModal.value = false;
    fileToDelete.value = null;
  }
}

const breadcrumbs = computed(() => {
  const pathParts = currentPath.value.split("/").filter(p => p);
  if (pathParts.length === 0 || pathParts[0] !== '.') {
    pathParts.unshift(".");
  }

  return pathParts.map((part, index) => {
    const path = pathParts.slice(0, index + 1).join("/");
    return {
      label: part === '.' ? '系统目录' : part,
      path: path,
      isCurrent: index === pathParts.length - 1,
    };
  });
});

const navigateToPath = (path: string) => {
  router.push({ query: { path } });
};

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(2)} MB`;
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
};

const formatDate = (timestamp: number) => {
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};
</script>