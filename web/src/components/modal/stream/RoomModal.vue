<template>
  <Dialog :open="isOpen" @update:open="handleClose">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>{{ dialogTitle }}</DialogTitle>
      </DialogHeader>
      <form @submit="onSubmit">
        <div class="grid gap-4 py-4">
          <FormField v-slot="{ componentField, errorMessage }" name="roomUrl">
            <FormItem>
              <FormLabel>直播链接</FormLabel>
              <FormControl>
                <Input type="text" placeholder="请输入直播链接" v-bind="componentField" :disabled="isEditMode"
                  :class="{ 'border-red-500': errorMessage }" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <div class="grid grid-cols-2 gap-4">
            <FormField v-slot="{ componentField, errorMessage }" name="interval">
              <FormItem>
                <FormLabel>间隔时间 (秒)</FormLabel>
                <FormControl>
                  <Input type="number" placeholder="30-600" v-bind="componentField"
                    :class="{ 'border-red-500': errorMessage }" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField, errorMessage }" name="format">
              <FormItem>
                <FormLabel>录制格式</FormLabel>
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger class="w-full" :class="{ 'border-red-500': errorMessage }">
                      <SelectValue placeholder="选择录制格式" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent class="w-[--radix-select-trigger-width]">
                    <SelectItem value="flv">FLV</SelectItem>
                    <SelectItem value="mp4">MP4</SelectItem>
                    <SelectItem value="mp3">MP3(仅音频)</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>

          <FormField v-slot="{ field, errorMessage }" name="monitorType">
            <FormItem>
              <FormLabel class="flex items-center">
                监控类型
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Info class="w-4 h-4 ml-1" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>实时监控：根据设置的间隔时间轮询获取直播状态。</p>
                      <p>定时监控：在指定时间段内根据间隔时间获取直播状态。</p>
                      <p>智能监控：根据有效直播历史动态调整间隔时间大小，暂无历史时降级为实时监控。</p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </FormLabel>
              <FormControl>
                <RadioGroup class="grid grid-cols-4 gap-2" :model-value="field.value"
                  @update:model-value="field.onChange">
                  <div v-for="item in monitorTypeOptions" :key="item.value">
                    <RadioGroupItem :id="`monitor-type-${item.value}`" :value="item.value" class="sr-only" />
                    <label :for="`monitor-type-${item.value}`" :class="cn(
                      buttonVariants({ variant: field.value === item.value ? 'default' : 'outline' }),
                      'w-full cursor-pointer flex items-center justify-center h-9',
                      errorMessage ? 'border-red-500' : ''
                    )">
                      {{ item.label }}
                    </label>
                  </div>
                </RadioGroup>
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <template v-if="values.monitorType === 2">
            <div class="grid grid-cols-2 gap-4">
              <FormField v-slot="{ componentField, errorMessage }" name="monitorStartAt">
                <FormItem>
                  <FormLabel>监控开始时间</FormLabel>
                  <FormControl>
                    <VueDatePicker v-bind="componentField" time-picker auto-apply model-type="HH:mm" :format="'HH:mm'"
                      text-input placeholder="HH:mm" :class="{ 'border-red-500': errorMessage }" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
              <FormField v-slot="{ componentField, errorMessage }" name="monitorStopAt">
                <FormItem>
                  <FormLabel>监控结束时间</FormLabel>
                  <FormControl>
                    <VueDatePicker v-bind="componentField" time-picker auto-apply model-type="HH:mm" :format="'HH:mm'"
                      text-input placeholder="HH:mm" :class="{ 'border-red-500': errorMessage }" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </template>

          <FormField v-slot="{ componentField, errorMessage }" name="remark">
            <FormItem>
              <FormLabel>备注</FormLabel>
              <FormControl>
                <Input type="text" placeholder="请输入备注" v-bind="componentField"
                  :class="{ 'border-red-500': errorMessage }" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>
        <DialogFooter>
          <Button type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? '保存中...' : '保存' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose, defineEmits, computed } from 'vue';
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import { useForm } from 'vee-validate';
import { addRoom, updateRoom, roomDetail } from '@/api/stream/live_manage';
import { toast } from 'vue-sonner';
import { cn } from '@/lib/utils';
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
import type { LiveManage } from '@/types/stream';

import { Button, buttonVariants } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Info } from 'lucide-vue-next';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/ui/tooltip';
import { Input } from '@/components/ui/input';

const isOpen = ref(false);
const emit = defineEmits(['refresh']);
const roomId = ref<number | null>(null);
const isSubmitting = ref(false);

const isEditMode = computed(() => !!roomId.value);
const dialogTitle = computed(() => isEditMode.value ? '编辑房间' : '添加房间');

const monitorTypeOptions = [
  { value: 0, label: '停止监控' },
  { value: 1, label: '实时监控' },
  { value: 2, label: '定时监控' },
  { value: 3, label: '智能监控' }
];

const formSchema = toTypedSchema(
  z.object({
    roomUrl: z.string().url({ message: '请输入有效的URL' }),
    interval: z.coerce.number().min(30, { message: '间隔时间最小为30秒' }).max(600, { message: '间隔时间最大为600秒' }),
    format: z.enum(['flv', 'mp4', 'mp3']),
    monitorType: z.coerce.number(),
    monitorStartAt: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, { message: '时间格式必须为 HH:mm' }).optional().nullable(),
    monitorStopAt: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, { message: '时间格式必须为 HH:mm' }).optional().nullable(),
    remark: z.string().max(45, { message: '备注最长为45个字符' }).optional().nullable(),
  }).refine(data => {
    if (data.monitorType === 2) {
      return !!data.monitorStartAt && !!data.monitorStopAt;
    }
    return true;
  }, {
    message: '定时监控需要设置开始和结束时间',
    path: ['monitorStartAt'],
  })
);

const { handleSubmit, values, setValues, resetForm } = useForm({
  validationSchema: formSchema,
  initialValues: {
    roomUrl: '',
    interval: 30,
    format: 'flv',
    monitorType: 0,
    remark: '',
    monitorStartAt: null,
    monitorStopAt: null,
  }
});

const onSubmit = handleSubmit(async (formValues) => {
  if (isSubmitting.value) return;

  isSubmitting.value = true;

  try {
    const payload = {
      ...formValues,
      monitorType: Number(formValues.monitorType)
    };

    if (isEditMode.value && roomId.value) {
      const res: any = await updateRoom({ id: roomId.value, ...payload });
      if (res.code !== 0) {
        toast.error(res.data.msg || '更新失败');
        return;
      }
      toast.success('更新成功');
    } else {
      const res: any = await addRoom(payload);
      if (res.code !== 0) {
        toast.error(res.data.msg || '添加失败');
        return;
      }
      toast.success('添加成功');
    }

    isOpen.value = false;
    emit('refresh');
  } catch (error) {
    console.error('Failed to save room:', error);
  } finally {
    isSubmitting.value = false;
  }
});

const openModal = async (id?: number) => {
  resetForm();
  if (id) {
    roomId.value = id;
    try {
      const response = await roomDetail(id);
      const roomData = response.data.data as LiveManage;

      const formattedData = {
        ...roomData,
        monitorType: Number(roomData.monitorType),
        monitorStartAt: roomData.monitorStartAt || null,
        monitorStopAt: roomData.monitorStopAt || null,
        remark: roomData.remark || ''
      };

      setValues(formattedData);
    } catch (error) {
      console.error('Failed to fetch room details:', error);
      toast.error('获取房间详情失败');
      return;
    }
  } else {
    roomId.value = null;
  }
  isOpen.value = true;
};

const handleClose = () => {
  isOpen.value = false;
};

defineExpose({
  openModal,
});
</script>