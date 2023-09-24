<template>
  <el-form :model="ruleForm" :rules="rules" ref="formRef" label-width="100px" style="max-width: 460px">
    <el-form-item label="线路名称">
      <el-input v-model="ruleForm.name" />
    </el-form-item>
    <el-form-item label="线路链接">
      <el-input v-model="ruleForm.url" />
    </el-form-item>
    <el-form-item label="说明">
      <el-input v-model="ruleForm.note" />
    </el-form-item>
    <div class="flex justify-center items-center">
      <el-button>取消</el-button>
      <el-button @click="confirm">确定</el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { reactive } from 'vue'

const props = defineProps<Partial<JiexiTable>>()

const formRef = ref<FormInstance>()

const ruleForm = reactive(props)

const rules = reactive<FormRules<typeof ruleForm>>({
  name: [
    {
      required: true,
    }
  ],
  url: [
    {
      required: true,
    }
  ],
})

async function confirm() {
  const isNext = await formRef.value?.validate()
  debugger
}

</script>