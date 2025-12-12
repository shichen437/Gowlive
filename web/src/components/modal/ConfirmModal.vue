<script setup lang="ts">
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import type { PropType } from 'vue';
import { Button } from "@/components/ui/button";
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps({
    open: {
        type: Boolean,
        default: false,
    },
    onOpenChange: {
        type: Function as PropType<(open: boolean) => void>,
        default: () => { },
    },
    onConfirm: {
        type: Function as PropType<() => void>,
        default: () => { },
    },
    title: {
        type: String,
        required: true,
    },
    description: {
        type: String,
        required: true,
    },
})

const handleConfirm = () => {
    props.onConfirm();
    props.onOpenChange(false);
}
</script>

<template>
    <Dialog :open="open" @update:open="onOpenChange">
        <DialogContent>
            <DialogHeader>
                <DialogTitle>{{ title }}</DialogTitle>
                <DialogDescription>
                    {{ description }}
                </DialogDescription>
            </DialogHeader>
            <DialogFooter>
                <Button variant="outline" @click="onOpenChange(false)">{{ t('common.operation.cancel') }}</Button>
                <Button @click="handleConfirm">{{ t('common.operation.confirm') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>
