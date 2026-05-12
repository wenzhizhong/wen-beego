<template>
  <div class="main">
    <el-card shadow="never">
      <div class="search-box">
        <el-form :inline="true" :model="searchForm" ref="searchFormRef">
          <el-form-item label="关键词">
            <el-input v-model="searchForm.keyword" placeholder="请输入关键词" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="useRenderIcon('search')" @click="onSearch">查询</el-button>
            <el-button :icon="useRenderIcon('refresh')" @click="onReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <el-card shadow="never" class="mt-4">
      <div class="mb-4">
        <el-button type="primary" :icon="useRenderIcon('add')" @click="openDialog()">新增</el-button>
        <el-button type="danger" :icon="useRenderIcon('delete')" @click="handleBatchDelete" :disabled="!selectedRows.length">批量删除</el-button>
      </div>

      <vxe-grid
        ref="xGridRef"
        v-bind="gridOptions"
        :data="dataList"
        :loading="loading"
        :checkbox-config="{}"
        @checkbox-change="handleSelectionChange"
        @page-change="handlePageChange"
      >
        <vxe-column type="checkbox" width="50" />
        {{range .Columns}}<vxe-column field="{{.Name}}" title="{{.Comment}}" />
        {{end}}
        <vxe-column title="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </vxe-column>
      </vxe-grid>
    </el-card>

    <{{.ModelName}}Form ref="{{.ModelNameLower}}FormRef" @refresh="onSearch" />
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import {{.ModelName}}Form from "./form.vue";
import { get{{.ModelName}}List, del{{.ModelName}} } from "./api";
import { message, messageBox } from "@/utils/message";

const {{.ModelNameLower}}FormRef = ref();
const loading = ref(false);
const dataList = ref([]);
const selectedRows = ref([]);

const searchForm = reactive({
  keyword: ""
});

const gridOptions = reactive({
  border: true,
  stripe: true,
  height: "auto",
  pager: {
    currentPage: 1,
    pageSize: 10,
    total: 0
  },
  columns: [] as any[]
});

onMounted(() => {
  onSearch();
});

async function onSearch() {
  loading.value = true;
  try {
    const res = await get{{.ModelName}}List({
      currentPage: gridOptions.pager.currentPage,
      pageSize: gridOptions.pager.pageSize,
      keyword: searchForm.keyword
    });
    const { list, total } = res.data ?? { list: [], total: 0 };
    dataList.value = list ?? [];
    gridOptions.pager.total = total ?? 0;
    selectedRows.value = [];
  } finally {
    loading.value = false;
  }
}

function onReset() {
  searchForm.keyword = "";
  onSearch();
}

function openDialog(row?: any) {
  {{.ModelNameLower}}FormRef.value?.open(row?.id);
}

function handleSelectionChange({ records }: any) {
  selectedRows.value = records;
}

function handlePageChange({ currentPage, pageSize }: any) {
  gridOptions.pager.currentPage = currentPage;
  gridOptions.pager.pageSize = pageSize;
  onSearch();
}

async function handleDelete(row: any) {
  await messageBox.confirm("确认删除该记录？");
  await del{{.ModelName}}({ id: row.id });
  message.success("删除成功");
  onSearch();
}

async function handleBatchDelete() {
  if (!selectedRows.value.length) {
    message.warning("请选择要删除的记录");
    return;
  }
  const ids = selectedRows.value.map((item: any) => item.id);
  await messageBox.confirm("确认删除选中的 " + ids.length + " 条记录？");
  await del{{.ModelName}}({ ids });
  message.success("批量删除成功");
  onSearch();
}
</script>

<style scoped>
.main {
  padding: 0;
}
.search-box {
  display: flex;
  align-items: center;
}
</style>
