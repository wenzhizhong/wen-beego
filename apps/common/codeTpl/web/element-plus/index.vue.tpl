<script setup lang="ts">
import { ref } from "vue";
import { use{{.ModelName}} } from "./utils/hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import {{.ModelName}}Form from "./form.vue";

import Delete from "~icons/ep/delete";
import EditPen from "~icons/ep/edit-pen";
import Refresh from "~icons/ep/refresh";
import AddFill from "~icons/ri/add-circle-line";

defineOptions({
  name: "{{.MenuModule}}_{{.ModelName}}"
});

const formRef = ref();
const tableRef = ref();
const {
  loading,
  dataList,
  selectedNum,
  searchForm,
  pagination,
  columns,
  onSearch,
  resetForm,
  handleSizeChange,
  handleCurrentChange,
  handleSelectionChange,
  onSelectionCancel,
  openDialog,
  handleDelete,
  onBatchDel
} = use{{.ModelName}}(formRef);
</script>

<template>
  <div class="main">
    <el-form
      ref="searchFormRef"
      :inline="true"
      :model="searchForm"
      class="search-form bg-bg_color w-full pl-8 pt-[12px]"
    >
      <el-form-item label="关键词：">
        <el-input
          v-model="searchForm.keyword"
          placeholder="请输入关键词"
          clearable
          class="!w-[200px]"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :icon="useRenderIcon('ri/search-line')" :loading="loading" @click="onSearch">搜索</el-button>
        <el-button :icon="useRenderIcon(Refresh)" @click="resetForm(formRef)">重置</el-button>
      </el-form-item>
    </el-form>

    <PureTableBar title="{{.MenuName}}" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button type="primary" :icon="useRenderIcon(AddFill)" @click="openDialog()">新增</el-button>
        <el-button type="danger" :icon="useRenderIcon(Delete)" :disabled="selectedNum === 0" @click="onBatchDel(tableRef?.getTableRef()?.getSelectionRows())">批量删除</el-button>
      </template>
      <template v-slot="{ size, dynamicColumns }">
        <div
          v-if="selectedNum > 0"
          v-motion-fade
          class="bg-[var(--el-fill-color-light)] w-full h-[46px] mb-2 pl-4 flex items-center"
        >
          <div class="flex-auto">
            <span class="text-[rgba(42,46,54,0.5)] dark:text-[rgba(220,220,242,0.5)]">
              已选 {{"{{"}} selectedNum {{"}}"}} 项
            </span>
            <el-button type="primary" text @click="onSelectionCancel">取消选择</el-button>
          </div>
        </div>
        <pure-table
          ref="tableRef"
          row-key="id"
          adaptive
          :adaptiveConfig="{ offsetBottom: 108 }"
          align-whole="center"
          table-layout="auto"
          :loading="loading"
          :size="size"
          :data="dataList"
          :columns="dynamicColumns"
          :pagination="{ ...pagination, size }"
          :header-cell-style="{
            background: 'var(--el-fill-color-light)',
            color: 'var(--el-text-color-primary)'
          }"
          @selection-change="handleSelectionChange"
          @page-size-change="handleSizeChange"
          @page-current-change="handleCurrentChange"
        >
          <template #operation="{ row }">
            <el-button class="reset-margin" link type="primary" :size="size" :icon="useRenderIcon(EditPen)" @click="openDialog('编辑', row)">编辑</el-button>
            <el-button class="reset-margin" link type="danger" :size="size" :icon="useRenderIcon(Delete)" @click="handleDelete(row)">删除</el-button>
          </template>
        </pure-table>
      </template>
    </PureTableBar>

    <{{.ModelName}}Form ref="formRef" @refresh="onSearch" />
  </div>
</template>

<style scoped>
.main {
  margin: 0;
}
.search-form :deep(.el-form-item) {
  margin-bottom: 12px;
}
</style>
