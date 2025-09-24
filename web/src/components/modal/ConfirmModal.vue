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

const props = defineProps({
    open: {
        type: Boolean,
        default: false,
    },
    onOpenChange: {
        type: Function as PropType<(open: boolean) => void>,
        default: () => {},
    },
    onConfirm: {
        type: Function as PropType<() => void>,
        default: () => {},
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
                <Button variant="outline" @click="onOpenChange(false)">取消</Button>
                <Button @click="handleConfirm">确定</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>