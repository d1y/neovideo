<template>
  <el-button-group>
    <el-button @click="create">创建线路</el-button>
    <el-button>批量导入</el-button>
  </el-button-group>
  <el-table ref="tableRef" row-key="date" :data="tableData" style="width: 100%">
    <el-table-column prop="name" label="线路名称" width="120" />
    <el-table-column prop="url" label="线路链接" width="320" />
    <el-table-column prop="note" label="说明" width="auto" />
    <el-table-column label="操作" width="320">
      <template #default="scope">
        <el-button>修改</el-button>
        <el-button @click="del(scope.$index)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox, type TableInstance } from 'element-plus'
import * as jiexiApi from '@/api/jiexi'
import { onMounted } from 'vue'
import jiexiForm from './jiexi_form.vue'
import { JiexiTable } from '@t/jiexi'
import { h } from 'vue'

const tableRef = ref<TableInstance>()

const tableData = ref<JiexiTable[]>([])

async function getData() {
  const resp = await jiexiApi.getList()
  tableData.value = [...resp.data, {
    name: "d1y",
    url: "test",
    note: "白嫖",
    id: 12,
  }]
}

async function create() {
  ElMessageBox({
    title: "新建线路",
    showConfirmButton: false,
    message: h(jiexiForm),
  })
}

async function del(idx: number) {
  const item = tableData.value[idx]
  const resp = await jiexiApi.del(item.id!)
  ElMessage({
    duration: 1200,
    message: resp.success ? '删除成功' : '删除成功',
  })
  if (resp.success) {
    tableData.value.splice(idx, 1)
  }
}

onMounted(getData)

</script>