import { reactive, ref, onMounted } from "vue";
import { get{{.ModelName}}List, del{{.ModelName}}, add{{.ModelName}}, edit{{.ModelName}}, get{{.ModelName}}Detail } from "./api";
import { message, messageBox } from "@/utils/message";

export function use{{.ModelName}}() {
  const loading = ref(false);
  const dataList = ref([]);
  const selectedRows = ref([]);

  const searchForm = reactive({
    keyword: ""
  });

  const pagination = reactive({
    currentPage: 1,
    pageSize: 10,
    total: 0
  });

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
      selectedRows.value = [];
    } finally {
      loading.value = false;
    }
  }

  function onReset() {
    searchForm.keyword = "";
    onSearch();
  }

  function onPageChange({ currentPage, pageSize }: any) {
    pagination.currentPage = currentPage;
    pagination.pageSize = pageSize;
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

  return {
    loading,
    dataList,
    selectedRows,
    searchForm,
    pagination,
    onSearch,
    onReset,
    onPageChange,
    handleDelete,
    handleBatchDelete
  };
}
