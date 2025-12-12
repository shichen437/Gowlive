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
                        <TableHead>{{ t('media.common.fields.filename') }}</TableHead>
                        <TableHead>{{ t('media.file.fields.filesize') }}</TableHead>
                        <TableHead>{{ t('media.file.fields.lastModified') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
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
                            <span v-if="!canPlay(file) && !showButtons">
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
                                @click="openConfirmModal(file)" v-if="showButtons">
                                <Trash2 class="w-4 h-4" />
                            </Button>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </div>
    </div>
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => showConfirmModal = open"
        :onConfirm="handleDeleteFile" :title="t('common.operation.deleteConfirm')" :description="confirmModalDescription" />
    <ConfirmModal :open="showFileCheckConfirmModal" :onOpenChange="(open: any) => showFileCheckConfirmModal = open"
        :onConfirm="performFileCheck" :title="t('media.file.confirmCheck.title')" :description="t('media.file.confirmCheck.desc')" />
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { useRoute, useRouter, onBeforeRouteLeave } from "vue-router";
import { listFiles, deleteFile, getEmptyDir } from "@/api/media/file_manage";
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
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const files = ref<FileInfo[]>([]);
const currentPath = ref(".");
const showConfirmModal = ref(false);
const fileToDelete = ref<FileInfo | null>(null);
const showFileCheckConfirmModal = ref(false);
const fileToCheck = ref<FileInfo | null>(null);
const confirmModalDescription = ref("");

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

const openFileCheckConfirmModal = (file: FileInfo) => {
    fileToCheck.value = file;
    showFileCheckConfirmModal.value = true;
};

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

const openConfirmModal = async (file: FileInfo) => {
    fileToDelete.value = file;
    if (file.isFolder) {
        try {
            const res: any = await getEmptyDir({ path: `${currentPath.value}/${file.filename}`.replace(/^\.\//, '') });
            if (res.code !== 0) {
                toast.error(res.msg || t('media.file.toast.folderStatus'));
                return;
            }
            if (res.data.isEmpty) {
                confirmModalDescription.value = t('media.file.delete.folder');
            } else {
                confirmModalDescription.value = t('media.file.delete.folderNotEmpty');
            }
        } catch (error) {
            console.error("Failed to check if dir is empty:", error);
            toast.error(t('media.file.toast.folderStatus'));
            return;
        }
    } else {
        confirmModalDescription.value = t('media.file.delete.file');
    }
    showConfirmModal.value = true;
};

const fetchFiles = async (path: string) => {
    try {
        const res: any = await listFiles({ path: path });
        if (res.code !== 0) {
            toast.error(res.msg || t('media.file.toast.list'))
            return
        }
        files.value = res.data.rows;
        currentPath.value = path;
    } catch (error) {
        console.error("Failed to fetch files:", error);
    }
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
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        fetchFiles(currentPath.value);
        toast.success(t('common.toast.deleteSuccess'));
    } catch (error) {
        console.error("Failed to delete file:", error);
    } finally {
        showConfirmModal.value = false;
        fileToDelete.value = null;
    }
}

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
            toast.error(res.msg || t('media.file.toast.createcCheck'));
            return;
        }
        toast.success(t('media.file.createcCheckSuccess'));
    } catch (error) {
        console.error("Failed to create file check task:", error);
    } finally {
        showFileCheckConfirmModal.value = false;
        fileToCheck.value = null;
    }
}

const showButtons = computed(() => {
    const pathParts = currentPath.value.split("/").filter(p => p);
    if (pathParts.length === 0 || pathParts[0] !== '.') {
        pathParts.unshift(".");
    }
    return pathParts.length > 2;
});

const breadcrumbs = computed(() => {
    const pathParts = currentPath.value.split("/").filter(p => p);
    if (pathParts.length === 0 || pathParts[0] !== '.') {
        pathParts.unshift(".");
    }

    return pathParts.map((part, index) => {
        const path = pathParts.slice(0, index + 1).join("/");
        return {
            label: part === '.' ? t('media.file.fields.systemDir') : part,
            path: path,
            isCurrent: index === pathParts.length - 1,
        };
    });
});

const navigateToPath = (path: string) => {
    router.push({ query: { path } });
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

</script>
