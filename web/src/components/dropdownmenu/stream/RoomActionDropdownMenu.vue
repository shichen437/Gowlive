<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon">
                <MoreHorizontal class="w-4 h-4" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
            <DropdownMenuItem v-if="!room.isTop" @click="$emit('topRoom', room.liveId)">
                <Pin class="w-4 h-4 mr-2" />
                <span>{{ t('stream.rooms.buttons.top') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem v-else @click="$emit('unTopRoom', room.liveId)">
                <PinOff class="w-4 h-4 mr-2" />
                <span>{{ t('stream.rooms.buttons.unTop') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem v-if="room.isRecording" @click="$emit('streamPreview', room.liveId)">
                <Eye class="w-4 h-4 mr-2" />
                <span>{{ t('stream.rooms.buttons.preview') }}</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="$emit('goToFolder', room)">
                <Folder class="w-4 h-4 mr-2" />
                <span>{{ t('stream.rooms.buttons.openFolder') }}</span>
            </DropdownMenuItem>
            <Separator class="my-1" />
            <DropdownMenuItem @click="$emit('deleteRoom', room.liveId)" variant="destructive">
                <Trash2 class="w-4 h-4 mr-2" />
                <span>{{ t('common.operation.delete') }}</span>
            </DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>

<script setup lang="ts">
import type { RoomInfo } from "@/types/stream";
import { useI18n } from "vue-i18n";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { MoreHorizontal, Pin, PinOff, Folder, Trash2, Eye } from "lucide-vue-next";

const { t } = useI18n();

defineProps<{
    room: RoomInfo;
}>();

defineEmits<{
    (e: "topRoom", id: number): void;
    (e: "unTopRoom", id: number): void;
    (e: "streamPreview", id: number): void;
    (e: "goToFolder", room: RoomInfo): void;
    (e: "deleteRoom", id: number): void;
}>();
</script>
