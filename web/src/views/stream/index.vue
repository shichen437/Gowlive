<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <div class="flex items-center space-x-2">
                <Button @click="openAddRoomModal">
                    <Plus class="w-4 h-4 mr-2" />
                    {{ t('stream.rooms.add.title') }}
                </Button>
                <Button @click="openBatchAddRoomModal" variant="secondary">
                    <CopyPlus class="w-4 h-4 mr-2" />
                    {{ t('stream.rooms.batch.button') }}
                </Button>
            </div>
            <div class="flex items-center space-x-2">
                <div class="flex items-center">
                    <Button variant="outline" size="icon"
                        @click="setDisplayMode(displayMode === 'list' ? 'card' : 'list')"
                        :aria-label="displayMode === 'list' ? t('stream.rooms.buttons.cardView') : t('stream.rooms.buttons.listView')">
                        <LayoutGrid v-if="displayMode === 'card'" class="w-4 h-4" />
                        <List v-else class="w-4 h-4" />
                    </Button>
                </div>
                <ExportRoomDropDownMenu @export="handleExport" />
                <SortRoomDropDownMenu v-model:sort="sort" @update:sort="handleSortChange" />
                <FilterRoomDropDownMenu v-model:filter="filter" @update:filter="handleFilterChange" />
            </div>
        </div>

        <div v-if="displayMode === 'list'" class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('stream.common.fields.anchorName') }}</TableHead>
                        <TableHead class="text-center">{{ t('stream.rooms.fields.roomName') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.platform') }}</TableHead>
                        <TableHead class="text-center">{{ t('stream.rooms.fields.monitorStatus.title') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="rooms.length > 0">
                        <TableRow v-for="room in rooms" :key="room.id" :class="{ 'bg-top/80': room.isTop }">
                            <TableCell class="text-center">
                                {{ room.anchor }}
                            </TableCell>
                            <TableCell class="text-center">
                                {{ room.roomName }}
                            </TableCell>
                            <TableCell class="text-center">
                                <Badge variant="outline">
                                    {{ getPlatformLabel(room.platform) }}
                                </Badge>
                            </TableCell>
                            <TableCell class="text-center" :class="getStatusColor(
                                room.status,
                                room.isRecording,
                                room.isLiving,
                            )
                                ">{{
                                    getStatusText(room.status, room.isRecording, room.isLiving)
                                }}</TableCell>
                            <TableCell class="text-center space-x-2">
                                <Button v-if="room.status === 0" variant="ghost" size="icon"
                                    @click="openStartConfirmModal(room.liveId)">
                                    <Play class="w-4 h-4" />
                                </Button>
                                <Button v-else variant="ghost" size="icon"
                                    class="text-destructive hover:text-destructive"
                                    @click="openStopConfirmModal(room.liveId)">
                                    <Square class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" @click="openEditRoomModal(room.liveId)">
                                    <Pencil class="w-4 h-4" />
                                </Button>

                                <RoomActionDropdownMenu :room="room" @top-room="handleTopRoom"
                                    @un-top-room="handleUnTopRoom" @go-to-folder="handleGoToRoomFolder"
                                    @delete-room="openConfirmModal" @stream-preview="handlePreview" />
                            </TableCell>
                        </TableRow>
                    </template>
                    <template v-else>
                        <TableRow>
                            <TableCell :colspan="5" class="h-24 text-center">
                                {{ t('common.noData') }}
                            </TableCell>
                        </TableRow>
                    </template>
                </TableBody>
            </Table>
        </div>
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <template v-if="rooms.length > 0">
                <Card v-for="room in rooms" :key="room.id" :class="{ 'bg-top/80': room.isTop }">
                    <CardHeader class="flex flex-row items-center justify-between space-y-0">
                        <CardTitle class="text-lg truncate" :title="room.anchor">{{ room.anchor }}</CardTitle>
                        <div class="flex items-center space-x-2">
                            <Badge variant="outline">{{ getPlatformLabel(room.platform) }}</Badge>
                        </div>
                    </CardHeader>
                    <CardContent>
                        <CardDescription class="line-clamp-1" :title="room.roomName">
                            {{ room.roomName == "" ? "-" : room.roomName }}
                        </CardDescription>
                    </CardContent>
                    <CardFooter class="flex justify-between items-center pt-0">
                        <Badge variant="secondary" class="text-sm" :class="getStatusColor(
                            room.status,
                            room.isRecording,
                            room.isLiving,
                        )
                            ">{{
                                getStatusText(room.status, room.isRecording, room.isLiving)
                            }}</Badge>
                        <div class="flex space-x-1">
                            <Button v-if="room.status === 0" variant="ghost" size="icon"
                                @click="openStartConfirmModal(room.liveId)">
                                <Play class="w-4 h-4" />
                            </Button>
                            <Button v-else variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                @click="openStopConfirmModal(room.liveId)">
                                <Square class="w-4 h-4" />
                            </Button>
                            <Button variant="ghost" size="icon" @click="openEditRoomModal(room.liveId)">
                                <Pencil class="w-4 h-4" />
                            </Button>

                            <RoomActionDropdownMenu :room="room" @top-room="handleTopRoom"
                                @un-top-room="handleUnTopRoom" @go-to-folder="handleGoToRoomFolder"
                                @delete-room="openConfirmModal" @stream-preview="handlePreview" />
                        </div>
                    </CardFooter>
                </Card>
            </template>
            <template v-else>
                <div class="col-span-full">
                    <div class="border rounded-md h-32 flex justify-center items-center">
                        {{ t('common.noData') }}
                    </div>
                </div>
            </template>
        </div>


        <Pagination v-if="total > 0" v-slot="{ page: currentPage }" :total="total" :items-per-page="pageSize"
            :sibling-count="1" show-edges :page="pageNum" @update:page="handlePageChange">
            <PaginationContent v-slot="{ items }">
                <PaginationPrevious />

                <template v-for="(item, index) in items">
                    <PaginationItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                        <Button class="w-10 h-10 p-0" :variant="item.value === currentPage
                            ? 'secondary'
                            : 'outline'
                            " :disabled="item.value === currentPage">
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
    <BatchAddRoomModal ref="batchAddRoomModal" @refresh="getRooms" />
    <ConfirmModal :open="showConfirmModal" :onOpenChange="(open: any) => (showConfirmModal = open)"
        :onConfirm="handleDeleteRoom" :title="t('common.operation.deleteConfirm')"
        :description="t('stream.rooms.deleteDesc')" />
    <ConfirmModal :open="showStartConfirmModal" :onOpenChange="(open: any) => (showStartConfirmModal = open)"
        :onConfirm="handleStartRoom" :title="t('stream.rooms.startMonitor.title')"
        :description="t('stream.rooms.startMonitor.desc')" />
    <ConfirmModal :open="showStopConfirmModal" :onOpenChange="(open: any) => (showStopConfirmModal = open)"
        :onConfirm="handleStopRoom" :title="t('stream.rooms.stopMonitor.title')"
        :description="t('stream.rooms.stopMonitor.desc')" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import {
    roomList,
    deleteRoom,
    startRoom,
    stopRoom,
    topRoom,
    unTopRoom,
    exportRoomInfo,
} from "@/api/stream/live_manage";
import { getRoomFilePath } from "@/api/media/file_manage";
import { getStreamDisplayMode, setStreamDisplayMode } from "@/store/cache";
import type { RoomInfo } from "@/types/stream";
import RoomModal from "@/components/modal/stream/RoomModal.vue";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import BatchAddRoomModal from "@/components/modal/stream/BatchAddRoomModal.vue";
import SortRoomDropDownMenu from "@/components/dropdownmenu/stream/SortRoomDropDownMenu.vue";
import FilterRoomDropDownMenu from "@/components/dropdownmenu/stream/FilterRoomDropDownMenu.vue";
import ExportRoomDropDownMenu from "@/components/dropdownmenu/stream/ExportRoomDropDownMenu.vue";
import RoomActionDropdownMenu from "@/components/dropdownmenu/stream/RoomActionDropdownMenu.vue";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationNext,
    PaginationPrevious,
    PaginationItem,
} from "@/components/ui/pagination";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Plus, Pencil, Trash2, Play, Square, CopyPlus, Folder, Pin, PinOff, MoreHorizontal, List, LayoutGrid } from "lucide-vue-next";
import { toast } from "vue-sonner";
import { Badge } from "@/components/ui/badge";
import { useDict } from "@/utils/useDict";

const { t } = useI18n();
const router = useRouter();

const showConfirmModal = ref(false);
const rooms = ref<RoomInfo[]>([]);
const pageNum = ref(1);
const pageSize = ref(12);
const total = ref(0);
const roomModal = ref<InstanceType<typeof RoomModal> | null>(null);
const batchAddRoomModal = ref<InstanceType<typeof BatchAddRoomModal> | null>(
    null,
);
const roomToDelete = ref<number | null>(null);
const sort = ref("");
const filter = ref({
    anchor: "",
    roomName: "",
    platform: "",
});

const displayMode = ref(getStreamDisplayMode() || "list");

const setDisplayMode = (mode: "list" | "card") => {
    displayMode.value = mode;
    setStreamDisplayMode(mode);
};

const showStartConfirmModal = ref(false);
const roomToStart = ref<number | null>(null);
const showStopConfirmModal = ref(false);
const roomToStop = ref<number | null>(null);

const { getLabel: getPlatformLabel } = useDict("live_platform");

const getRooms = async () => {
    try {
        const params = {
            pageNum: pageNum.value,
            pageSize: pageSize.value,
            sort: sort.value,
            ...filter.value,
        };
        const response: any = await roomList(params);
        rooms.value = response.data.rows || [];
        total.value = response.data.total || 0;
    } catch (error) {
        console.error("Failed to fetch rooms:", error);
    } finally {
        scheduleRefresh();
    }
};

onMounted(async () => {
    getRooms();
});

const getStatusColor = (status: number, isRecording: boolean, isLiving: boolean) => {
    if (isRecording || isLiving) {
        return "text-red-600";
    }
    switch (status) {
        case 1:
            return "text-cyan-600";
        case 2:
            return "text-green-600";
        case 3:
            return "text-violet-600";
        default:
            return "";
    }
};

const getStatusText = (status: number, isRecording: boolean, isLiving: boolean) => {
    if (isRecording) {
        return t('stream.rooms.fields.monitorStatus.recording');
    }
    if (isLiving) {
        return t('stream.rooms.fields.monitorStatus.living');
    }
    switch (status) {
        case 1:
            return t('stream.rooms.fields.monitorStatus.realtime');
        case 2:
            return t('stream.rooms.fields.monitorStatus.cron');
        case 3:
            return t('stream.rooms.fields.monitorStatus.intelligent');
        default:
            return t('stream.rooms.fields.monitorStatus.none');
    }
};

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getRooms();
};

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

const openBatchAddRoomModal = () => {
    batchAddRoomModal.value?.openModal();
};

const openEditRoomModal = (id: number) => {
    roomModal.value?.openModal(id);
};

const openConfirmModal = (id: number) => {
    roomToDelete.value = id;
    showConfirmModal.value = true;
};

const handleGoToRoomFolder = async (room: RoomInfo) => {
    let path = "";
    try {
        const res: any = await getRoomFilePath({
            platform: room.platform,
            anchor: room.anchor,
        });
        if (res.code === 0) {
            path = res.data.path;
        }
    } catch (error) {
        console.error("Failed to get room file path:", error);
    } finally {
        router.push({ path: "/media/file", query: { path } });
    }
};

const handlePreview = (liveId: number) => {
    const routeData = router.resolve({
        name: 'StreamPreview',
        query: {
            id: liveId,
        }
    });
    window.open(routeData.href, '_blank');
};

async function handleDeleteRoom() {
    if (!roomToDelete.value) {
        return;
    }
    try {
        const res: any = await deleteRoom(roomToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        getRooms();
        toast.success(t('common.toast.deleteSuccess'));
    } catch (error) {
        console.error("Failed to delete room:", error);
    } finally {
        showConfirmModal.value = false;
        roomToDelete.value = null;
    }
}

const openStartConfirmModal = (id: number) => {
    roomToStart.value = id;
    showStartConfirmModal.value = true;
};

const openStopConfirmModal = (id: number) => {
    roomToStop.value = id;
    showStopConfirmModal.value = true;
};

async function handleStartRoom() {
    if (!roomToStart.value) {
        return;
    }
    try {
        const res: any = await startRoom(roomToStart.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.operationFailed'));
            return;
        }
        getRooms();
        toast.success(t('stream.rooms.toast.startMonitor'));
    } catch (error) {
        console.error("Failed to start room:", error);
    } finally {
        showStartConfirmModal.value = false;
        roomToStart.value = null;
    }
}

async function handleStopRoom() {
    if (!roomToStop.value) {
        return;
    }
    try {
        const res: any = await stopRoom(roomToStop.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.operationFailed'));
            return;
        }
        getRooms();
        toast.success(t('stream.rooms.toast.stopMonitor'));
    } catch (error) {
        console.error("Failed to stop room:", error);
    } finally {
        showStopConfirmModal.value = false;
        roomToStop.value = null;
    }
}

async function handleTopRoom(id: number) {
    try {
        const res: any = await topRoom(id);
        if (res.code !== 0) {
            toast.error(res.msg || t('stream.rooms.toast.topErr'));
            return;
        }
        getRooms();
        toast.success(t('stream.rooms.toast.topSuccess'));
    } catch (error) {
        console.error("Failed to top room:", error);
    }
}

async function handleUnTopRoom(id: number) {
    try {
        const res: any = await unTopRoom(id);
        if (res.code !== 0) {
            toast.error(res.msg || t('stream.rooms.toast.unTopErr'));
            return;
        }
        getRooms();
        toast.success(t('stream.rooms.toast.unTopSuccess'));
    } catch (error) {
        console.error("Failed to un-top room:", error);
    }
}

const handleExport = async (exportType: number) => {
    try {
        const params = {
            sort: sort.value,
            ...filter.value,
            exportType: exportType,
        };
        const response: any = await exportRoomInfo(params);
        const ct: string = response.headers?.['content-type'] || '';
        if (ct == "application/json") {
            let text: string;
            if (response.data instanceof Blob) {
                text = await response.data.text();
            } else if (response.data instanceof ArrayBuffer) {
                text = new TextDecoder('utf-8').decode(response.data);
            } else {
                text = String(response.data);
            }
            const payload = JSON.parse(text);
            const code = payload?.code ?? -1;
            const msg = payload?.msg || t('stream.rooms.toast.exportErr');
            if (code !== 0) {
                toast.error(msg);
                return;
            }
        }
        const cd: string = response.headers?.['content-disposition'] || '';
        const star = cd.match(/filename\*=UTF-8''([^;]+)/i);
        const normal = cd.match(/filename="([^"]+)"/i);
        let fileName = star?.[1] ? decodeURIComponent(star[1]) : (normal?.[1] || 'rooms');

        const wantedExt = exportType === 1 ? '.xlsx' : '.txt';
        const wantedMime =
            exportType === 1
                ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
                : 'text/plain; charset=utf-8';

        if (!fileName.toLowerCase().endsWith(wantedExt)) {
            fileName = fileName.replace(/\.[^./\\]+$/, '');
            fileName += wantedExt;
        }

        const blob = new Blob([response.data], { type: wantedMime });

        const link = document.createElement('a');
        const url = window.URL.createObjectURL(blob);
        link.href = url;
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
    } catch (error) {
        console.error("Failed to export rooms:", error);
    }
};

const refreshInterval = ref<number | null>(null);

const clearRefreshTimers = () => {
    if (refreshInterval.value !== null) {
        clearInterval(refreshInterval.value);
        refreshInterval.value = null;
    }
};

const scheduleRefresh = () => {
    clearRefreshTimers();
    refreshInterval.value = window.setInterval(() => {
        getRooms();
    }, 20000);
};

onUnmounted(() => {
    clearRefreshTimers();
});
</script>
