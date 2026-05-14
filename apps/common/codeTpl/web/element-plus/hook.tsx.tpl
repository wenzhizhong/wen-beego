import { reactive, ref, onMounted, nextTick } from "vue";
import dayjs from "dayjs";
import { message } from "@/utils/message";
import { ElMessageBox } from "element-plus";
import { get{{.ModelName}}List, del{{.ModelName}} } from "@/api/{{.ModelNameLower}}";
import type { PaginationProps } from "@pureadmin/table";
import type { FormItemProps } from "./types";

export function use{{.ModelName}}(formRef: any) {
  const loading = ref(false);
  const dataList = ref([]);
  const selectedNum = ref(0);

  const searchForm = reactive({
    keyword: ""
  });

  const pagination = reactive<PaginationProps>({
    total: 0,
    pageSize: 10,
    currentPage: 1,
    background: true
  });

  const columns: TableColumnList = [
    {
      label: "勾选列",
      type: "selection",
      fixed: "left",
      reserveSelection: true
    },
{{range .Columns}}    {
      label: "{{.Comment}}",
      prop: "{{.Name}}",
      minWidth: 120
    },
{{end}}    {
      label: "创建时间",
      prop: "created_at",
      minWidth: 180,
      formatter: ({ created_at }) => created_at ? dayjs(created_at).format("YYYY-MM-DD HH:mm:ss") : ""
    },
    {
      label: "操作",
      fixed: "right",
      width: 180,
      slot: "operation"
    }
  ];

  onMounted(() => {
    onSearch();
  });

  async function onSearch() {
    loading.value = true;
    try {
      const res = await get{{.ModelName}}List({
        currentPage: pagination.currentPage,
        pageSize: pagination.pageSize,
        keyword: searchForm.keyword
      });
      const { list, total } = res.data ?? { list: [], total: 0 };
      dataList.value = list ?? [];
      pagination.total = total ?? 0;
      selectedNum.value = 0;
    } finally {
      loading.value = false;
    }
  }

  function resetForm(formEl: any) {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  }

  function handleSizeChange(val: number) {
    pagination.pageSize = val;
    onSearch();
  }

  function handleCurrentChange(val: number) {
    pagination.currentPage = val;
    onSearch();
  }

  function handleSelectionChange(val: any) {
    selectedNum.value = val.length;
  }

  function onSelectionCancel() {
    selectedNum.value = 0;
  }

  function openDialog(title = "新增", row?: FormItemProps) {
    nextTick(() => {
      formRef.value?.open(title, row);
    });
  }

  async function handleDelete(row: any) {
    try {
      await ElMessageBox.confirm("确认删除该记录？");
      const res = await del{{.ModelName}}({ id: row.id });
      if (res.code === 200) {
        message("删除成功", { type: "success" });
        onSearch();
      } else {
        message(res.message || "删除失败", { type: "error" });
      }
    } catch {
      // cancelled
    }
  }

  async function onBatchDel(rows: any[]) {
    if (!rows.length) {
      message("请选择要删除的记录", { type: "warning" });
      return;
    }
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${rows.length} 条记录？`);
      const ids = rows.map((item: any) => item.id);
      const res = await del{{.ModelName}}({ ids });
      if (res.code === 200) {
        message("批量删除成功", { type: "success" });
        onSearch();
      } else {
        message(res.message || "删除失败", { type: "error" });
      }
    } catch {
      // cancelled
    }
  }

  return {
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
  };
}
