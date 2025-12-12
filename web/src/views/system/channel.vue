<template>
    <div class="space-y-4">
        <div class="flex justify-between items-center">
            <Button @click="handleAdd">
                <Plus class="w-4 h-4 mr-2" />
                {{ t('system.channel.add.button') }}
            </Button>
        </div>
        <div class="border rounded-lg">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead class="text-center">{{ t('system.channel.fields.name') }}</TableHead>
                        <TableHead class="text-center">{{ t('system.channel.fields.type') }}</TableHead>
                        <TableHead class="text-center">{{ t('system.channel.fields.status') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.fields.remark') }}</TableHead>
                        <TableHead class="text-center">{{ t('common.operation.title') }}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <template v-if="channels.length === 0">
                        <TableRow>
                            <TableCell :colspan="5" class="h-24 text-center">
                                {{ t('common.noData') }}
                            </TableCell>
                        </TableRow>
                    </template>
                    <template v-else>
                        <TableRow v-for="channel in channels" :key="channel.id">
                            <TableCell class="text-center">{{ channel.name }}</TableCell>
                            <TableCell class="text-center">{{ getChannelTypeName(channel.type) }}</TableCell>
                            <TableCell class="text-center">
                                <Badge variant="outline"
                                    :class="channel.status === 1 ? 'text-green-600' : 'text-red-600'">
                                    {{ channel.status === 1 ? t('system.channel.fields.statusActive') : "禁用" }}
                                </Badge>
                            </TableCell>
                            <TableCell class="text-center">{{ channel.remark }}</TableCell>
                            <TableCell class="text-center space-x-2">
                                <Button variant="ghost" size="icon" @click="handleEdit(channel)">
                                    <Pencil class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive hover:text-destructive"
                                    @click="openDeleteConfirm(channel)">
                                    <Trash2 class="w-4 h-4" />
                                </Button>
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

        <ChannelModal v-model="isChannelModalOpen" :channel="selectedChannel" @success="handleModalSuccess" />

        <ConfirmModal :open="isConfirmModalOpen" :onOpenChange="(v) => (isConfirmModalOpen = v)"
            :onConfirm="handleDeleteConfirm" :title="t('common.operation.deleteConfirm')"
            :description="t('system.channel.deleteDesc')" />
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { Button } from "@/components/ui/button";
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
import { Badge } from "@/components/ui/badge";
import { Pencil, Trash2, Plus } from "lucide-vue-next";
import { listPushChannel, deletePushChannel } from "@/api/system/push_channel";
import type { PushChannel } from "@/types/system";
import { useDict } from "@/utils/useDict";
import ChannelModal from "@/components/modal/system/ChannelModal.vue";
import ConfirmModal from "@/components/modal/ConfirmModal.vue";
import { toast } from "vue-sonner";

const { t } = useI18n();
const { dict } = useDict("channel_type");

const channels = ref<PushChannel[]>([]);
const loading = ref(true);
const pageNum = ref(1);
const pageSize = ref(10);
const total = ref(0);

const isChannelModalOpen = ref(false);
const selectedChannel = ref<PushChannel | null>(null);

const isConfirmModalOpen = ref(false);
const channelToDelete = ref<number | null>(null);

async function getChannels() {
    loading.value = true;
    try {
        const params = {
            pageNum: pageNum.value,
            pageSize: pageSize.value,
        };
        const res: any = await listPushChannel(params);
        if (res.code !== 0) {
            toast.error(res.msg || t('system.channel.toast.listErr'));
        }
        channels.value = res.data.rows || [];
        total.value = res.data.total || 0;
    } catch (error) {
        console.error("Error fetching push channels:", error);
    } finally {
        loading.value = false;
    }
}

function getChannelTypeName(type: string) {
    const found = dict.value.find((d: any) => d.dictValue === type);
    return found ? found.dictLabel : type;
}

const handlePageChange = (newPage: number) => {
    pageNum.value = newPage;
    getChannels();
}

function handleAdd() {
    selectedChannel.value = null;
    isChannelModalOpen.value = true;
}

function handleEdit(channel: PushChannel) {
    selectedChannel.value = channel;
    isChannelModalOpen.value = true;
}

function openDeleteConfirm(channel: PushChannel) {
    channelToDelete.value = channel.id;
    isConfirmModalOpen.value = true;
}

async function handleDeleteConfirm() {
    if (!channelToDelete.value) return;
    try {
        const res: any = await deletePushChannel(channelToDelete.value);
        if (res.code !== 0) {
            toast.error(res.msg || t('common.toast.deleteFailed'));
            return;
        }
        getChannels();
        toast.success(t('common.toast.deleteSuccess'));
    } catch (error) {
        console.error("Failed to delete channel:", error);
    } finally {
        isConfirmModalOpen.value = false;
        channelToDelete.value = null;
    }
}

function handleModalSuccess() {
    getChannels();
}

onMounted(() => {
    getChannels();
});
</script>
