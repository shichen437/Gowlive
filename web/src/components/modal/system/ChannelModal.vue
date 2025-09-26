<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>{{ isEdit ? "编辑" : "添加" }}渠道</DialogTitle>
      </DialogHeader>
      <form @submit="onSubmit" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <FormField name="name" v-slot="{ componentField }">
            <FormItem>
              <FormLabel>渠道名称</FormLabel>
              <FormControl>
                <Input type="text" placeholder="请输入渠道名称" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
          <FormField name="status" v-slot="{ value, handleChange }">
            <FormItem>
              <FormLabel>渠道状态</FormLabel>
              <FormControl>
                <div class="flex gap-2">
                  <Button type="button" class="flex-1" :variant="value === 1 ? 'default' : 'outline'"
                    @click="handleChange(1)">
                    启用
                  </Button>
                  <Button type="button" class="flex-1" :variant="value === 0 ? 'default' : 'outline'"
                    @click="handleChange(0)">
                    禁用
                  </Button>
                </div>
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>

        <Tabs :model-value="values.type" @update:model-value="(val) => setFieldValue('type', String(val))"
          default-value="email" class="w-full">
          <TabsList class="grid w-full grid-cols-2" :class="{ 'pointer-events-none opacity-50': isEdit }">
            <TabsTrigger value="email"> 邮箱 </TabsTrigger>
            <TabsTrigger value="gotify"> Gotify </TabsTrigger>
          </TabsList>
          <TabsContent value="email" class="space-y-4 mt-4 border-0 p-0">
            <div class="grid grid-cols-2 gap-4">
              <FormField name="email.server" v-slot="{ componentField }">
                <FormItem>
                  <FormLabel>SMTP服务器</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="smtp.example.com" v-bind="componentField" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
              <FormField name="email.port" v-slot="{ componentField }">
                <FormItem>
                  <FormLabel>端口</FormLabel>
                  <FormControl>
                    <Input type="number" placeholder="465" v-bind="componentField" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
            <FormField name="email.sender" v-slot="{ componentField }">
              <FormItem>
                <FormLabel>发送人</FormLabel>
                <FormControl>
                  <Input type="email" placeholder="sender@example.com" v-bind="componentField" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            <FormField name="email.authCode" v-slot="{ componentField }">
              <FormItem>
                <FormLabel>授权码</FormLabel>
                <FormControl>
                  <Input type="password" placeholder="请输入授权码" v-bind="componentField" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            <FormField name="email.receiver" v-slot="{ componentField }">
              <FormItem>
                <FormLabel>收件人</FormLabel>
                <FormControl>
                  <Input type="email" placeholder="receiver@example.com" v-bind="componentField" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </TabsContent>
          <TabsContent value="gotify" class="space-y-4 mt-4 border-0 p-0">
            <FormField name="url" v-slot="{ componentField }">
              <FormItem>
                <FormLabel>Gotify URL</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="https://gotify.example.com" v-bind="componentField" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </TabsContent>
        </Tabs>

        <FormField name="remark" v-slot="{ componentField }">
          <FormItem>
            <FormLabel>备注</FormLabel>
            <FormControl>
              <Input placeholder="请输入备注(可选)" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter>
          <Button type="button" variant="outline" @click="open = false">取消</Button>
          <Button type="submit">保存</Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { watch, computed } from "vue";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  addPushChannel,
  updatePushChannel,
  getPushChannel,
} from "@/api/system/push_channel";
import type { PushChannel } from "@/types/system";
import { toast } from "vue-sonner";

const props = defineProps<{
  modelValue: boolean;
  channel: PushChannel | null;
}>();

const emit = defineEmits(["update:modelValue", "success"]);

const open = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value),
});

const isEdit = computed(() => !!props.channel);

const formSchema = toTypedSchema(
  z
    .object({
      name: z.string().min(1, "渠道名称不能为空"),
      type: z.string().min(1, "渠道类型不能为空"),
      status: z.coerce.number(),
      remark: z.string().optional(),
      url: z.string().optional(),
      email: z
        .object({
          sender: z.string().optional(),
          receiver: z.string().optional(),
          server: z.string().optional(),
          port: z.coerce.number().optional(),
          authCode: z.string().optional(),
        })
        .optional(),
    })
    .superRefine((data, ctx) => {
      if (data.type === "gotify") {
        if (!data.url || data.url.trim() === "") {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Gotify URL不能为空",
            path: ["url"],
          });
        } else {
          const urlResult = z
            .string()
            .url({ message: "请输入有效的URL" })
            .safeParse(data.url);
          if (!urlResult.success) {
            ctx.addIssue({
              code: z.ZodIssueCode.custom,
              message: urlResult.error.issues[0].message,
              path: ["url"],
            });
          }
        }
      }
    })
);

const { handleSubmit, resetForm, values, setFieldValue, setValues } = useForm({
  validationSchema: formSchema,
});

watch(open, async (isOpen) => {
  if (isOpen) {
    if (isEdit.value && props.channel) {
      try {
        const res: any = await getPushChannel(props.channel.id);
        if (res.code !== 0) {
          toast.error(res.msg || "获取详情失败")
        }
        setValues(res.data);
      } catch (error) {
        console.error("Failed to fetch channel details:", error);
      }
    } else {
      resetForm({
        values: { status: 1, type: "email", email: {}, url: "" },
      });
    }
  } else {
    resetForm({
      values: { status: 1, type: "email", email: {}, url: "" },
    });
  }
});

const onSubmit = handleSubmit(async (formValues) => {
  try {
    if (isEdit.value && props.channel) {
      await updatePushChannel({ ...formValues, id: props.channel.id });
    } else {
      await addPushChannel(formValues);
    }
    emit("success");
    open.value = false;
  } catch (error) {
    console.error("Failed to save channel:", error);
  }
});
</script>