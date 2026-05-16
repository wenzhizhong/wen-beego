<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    :close-on-click-modal="false"
    width="600px"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
      {{range .Columns}}{{if not .SkipForm}}
      {{if eq .FormType "textarea"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-input type="textarea" v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} />
      </el-form-item>
      {{else if eq .FormType "datetime"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-date-picker v-model="form.{{.Name}}" type="datetime" placeholder="请选择{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} style="width: 100%" />
      </el-form-item>
      {{else if eq .FormType "switch"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-switch v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} />
      </el-form-item>
      {{else if eq .FormType "select"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-select v-model="form.{{.Name}}" placeholder="请选择{{.Comment}}" style="width: 100%" {{if .FormParamVue}}{{.FormParamVue}}{{end}}>
          <el-option label="选项1" value="1" />
          <el-option label="选项2" value="2" />
        </el-select>
      </el-form-item>
      {{else if eq .FormType "radio"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-radio-group v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}}>
        </el-radio-group>
      </el-form-item>
      {{else if eq .FormType "checkbox"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-checkbox-group v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}}>
        </el-checkbox-group>
      </el-form-item>
      {{else if eq .FormType "imageUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-upload action="#" list-type="picture-card" :auto-upload="false" {{if .FormParamVue}}{{.FormParamVue}}{{end}}  @change="e=>{uploadFile(e, '{{.Name}}', {{if isMultipleCompt .FormParam}}true{{else}}false{{end}} )}"  :on-remove="(e)=>{removeFile(e, '{{.Name}}')}"  :file-list="fileListObj['{{.Name}}']">
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
      {{else if eq .FormType "fileUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-upload action="#" :auto-upload="false" {{if .FormParamVue}}{{.FormParamVue}}{{end}}  @change="e=>{uploadFile(e, '{{.Name}}', {{if isMultipleCompt .FormParam}}true{{else}}false{{end}} )}"  :on-remove="(e)=>{removeFile(e, '{{.Name}}')}"  :file-list="fileListObj['{{.Name}}']">
          <el-button type="primary">上传文件</el-button>
        </el-upload>
      </el-form-item>
      {{else if eq .FormType "editor"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <div class="editor-placeholder" {{if .FormParamVue}}{{.FormParamVue}}{{end}}>富文本编辑器（请自行引入组件）</div>
      </el-form-item>
      {{else}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-input v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} />
      </el-form-item>
      {{end}}
      {{end}}{{end}}    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick } from "vue";
import { get{{.ModelName}}Detail, add{{.ModelName}}, edit{{.ModelName}} } from "@/api/{{.ModelNameLower}}";
import { message } from "@/utils/message";
import { ElMessageBox } from "element-plus";
import { formRules } from "./utils/rule";
import type { FormInstance } from "element-plus";
import {upload} from "@/api/upload";
import { getVarType } from "@/components/sliceUpload/common";

{{range .Columns}}{{if .FormParamTs}}
{{.FormParamTs}}
{{end}}{{end}}

const emit = defineEmits(["refresh"]);
const visible = ref(false);
const dialogTitle = ref("新增");
const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const editId = ref("");

const fileListObj = ref({
  {{range .Columns}}{{if eq .FormType "imageUpload"}} {{.Name}}: [],
  {{else if eq .FormType "fileUpload"}} {{.Name}}: [],
  {{end}}{{end}}
});
const needStringifyFields = reactive([
  {{range .Columns}}{{if eq .FormType "checkbox"}}  "{{.Name}}",
  {{else if eq .FormType "select"}}  "{{.Name}}",
  {{end}}{{end}}
  ...(Object.keys(fileListObj.value))
]);
const needConvertBoolFields = reactive([
  {{range .Columns}}{{if eq .FormType "switch"}}  "{{.Name}}",
  {{end}}{{end}}
]);

const form = reactive({
  {{range .Columns}}{{if not .SkipForm}}
    {{- if isMultipleComptFields .FormType -}}
      {{- if and (eq .FormType "select") (isMultipleCompt .FormParam) -}}
        {{.Name}}: {{if and (eq .TsType "string") (eq .DefVal "")}}{{.DefVal}}{{else}}{{.DefVal}}{{end}},
      {{- else -}}
        {{.Name}}: [],
      {{- end -}}
    {{- else -}}
      {{.Name}}: {{if and (eq .TsType "string") (eq .DefVal "")}}""{{else}}{{.DefVal}}{{end}},
    {{- end}}
  {{end}}{{end}}
});

function open(title = "新增", row?: any) {
  dialogTitle.value = title;
  visible.value = true;
  nextTick(async () => {
    if (row && row.id) {
      editId.value = row.id;
      const res = await get{{.ModelName}}Detail({ id: row.id });
      const data = res.data ?? {};
      Object.assign(form, data);
    } else {
      editId.value = "";
      formRef.value?.resetFields();
    }
  });
}

function handleClose() {
  formRef.value?.resetFields();
  editId.value = "";
  visible.value = false;
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;
  submitLoading.value = true;
  try {
    await ElMessageBox.confirm(editId.value ? "是否修改?" : "是否新增?");

    let tmpForm = JSON.parse(JSON.stringify(form));
    for (let i = 0; i < needStringifyFields.length; i++) {
      if (!form[needStringifyFields[i]] || getVarType(form[needStringifyFields[i]]) == 'array') {
        console.log("needStringifyFields", needStringifyFields[i], getVarType(form[needStringifyFields[i]]));
        tmpForm[needStringifyFields[i]] = form[needStringifyFields[i]].join(",");
      }
    }
    for (let i = 0; i < needConvertBoolFields.length; i++) {
      tmpForm[needConvertBoolFields[i]] = Number(!!(tmpForm[needConvertBoolFields[i]]))
    }
    
    let result
    if (editId.value) {
      result =  await editGenerateForm({ id: editId.value, ...tmpForm })
    } else {
      result =  await addGenerateForm(tmpForm)
    }
    let msg = ""
    if  (result && result.code === 200) {
      msg = editId.value ? "修改成功" : "新增成功"
      message(msg, { type: "success" });
    } else {
      msg = result.message || "操作失败"
      message(msg, { type: "error" });
    }
    emit("refresh");
    handleClose();
  }catch(error) {
    let errMsg = typeof error === "string" ? error : error.message;
    message(errMsg || "修改失败", { type: "error" });
  } finally {
    submitLoading.value = false;
  }
}

/**
* 上传文件
* file 文件对象
* key 字段名
* multiple 是否多文件上传
*/
async function uploadFile(file: any, key: string, multiple: boolean = false) { 
  if (!file) return;
  upload(file.raw, file.name, file.size).then(async (res) => { 
    if (res.code === 200 && res?.data?.filePath) {
      message("上传成功", { type: "success" });

      let fileItem = {id: res.data.fileId, name: file.name, url: res.data.filePath}
      form[key] = multiple ? [...form[key], res.data.fileId] : [res.data.fileId];
      multiple? fileListObj.value[key].push(fileItem) : fileListObj.value[key] = [fileItem];
    } else {
      message(res.message ||"上传失败", { type: "error" });
    }
  });
}
function removeFile(e :any, key:string){
  if (fileListObj.value[key]){
    let fileId = e?.raw && e.raw instanceof File ? e.raw.fileId : e.id;
    if (!fileId) {
      message("请选择文件", { type: "error" });
      return
    }
    
    fileListObj.value[key] = fileListObj.value[key].filter(item=>item.id != fileId);
    form[key] = form[key].filter(item=>item != fileId);
  }
}

defineExpose({ open });
</script>

<style scoped>
.editor-placeholder {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 8px;
  color: var(--el-text-color-placeholder);
}
</style>
