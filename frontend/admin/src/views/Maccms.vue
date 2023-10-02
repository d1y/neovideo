<template>
  <div class="p-5">
    <el-button-group>
      <el-button class="pl-2" @click="create">创建线路</el-button>
      <el-button @click="createBybatchImport">批量导入</el-button>
    </el-button-group>
    <el-table ref="tableRef" :data="tableData">
      <el-table-column label="可用" width="60">
        <template #default="{row}">
          <div class="w-[24px] h-[24px] rounded-[24px]" :style="{
            background: row.available ? '#4d70ff' : '#ff4d4f'
          }"></div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="cms名称" width="120" />
      <el-table-column prop="api" label="api" width="420" />
      <el-table-column label="操作" width="320">
        <template #default="scope">
          <el-button>检测可用性</el-button>
          <el-button>修改</el-button>
          <el-button @click="del(scope.$index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage, ElMessageBox, type TableInstance } from "element-plus";
import * as maccmsApi from "@/api/maccms";
import { onMounted } from "vue";
import { MacCMSRepo } from "@/types/maccms";

const tableRef = ref<TableInstance>();

const tableData = ref<MacCMSRepo[]>([]);

async function getData() {
  const resp = await maccmsApi.getList();
  tableData.value = resp.data;
}

async function createBybatchImport() {
}

async function create() {
}

async function del(idx: number) {
  const item = tableData.value[idx];
  const resp = await maccmsApi.del(item.id!);
  ElMessage({
    duration: 1200,
    message: resp.success ? "删除成功" : "删除成功",
  });
  if (resp.success) {
    tableData.value.splice(idx, 1);
  }
  await getData();
}

onMounted(getData);
</script>
