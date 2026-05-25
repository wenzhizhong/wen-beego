import { reactive, ref, onMounted, nextTick } from "vue";
import dayjs from "dayjs";
import { message } from "@/utils/message";
import { ElMessageBox, ElDialog, ElImage, ElLink } from "element-plus";
import { listToMap, getFileSuffix, getFileIcon } from "@/utils/util.js";
import { dateParsingInZone } from "@/utils/time";
import { get{{.ModelName}}List, del{{.ModelName}} } from "@/api/{{.ModelNameLower}}";
import { getLinkById } from "@/api/upload";
import type { PaginationProps } from "@pureadmin/table";
import type { FormItemProps } from "./types";

{{range .Columns}}{{if .FormParamTs}}
{{.FormParamTs}}
{{end}}{{end}}

{{range .Columns}}{{if .FormParamTs}}
const {{.Name}}OptionsMap = listToMap({{.Name}}Options, 'id', 'name');
{{end}}{{end}}

export const previewVisible = ref(false);
export const previewFileList = ref<any[]>([]);

export async function handlePreview(row: any, key: string) {
  if (!row[key]) return;
  const ids = typeof row[key] === 'string' ? row[key].split(',').filter(Boolean) : [row[key]];
  if (!ids.length) return;
  const res: any = await getLinkById(ids);
  if (res?.code !== 200) { message("获取文件链接失败", { type: "error" }); return; }
  previewFileList.value = (res?.data?.list || []).map((item: any) => ({
    id: item.id, name: item.realName, path: item.path, suffix: getFileSuffix(item.realName)
  }));
  previewVisible.value = true;
}

export function use{{.ModelName}}(formRef: any) {
  const loading = ref(false);
  const dataList = ref([]);
  const selectedNum = ref(0);

  const searchForm = reactive({
{{range .Columns}}{{if not (eq .Name "id")}}{{if not (isHasDeletedFields .Name)}}{{if not (eq .FormType "editor")}}{{if not (eq .FormType "imageUpload")}}{{if not (eq .FormType "fileUpload")}}    {{.Name}}: {{if isMultipleCompt .FormType}}[]{{else}}""{{end}},
{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}  
    {{if .HasUnitId}}selectUnitIds:"",{{end}} });

  const pagination = reactive<PaginationProps>({
    total: 0, pageSize: 10, currentPage: 1, background: true
  });

  const columns: TableColumnList = [
    { label: "勾选列", type: "selection", fixed: "left", reserveSelection: true },
{{range .Columns}}{{if not (eq .Name "id") }}{{if not (eq .Name "unit_id") }}{{if not (isHasDeletedFields .Name)}}{{if not (isDeleteUserIdFields .Name)}}{{if not (isCreateUserIdFields .Name)}}{{if not (isUpdateUserIdFields .Name)}}{{if not (isDeletedTimeFields .Name)}}{{if not (eq .FormType "editor")}}    {
      label: "{{.Comment}}",
      prop: "{{.Name}}",
{{if .FormParamTs}}{{if or (eq .FormType "switch") (eq .FormType "radio")}}      cellRenderer: ({ row, props }: any) => <span>{ {{.Name}}OptionsMap[row.{{.Name}}] ?? row.{{.Name}} }</span>,
{{else if or (eq .FormType "select") (eq .FormType "checkbox")}}      cellRenderer: ({ row, props }: any) => {
        if (!row.{{.Name}}) return <span>-</span>;
        return <span>{ (row.{{.Name}} || '').split(',').filter(Boolean).map((v: string) => {{.Name}}OptionsMap[v] ?? v).join(', ') }</span>;
      },
{{end}}{{else if or (eq .FormType "imageUpload") (eq .FormType "fileUpload")}}      cellRenderer: ({ row, props }: any) => {
        if (!row.{{.Name}}) return <span>-</span>;
        const count = (row.{{.Name}} || '').split(',').filter(Boolean).length;
        return <el-link type="primary" onClick={() => handlePreview(row, '{{.Name}}')}>{ count } 个文件</el-link>;
      },
{{else if or (eq .Type "timestamp without time zone") (eq .Type "timestamp with time zone") (eq .Type "timestamptz") (eq .Type "timestamp")}}      formatter: ({ {{.Name}} }: any) => {{.Name}} ? dayjs({{.Name}}).format("YYYY-MM-DD HH:mm:ss") : "",
{{else if eq .Type "date"}}      formatter: ({ {{.Name}} }: any) => {{.Name}} ? dayjs({{.Name}}).format("YYYY-MM-DD") : "",
{{else if eq .FormType "datetime"}}      formatter: ({ {{.Name}} }: any) => {{.Name}} ? dayjs({{.Name}}).format("YYYY-MM-DD HH:mm:ss") : "",
{{else if eq .FormType "time"}}      formatter: ({ {{.Name}} }: any) => {{.Name}} ? dayjs({{.Name}}).format("HH:mm:ss") : "",
{{end}}      minWidth: 120
    },
{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}
{{if .HasCreateUserId}}    { label: "创建人", prop: "created_by_name", minWidth: 100 },
{{end}}{{if .HasUpdateUserId}}    { label: "更新人", prop: "updated_by_name", minWidth: 100 },
{{end}}{{if .HasUnitId}}    { label: "组织单位", prop: "unit_name", minWidth: 120 },
{{end}}    { label: "操作", fixed: "right", width: 180, slot: "operation" }
  ];

  onMounted(() => { onSearch(); });

  async function onSearch() {
    loading.value = true;
    try {
      const searchDto: any = {};
      for (const k of Object.keys(searchForm)) {
        let val = (searchForm as any)[k];
        if (val === '' || val === null || val === undefined) continue;
        if (Array.isArray(val) && val.length === 2 && !isNaN(Date.parse(val[0]))) {
          searchDto[k + 'Start'] = dateParsingInZone(val[0]);
          searchDto[k + 'End']   = dateParsingInZone(val[1]);
        } else if (Array.isArray(val)) {
          searchDto[k] = val.join(',');
        } else {
          searchDto[k] = val;
        }
      }
      const res = await get{{.ModelName}}List({
        currentPage: pagination.currentPage,
        pageSize: pagination.pageSize,
        dto: encodeURIComponent(JSON.stringify(searchDto))
      });
      const { list, total } = res.data ?? { list: [], total: 0 };
      dataList.value = list ?? [];
      pagination.total = total ?? 0;
      selectedNum.value = 0;
    } finally { loading.value = false; }
  }

  function resetForm(formEl: any) { if (formEl) { formEl.resetFields(); onSearch(); } }
  function handleSizeChange(val: number) { pagination.pageSize = val; onSearch(); }
  function handleCurrentChange(val: number) { pagination.currentPage = val; onSearch(); }
  function handleSelectionChange(val: any) { selectedNum.value = val.length; }
  function onSelectionCancel() { selectedNum.value = 0; }

  function openDialog(title = "新增", row?: FormItemProps) {
    nextTick(() => { formRef.value?.open(title, row); });
  }

  async function handleDelete(row: any) {
    try {
      await ElMessageBox.confirm("确认删除该记录？");
      const res = await del{{.ModelName}}({ id: row.id });
      if (res.code === 200) { message("删除成功", { type: "success" }); onSearch(); }
      else { message(res.message || "删除失败", { type: "error" }); }
    } catch {}
  }

  async function onBatchDel(rows: any[]) {
    if (!rows.length) { message("请选择要删除的记录", { type: "warning" }); return; }
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${rows.length} 条记录？`);
      const ids = rows.map((item: any) => item.id);
      const res = await del{{.ModelName}}({ ids });
      if (res.code === 200) { message("批量删除成功", { type: "success" }); onSearch(); }
      else { message(res.message || "删除失败", { type: "error" }); }
    } catch {}
  }
  {{if .HasUnitId}}function onTreeSelect({ id, selected}){
    searchForm.selectUnitIds = selected ? id : "";
    onSearch();
  }{{end}} 

  return {
    loading, dataList, selectedNum, searchForm, pagination, columns,
    onSearch, resetForm, handleSizeChange, handleCurrentChange,
    handleSelectionChange, onSelectionCancel, openDialog, handleDelete, onBatchDel, 
    {{if .HasUnitId}}onTreeSelect,{{end}} 
  };
}
