<template>
    <div class="p-2 space-y-4">
        <Breadcrumb>
            <BreadcrumbList>
                <template v-for="(crumb, index) in breadcrumbs" :key="crumb.path">
                    <BreadcrumbItem>
                        <BreadcrumbLink v-if="!crumb.isCurrent" @click="navigateToPath(crumb.path)"
                            class="cursor-pointer">
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
                        <TableHead class="w-10"></TableHead>
                        <TableHead>文件名</TableHead>
                        <TableHead>文件大小</TableHead>
                        <TableHead>上次修改时间</TableHead>
                        <TableHead class="text-center">操作</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-for="file in files" :key="file.filename">
                        <TableCell @click="handleItemClick(file)" class="space-y-4 p-4"
                            :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">
                            <component :is="getFileIcon(file)" class="h-4 w-4" />
                        </TableCell>
                        <TableCell class="text-center truncate max-w-xs mx-auto" @click="handleItemClick(file)"
                            :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{ file.filename }}
                        </TableCell>
                        <TableCell class="text-center md:w-[100px]" @click="handleItemClick(file)"
                            :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{ file.isFolder ? "-" :
                                formatBytes(file.size)
                            }}</TableCell>
                        <TableCell class="text-center md:w-[220px]" @click="handleItemClick(file)"
                            :class="{ 'cursor-pointer hover:bg-muted/50': file.isFolder }">{{
                                formatDate(file.lastModified) }}
                        </TableCell>
                        <TableCell class="text-center md:w-[110px]">
                            <span v-if="file.isFolder || file.filename.endsWith('.db')">
                                -
                            </span>
                            <Button variant="ghost" size="icon" v-if="canPlay(file)" @click="handlePlay(file)">
                                <CirclePlay class="w-4 h-4" />
                            </Button>
                            <Button variant="ghost" size="icon" v-if="canPlay(file)"
                                @click="openFileCheckConfirmModal(file)">
                                <FileCheck class="w-4 h-4" />
                            </Button>
                            <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                @click="openConfirmModal(file)" v-if="!file.isFolder && !file.filename.endsWith('.db')">
                                <Trash2 class="w-4 h-4" />
                            </Button>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </div>
    </div>
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="handleDeleteFile" title="确认删除" description="你确定要删除该文件吗？此操作无法撤销。" />
    <ConfirmModal :open="showFileCheckConfirmModal" :onOpenChange="(open: any) => showFileCheckConfirmModal = open"
        :onConfirm="performFileCheck" title="确认检查" description="你确定要为该文件创建检测任务吗？" />
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { useRoute, useRouter, onBeforeRouteLeave } from "vue-router";
import { listFiles, deleteFile } from "@/api/media/file_manage";
import { postTask } from "@/api/media/file_check";
import type { FileInfo } from "@/types/media";
import { canPlay, isVideo, isAudio } from "@/types/media";
import { setLastFilePath, getLastFilePath } from "@/store/cache";
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
import { Folder, File, Trash2, CirclePlay, Database, Film, CassetteTape, FileCheck } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import { formatBytes, formatDate } from "@/utils/convert";

const route = useRoute();
const router = useRouter();

const files = ref<FileInfo[]>([]);
const currentPath = ref(".");
const showConfirmModal = ref(false);
const fileToDelete = ref<FileInfo | null>(null);
const showFileCheckConfirmModal = ref(false);
const fileToCheck = ref<FileInfo | null>(null);

function getFileIcon(file: FileInfo) {
    if (file.isFolder) {
        return Folder;
    } else if (isVideo(file)) {
        return Film;
    } else if (isAudio(file)) {
        return CassetteTape;
    } else if (file.filename.endsWith('.db')) {
        return Database;
    } else {
        return File;
    }
}

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
    let path = route.query.path as string;
    if (!path) {
        const lastPath = getLastFilePath();
        if (lastPath) {
            path = lastPath;
            router.replace({ query: { path: lastPath } });
        } else {
            path = ".";
        }
    }
    fetchFiles(path);
});

onBeforeRouteLeave(() => {
    setLastFilePath(currentPath.value);
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

const handlePlay = (file: FileInfo) => {
    const routeData = router.resolve({
        name: 'MediaPlay',
        query: {
            path: currentPath.value,
            filename: file.filename
        }
    });
    window.open(routeData.href, '_blank');
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

const openFileCheckConfirmModal = (file: FileInfo) => {
    fileToCheck.value = file;
    showFileCheckConfirmModal.value = true;
};

async function performFileCheck() {
    if (!fileToCheck.value) {
        return;
    }
    try {
        const res: any = await postTask({
            path: currentPath.value,
            filename: fileToCheck.value.filename,
        });
        if (res.code !== 0) {
            toast.error(res.msg || "创建任务失败");
            return;
        }
        toast.success("创建任务成功");
    } catch (error) {
        console.error("Failed to create file check task:", error);
    } finally {
        showFileCheckConfirmModal.value = false;
        fileToCheck.value = null;
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

</script>
