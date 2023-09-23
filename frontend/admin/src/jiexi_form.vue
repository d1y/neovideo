<template>
  <el-form ref="formRef" label-width="100px" style="max-width: 460px">
    <el-form-item label="线路名称">
      <el-input v-model="name" />
    </el-form-item>
    <el-form-item label="线路链接">
      <el-input v-model="url" />
    </el-form-item>
    <el-form-item label="说明">
      <el-input v-model="note" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { reactive } from 'vue';

defineProps<{
  name: string,
  url: string,
  note: string,
}>()

// const rules: FormRules = reactive<FormRules<typeof ruleForm>>({
//   pass: [{ validator: validatePass, trigger: 'blur' }],
//   checkPass: [{ validator: validatePass2, trigger: 'blur' }],
//   age: [{ validator: checkAge, trigger: 'blur' }],
// })

const formRef = ref<FormInstance>()
</script>