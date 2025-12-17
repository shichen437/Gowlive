<template>
    <div class="flex h-screen">
        <div class="w-3/4 flex flex-col p-4">
            <div class="flex justify-between items-center mb-4">
                <h1 class="text-lg font-semibold truncate max-w-2xl">
                    <span v-if="filename">{{ filename }}</span>
                    <span v-else>{{ t('media.play.toast.loadVideoErr') }}</span>
                </h1>
                <div class="space-x-2" v-if="filename">
                    <Button variant="outline" size="sm" @click="openClipModal">
                        <Scissors class="w-4 h-4 mr-2" />
                        {{ t('media.play.clip.button') }}
                    </Button>
                </div>
            </div>
            
            <div class="grow flex items-center justify-center rounded-lg bg-black/5 overflow-hidden">
                <VideoPlayer ref="playerRef" v-if="videoUrl" :url="videoUrl" :format="videoFormat" :isLive="false" />
                <div v-else class="text-muted-foreground">{{ t('media.play.toast.loadVideoErr') }}</div>
            </div>
        </div>

        <div class="w-1/4 border-l p-2 flex flex-col">
            <div class="border rounded-lg grow">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{{ t('media.play.playlist') }}</TableHead>
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

    <Dialog :open="isClipModalOpen" @update:open="isClipModalOpen = $event">
        <DialogContent class="sm:max-w-[500px]">
            <DialogHeader>
                <DialogTitle>{{ t('media.play.clip.button') }}</DialogTitle>
                <DialogDescription>
                    {{ t('media.play.clip.desc') }}
                </DialogDescription>
            </DialogHeader>
            <div class="grid gap-6 py-4">
                <div class="grid grid-cols-2 gap-4">
                    <div class="space-y-2">
                        <Label>{{ t('media.play.clip.fields.startTime') }}</Label>
                        <div class="flex gap-2">
                            <Input v-model="clipStart" placeholder="00:00:00" />
                            <Button variant="secondary" size="icon" @click="setStartToCurrent" :title="t('media.play.clip.fields.setStartTime')">
                                <MapPin class="w-4 h-4" />
                            </Button>
                        </div>
                    </div>
                    <div class="space-y-2">
                        <Label>{{ t('media.play.clip.fields.endTime') }}</Label>
                        <div class="flex gap-2">
                            <Input v-model="clipEnd" placeholder="00:00:00" />
                            <Button variant="secondary" size="icon" @click="setEndToCurrent" :title="t('media.play.clip.fields.setEndTime')">
                                <MapPin class="w-4 h-4" />
                            </Button>
                        </div>
                    </div>
                </div>
                <div class="text-sm text-muted-foreground bg-muted p-2 rounded">
                   {{ t('media.play.clip.tooltip') }}
                </div>
            </div>
            <DialogFooter>
                <Button variant="outline" @click="isClipModalOpen = false">{{ t('common.operation.cancel') }}</Button>
                <Button @click="handleClipSubmit" :disabled="isClipping">
                    {{ isClipping ? t('media.play.clip.clipping') : t('media.play.clip.startClip') }}
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import VideoPlayer from '@/components/player/VideoPlayer.vue';
import { listFiles, clipFile } from "@/api/media/file_manage";
import type { FileInfo } from "@/types/media";
import { canPlay } from "@/types/media";
import { toast } from "vue-sonner";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { 
    Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter 
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from "@/components/ui/button";
import { Scissors, MapPin } from "lucide-vue-next"; 
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const filename = ref('');
const currentPath = ref('');
const files = ref<FileInfo[]>([]);
const playerRef = ref<InstanceType<typeof VideoPlayer> | null>(null);

// 剪辑相关状态
const isClipModalOpen = ref(false);
const clipStart = ref("00:00:00");
const clipEnd = ref("00:00:00");
const isClipping = ref(false);

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
        return url;
    }
    return '';
});

const fetchFiles = async (path: string) => {
    try {
        const res: any = await listFiles({ path: path });
        if (res.code !== 0) {
            toast.error(res.msg || t('media.play.toast.listErr'))
            return
        }
        files.value = (res.data.rows || []).filter(Boolean);
    } catch (error) {
        console.error("Failed to fetch files:", error);
        toast.error(t('media.play.toast.listErr'));
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

// 格式化时间 helper
const formatSeconds = (seconds: number) => {
    const date = new Date(0);
    date.setSeconds(seconds);
    return date.toISOString().substr(11, 8);
};

const openClipModal = () => {
    isClipModalOpen.value = true;
    const duration = playerRef.value?.getDuration() || 0;
    if (duration > 0 && clipEnd.value === "00:00:00") {
        clipEnd.value = formatSeconds(duration);
    }
};

const setStartToCurrent = () => {
    const time = playerRef.value?.getCurrentTime() || 0;
    clipStart.value = formatSeconds(time);
};

const setEndToCurrent = () => {
    const time = playerRef.value?.getCurrentTime() || 0;
    clipEnd.value = formatSeconds(time);
};

const handleClipSubmit = async () => {
    if (!filename.value) return;
    
    if (clipStart.value >= clipEnd.value) {
        toast.error(t('media.play.toast.endBeforeStart'));
        return;
    }

    isClipping.value = true;
    try {
        const res: any = await clipFile({
            path: currentPath.value,
            filename: filename.value,
            startTime: clipStart.value,
            endTime: clipEnd.value
        });

        if (res.code === 0) {
            toast.success(t('media.play.clip.success'));
            isClipModalOpen.value = false;
            // 刷新文件列表以显示新文件
            fetchFiles(currentPath.value);
        } else {
            toast.error(res.msg || t('media.play.toast.clipErr'));
        }
    } catch (error) {
        console.error("Clip failed:", error);
        toast.error(t('media.play.toast.requestFailed'));
    } finally {
        isClipping.value = false;
    }
};
</script>
