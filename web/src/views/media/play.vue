<template>
    <div class="flex h-screen">
        <div class="w-3/4 flex flex-col p-4">
            <h1 class="text-lg font-semibold mb-4 truncate">
                <span v-if="filename">{{ filename }}</span>
                <span v-else>No file selected</span>
            </h1>
            <div class="grow flex items-center justify-center rounded-lg">
                <VideoPlayer v-if="videoUrl" :url="videoUrl" :format="videoFormat" :isLive="false" />
                <div v-else class="text-muted-foreground">无法加载视频，缺少地址或参数。</div>
            </div>
        </div>

        <div class="w-1/4 border-l p-2 flex flex-col">
            <div class="border rounded-lg grow">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>播放列表</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <template v-for="(file, index) in files" :key="file ? file.filename : index">
                            <TableRow v-if="canPlay(file)" @click="handleItemClick(file)"
                                :class="{ 'text-blue-600 dark:text-blue-300': file.filename === filename, 'cursor-pointer': canPlay(file) && file.filename !== filename }">
                                <TableCell class="w-60 max-w-60 whitespace-nowrap overflow-hidden text-ellipsis">
                                    {{ file.filename }}
                                </TableCell>
                            </TableRow>
                        </template>
                    </TableBody>
                </Table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import VideoPlayer from '@/components/player/VideoPlayer.vue';
import { listFiles } from "@/api/media/file_manage";
import type { FileInfo } from "@/types/media";
import { toast } from "vue-sonner";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useTheme } from "@/composables/useTheme";

useTheme();
const route = useRoute();
const router = useRouter();

const filename = ref('');
const currentPath = ref('');
const files = ref<FileInfo[]>([]);

const videoFormat = computed<'flv' | 'mp4' | 'mp3' | 'mkv' | 'ts'>(() => {
    const ext = (filename.value.split('.').pop() || '').toLowerCase();
    if (ext === 'flv') return 'flv';
    if (ext === 'mp4') return 'mp4';
    if (ext === 'mkv') return 'mkv';
    if (ext === 'ts') return 'ts';
    if (ext === 'mp3') return 'mp3';
    return 'mp4';
});

function buildApiUrl(path: string, filename: string): string {
    const apiBaseUrl = (import.meta as any).env?.VITE_APP_BASE_API || '';
    const fullPath = `${path.replace(/\/+$/, '')}/${filename}`.replace(/^\.\//, '');
    const url = `${apiBaseUrl}/media/file/play?path=${encodeURIComponent(fullPath)}`;
    return url;
}

const videoUrl = computed(() => {
    if (currentPath.value && filename.value) {
        const url = buildApiUrl(currentPath.value, filename.value);
        console.log('Playback URL:', url);
        return url;
    }
    return '';
});

function canPlay(file: FileInfo) {
    return isVideo(file) || isAudio(file);
}

function isVideo(file: FileInfo) {
    return !file.isFolder && (file.filename.endsWith('.mp4') || file.filename.endsWith('.flv') || file.filename.endsWith('.mkv') || file.filename.endsWith('.ts'));
}

function isAudio(file: FileInfo) {
    return !file.isFolder && file.filename.endsWith('.mp3');
}

const fetchFiles = async (path: string) => {
    try {
        const res: any = await listFiles({ path: path });
        if (res.code !== 0) {
            toast.error(res.msg || "获取文件列表失败")
            return
        }
        files.value = (res.data.rows || []).filter(Boolean);
    } catch (error) {
        console.error("Failed to fetch files:", error);
        toast.error("获取文件列表失败");
    }
};

const handleItemClick = (file: FileInfo) => {
    if (canPlay(file)) {
        filename.value = file.filename;
        router.replace({ query: { path: currentPath.value, filename: file.filename } });
    }
};

onMounted(() => {
    const path = (route.query.path as string) || '.';
    const file = (route.query.filename as string) || '';

    currentPath.value = path;
    filename.value = file;
    fetchFiles(path);
});

watch(() => route.query.path, (newPath) => {
    const path = (newPath as string) || ".";
    currentPath.value = path;
    fetchFiles(path);
    if (!route.query.filename) {
        filename.value = '';
    }
});

watch(() => route.query.filename, (newFilename) => {
    filename.value = (newFilename as string) || '';
});
</script>
