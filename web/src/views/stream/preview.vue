<template>
    <div class="flex h-screen">
        <div class="w-3/4 flex flex-col p-4">
            <div class="flex justify-between items-center mb-4">
                <h1 class="text-lg font-semibold truncate max-w-2xl">
                    {{ roomName ? `${anchor} - ${roomName}` : anchor }}
                </h1>
            </div>
            <div class="grow flex items-center justify-center rounded-lg bg-black/5 overflow-hidden">
                <VideoPlayer v-if="videoUrl" ref="playerRef" :url="videoUrl" :format="videoFormat" :isLive="true"
                    :headers="headers" />
                <div v-else class="text-muted-foreground">{{ t('media.play.toast.loadVideoErr') }}</div>
            </div>
        </div>
        <div class="w-1/4 border-l p-2 flex flex-col">
            <div class="border rounded-lg grow">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{{ t('stream.preview.list') }}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <TableRow v-for="preview in previewListItems" :key="preview.id" @click="playPreview(preview.id)"
                            :class="{ 'text-blue-600 dark:text-blue-300': preview.id === currentPlayingId, 'cursor-pointer': preview.id !== currentPlayingId }">
                            <TableCell>
                                <div class="truncate">
                                    <Badge variant="outline">{{ getPlatformLabel(preview.platform) }}</Badge> {{
                                        preview.anchor }}
                                </div>
                            </TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { toast } from "vue-sonner";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import VideoPlayer from '@/components/player/VideoPlayer.vue';
import { previewRoom, previewList } from "@/api/stream/live_manage";
import { useDict } from "@/utils/useDict";
const { t } = useI18n();
const route = useRoute();
const { getLabel: getPlatformLabel } = useDict("live_platform");

const playerRef = ref<InstanceType<typeof VideoPlayer> | null>(null);
const anchor = ref('');
const roomName = ref('');
const videoUrl = ref('');
const headers = ref<Record<string, string>>({});
const previewListItems = ref<any[]>([]);
const currentPlayingId = ref<number | null>(null);

const videoFormat = computed<'flv' | 'm3u8'>(() => {
    if (videoUrl.value.includes('.flv')) {
        return 'flv';
    }
    return 'm3u8';
});

const playPreview = async (id: number) => {
    if (id === currentPlayingId.value) return;
    try {
        const res: any = await previewRoom(id);
        if (res.code === 0 && res.data && res.data.previewInfo) {
            const { anchor: anchorVal, roomName: roomNameVal, url, headers: headersVal } = res.data.previewInfo;
            anchor.value = anchorVal;
            roomName.value = roomNameVal;
            videoUrl.value = url;
            headers.value = headersVal;
            currentPlayingId.value = id;
        } else {
            videoUrl.value = '';
            toast.error(res.msg || t('stream.preview.toast.loadErr'));
        }
    } catch (error) {
        videoUrl.value = '';
        console.error("Failed to fetch preview data:", error);
        toast.error(t('stream.preview.toast.loadErr'));
    }
};

const fetchPreviewList = async () => {
    try {
        const res: any = await previewList();
        if (res.code === 0 && res.data && res.data.previewList) {
            previewListItems.value = res.data.previewList;
        }
    } catch (error) {
        console.error("Failed to fetch preview list:", error);
        toast.error(t('stream.preview.toast.listLoadErr'));
    }
}

onMounted(() => {
    const id = route.query.id as string;
    if (!id) {
        toast.error(t('stream.preview.toast.noId'));
        return;
    }

    playPreview(Number(id));
    fetchPreviewList();
});
</script>
