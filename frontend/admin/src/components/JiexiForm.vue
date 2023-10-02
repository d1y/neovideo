<template>
  <el-form
    :model="ruleForm"
    :rules="rules"
    ref="formRef"
    label-width="100px"
    style="max-width: 460px"
  >
    <el-form-item prop="name" label="线路名称">
      <el-input v-model="ruleForm.name" />
    </el-form-item>
    <el-form-item prop="url" label="线路链接">
      <el-input v-model="ruleForm.url" />
    </el-form-item>
    <el-form-item prop="note" label="说明">
      <el-input v-model="ruleForm.note" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { reactive } from "vue";
import { JiexiTable } from "@t/jiexi";

const props = defineProps<Partial<JiexiTable>>();

const formRef = ref<FormInstance>();

const ruleForm = reactive<JiexiTable>({
  name: props.name || "",
  url: props.url || "",
  note: props.note || "",
});

const rules = reactive<FormRules<typeof ruleForm>>({
  name: [{ required: true, message: "name is required", trigger: "blur" }],
  url: [{ required: true, message: "url is required", trigger: "blur" }],
});

defineExpose({
  async validate() {
    return await formRef.value?.validate();
  },
  getForm() {
    return ruleForm;
  },
});
</script>
