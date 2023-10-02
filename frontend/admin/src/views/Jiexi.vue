<template>
  <div class="p-5">
    <el-button-group>
      <el-button class="pl-2" @click="create">创建线路</el-button>
      <el-button @click="createBybatchImport">批量导入</el-button>
    </el-button-group>
    <el-table ref="tableRef" :data="tableData">
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
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage, ElMessageBox, type TableInstance } from "element-plus";
import * as jiexiApi from "@/api/jiexi";
import { onMounted } from "vue";
import jiexiForm from "@/components/JiexiForm.vue";
import jiexiImport from "@/components/JiexiImport.vue";
import { JiexiTable } from "@t/jiexi";
import { h } from "vue";

const tableRef = ref<TableInstance>();

const tableData = ref<JiexiTable[]>([]);

async function getData() {
  const resp = await jiexiApi.getList();
  tableData.value = resp.data;
}

async function createBybatchImport() {
  const val = await new Promise<false | string>((res) => {
    ElMessageBox({
      title: "批量导入",
      showCancelButton: true,
      showConfirmButton: true,
      cancelButtonText: "取消",
      confirmButtonText: "确定",
      message: h(jiexiImport, {
        data: "",
      }),
      async beforeClose(action, instance, done) {
        if (action == "cancel" || action == "close") {
          done();
          res(false);
          return;
        } else if (action == "confirm") {
          const m = instance.message as any;
          const val = m.component.exposed.getValue() as string;
          res(val);
          done();
        }
      },
    });
  });
  if (!val) {
    ElMessage({
      message: "空内容导入失败",
    });
    return;
  }
  const resp = await jiexiApi.batchImport(val);
  ElMessage({
    type: resp.success ? "success" : "error",
    message: resp.message,
  });
  await getData();
}

async function create() {
  const form = await new Promise<false | JiexiTable>(async (res) => {
    await ElMessageBox({
      title: "新建线路",
      showCancelButton: true,
      showConfirmButton: true,
      cancelButtonText: "取消",
      confirmButtonText: "确定",
      message: h(jiexiForm),
      async beforeClose(action, instance, done) {
        if (action == "cancel" || action == "close") {
          res(false);
          done();
          return;
        } else if (action == "confirm") {
          const m = instance.message as any;
          const status = await m.component.exposed.validate();
          if (!status) return;
          const data = m.component.exposed.getForm();
          res(data);
          done();
        }
      },
    });
  });
  if (!form) return;
  const createResp = await jiexiApi.create(form);
  ElMessage({
    duration: 1200,
    type: createResp.success ? "success" : "error",
    message: createResp.message,
  });
}

async function del(idx: number) {
  const item = tableData.value[idx];
  const resp = await jiexiApi.del(item.id!);
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
