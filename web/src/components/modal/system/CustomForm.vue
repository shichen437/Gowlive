<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Info } from 'lucide-vue-next';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const methodOptions = [
  { value: 0, label: 'GET' },
  { value: 1, label: 'POST' },
];
</script>

<template>
  <div class="space-y-4 mt-4 border-0 p-0">
    <FormField name="custom.webhookUrl" v-slot="{ componentField }">
      <FormItem>
        <FormLabel
          >{{ t('system.channel.fields.webhookUrl') }}
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger as-child>
                <Info class="w-4 h-4 ml-1" />
              </TooltipTrigger>
              <TooltipContent>
                <p>{{ t('system.channel.tooltip.pre') }}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider></FormLabel
        >
        <FormControl>
          <Input
            type="text"
            :placeholder="t('system.channel.placeholder.webhookUrl')"
            v-bind="componentField"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <div class="grid grid-cols-1 gap-4">
      <FormField name="custom.requestMethod" v-slot="{ value, handleChange }">
        <FormItem>
          <FormLabel>{{ t('system.channel.fields.requestMethod') }}</FormLabel>
          <Select
            :model-value="value !== undefined ? String(value) : undefined"
            @update:model-value="handleChange(Number($event))"
          >
            <FormControl>
              <SelectTrigger class="w-full">
                <SelectValue
                  :placeholder="t('system.channel.placeholder.requestMethod')"
                />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectItem
                v-for="opt in methodOptions"
                :key="opt.value"
                :value="String(opt.value)"
              >
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <FormField name="custom.requestHeaders" v-slot="{ componentField }">
      <FormItem>
        <FormLabel>{{ t('system.channel.fields.requestHeaders') }}</FormLabel>
        <FormControl>
          <Textarea
            :placeholder="t('system.channel.placeholder.requestHeaders')"
            v-bind="componentField"
            rows="3"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="custom.requestBody" v-slot="{ componentField }">
      <FormItem>
        <FormLabel>{{ t('system.channel.fields.requestBody') }}</FormLabel>
        <FormControl>
          <Textarea
            :placeholder="t('system.channel.placeholder.requestBody')"
            v-bind="componentField"
            rows="5"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
  </div>
</template>
