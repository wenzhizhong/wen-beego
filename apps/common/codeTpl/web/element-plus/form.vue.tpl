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
        <el-date-picker v-model="form.{{.Name}}" type="{{if eq .Type "date"}}{{.Type}}{{else}}datetime{{end}}" placeholder="请选择{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} format="{{if eq .Type "date"}}YYYY-MM-DD{{else}}YYYY-MM-DD HH:mm:ss{{end}}" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" @change="(dateStr: any)=>{form.{{.Name}} = dateParsingInZone(dateStr)}"/>
      </el-form-item>
      {{else if eq .FormType "time"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-time-picker v-model="form.{{.Name}}" placeholder="请选择{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} format="HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" @change="(dateStr: any)=>{form.{{.Name}} = dateParsingInZone(dateStr)}"/>
      </el-form-item>
      {{else if eq .FormType "switch"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-switch v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} />
      </el-form-item>
      {{else if eq .FormType "select"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-select v-model="form.{{.Name}}" placeholder="请选择{{.Comment}}" style="width: 100%" {{if .FormParamVue}}{{.FormParamVue}}{{end}}></el-select>
      </el-form-item>
      {{else if eq .FormType "radio"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
       {{/* <el-radio-group v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}}> */ -}}
        <el-radio-group v-model="form.{{.Name}}" >
          <template v-slot:default>
            <el-radio-button v-for="item in {{.Name}}Options" :key="item.id" :value="item.id">
              {{"{{"}} item.name {{"}}"}} 
            </el-radio-button>
          </template>
        </el-radio-group>
      </el-form-item>
      {{else if eq .FormType "checkbox"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        {{/* <el-checkbox-group v-model="form.{{.Name}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}}></el-checkbox-group> */ -}}
        <el-checkbox-group v-model="form.{{.Name}}">
          <el-checkbox-button v-for="item in {{.Name}}Options" :key="item.id" :value="item.id">
            {{"{{"}} item.name {{"}}"}} 
          </el-checkbox-button>
        </el-checkbox-group>
      </el-form-item>
      {{else if eq .FormType "imageUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <sliceUploadV2
          :bucketACL="'public'"
          :uploadSuccessCallback="(uploadResult: any, file: any)=>{uploadFileSuccessCallback(uploadResult, file, '{{.Name}}')}"
          :elementPlusUploader="{ 'list-type': 'picture-card', 'show-file-list': true, 'file-list': fileListObj['{{.Name}}'], 'on-remove': (e: any)=>{removeFile(e, '{{.Name}}')} {{if .FormParamVue}} , {{.FormParamVue}}{{end}} }">
          <template #default>
            <AddLargeLine style="width: 40px; height: 40px;"/>
          </template>
        </sliceUploadV2>
      </el-form-item>
      {{else if eq .FormType "fileUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <sliceUploadV2
          :bucketACL="'public'"
          :uploadSuccessCallback="(uploadResult: any, file: any)=>{uploadFileSuccessCallback(uploadResult, file, '{{.Name}}')}"
          :elementPlusUploader="{ 'list-type': 'text', 'show-file-list': true, 'file-list': fileListObj['{{.Name}}'], 'on-remove': (e: any)=>{removeFile(e, '{{.Name}}')} {{if .FormParamVue}} , {{.FormParamVue}}{{end}} }">
          <template #default>
            <el-button type="primary">上传文件</el-button>
          </template>
        </sliceUploadV2>
      </el-form-item>
      {{else if eq .FormType "editor"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <EditorBase v-if="editorVisible" v-model="form.{{.Name}}"/>
      </el-form-item>
      {{else}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-input v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" {{if .FormParamVue}}{{.FormParamVue}}{{end}} />
      </el-form-item>
      {{end}}
      {{end}}{{end}}    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit" :disabled="!hasPerms(['{{.AuthAdd}}', '{{.AuthEdit}}'])">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick } from "vue";
import type { FormInstance } from "element-plus";
import { ElMessageBox } from "element-plus";
import { hasPerms } from "@/utils/auth";
import { message } from "@/utils/message";
import { dateParsingInZone } from "@/utils/time";
import { getVarType } from "@/components/sliceUpload/common";
import sliceUploadV2 from "@/components/sliceUploadV2/index.vue";
import { EditorBase } from "@/components/Editor/components";
import AddLargeLine from "~icons/ri/add-large-line";
import { get{{.ModelName}}Detail, add{{.ModelName}}, edit{{.ModelName}} } from "@/api/{{.ModelNameLower}}";
import { formRules } from "./utils/rule";
import { getLinkById } from "@/api/upload";
{{range .Columns}}{{if .FormParamTs}}
{{if or (eq .FormType "select") (eq .FormType "radio") (eq .FormType "checkbox")}}import { {{.Name}}Options, {{.Name}}Props } from "./utils/hook";
{{else if eq .FormType "switch"}}import { {{.Name}}Options } from "./utils/hook";
{{end}}{{end}}{{end}}

const emit = defineEmits(["refresh"]);
const visible = ref(false);
const dialogTitle = ref("新增");
const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const editId = ref("");
const editorVisible = ref(false);

const fileListObj = ref<Record<string, any[]>>({
{{range .Columns}}{{if eq .FormType "imageUpload"}}  {{.Name}}: [],
{{else if eq .FormType "fileUpload"}}  {{.Name}}: [],
{{end}}{{end}}
});

const needStringifyFields = reactive([
{{range .Columns}}{{if eq .FormType "checkbox"}}  "{{.Name}}",
{{else if eq .FormType "select"}}  "{{.Name}}",
{{end}}{{end}}  ...(Object.keys(fileListObj.value))
]);
const needConvertBoolFields = reactive([
{{range .Columns}}{{if eq .FormType "switch"}}  "{{.Name}}",
{{end}}{{end}}
]);

const form = reactive({
{{range .Columns}}{{if not .SkipForm}}
    {{- if isMultipleCompt .FormType -}}
      {{- if and (eq .FormType "select") (not (hasMultipleProp .FormParam)) -}}
        {{.Name}}: {{if .DefVal}}{{.DefVal}}{{else}}""{{end}},
      {{- else -}}
        {{.Name}}: [],
      {{- end -}}
    {{- else if eq .FormType "switch" -}}
      {{.Name}}: {{if .DefVal}}{{.DefVal}}{{else}}0{{end}},
    {{- else -}}
      {{.Name}}: {{if .DefVal}}{{.DefVal}}{{else}}""{{end}},
    {{- end}}
{{end}}{{end}}
});

function open(title = "新增", row?: any) {
  dialogTitle.value = title;
  visible.value = true;
  editorVisible.value = false;

  nextTick(async () => {
    if (row && row.id) {
      editId.value = row.id;
      const res = await get{{.ModelName}}Detail({ id: row.id });
      const data = res.data ?? {};
      for (let i = 0; i < needStringifyFields.length; i++) {
        if (data[needStringifyFields[i]] === null || data[needStringifyFields[i]] === undefined) continue;
        if (getVarType(form[needStringifyFields[i]]) === 'array') {
          data[needStringifyFields[i]] = data[needStringifyFields[i]].split(",");
        }
      }
      for (let i = 0; i < needConvertBoolFields.length; i++) {
        if (data[needConvertBoolFields[i]] === null || data[needConvertBoolFields[i]] === undefined) continue;
        data[needConvertBoolFields[i]] = Boolean(data[needConvertBoolFields[i]]);
      }
      Object.assign(form, data);
{{range .Columns}}{{if or (eq .FormType "imageUpload") (eq .FormType "fileUpload")}}      fileListObj.value['{{.Name}}'] = await doGetLinkById("{{.Name}}");
{{end}}{{end}}    } else {
      editId.value = "";
      formRef.value?.resetFields();
    }

    nextTick(() => {
      editorVisible.value = true;
    });
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
    let tmpForm: any = JSON.parse(JSON.stringify(form));
    for (let i = 0; i < needStringifyFields.length; i++) {
      if (form[needStringifyFields[i]] === null || form[needStringifyFields[i]] === undefined) continue;
      if (getVarType(form[needStringifyFields[i]]) === 'array') {
        tmpForm[needStringifyFields[i]] = form[needStringifyFields[i]].join(",");
      }
    }
    for (let i = 0; i < needConvertBoolFields.length; i++) {
      if (form[needConvertBoolFields[i]] === null || form[needConvertBoolFields[i]] === undefined) continue;
      tmpForm[needConvertBoolFields[i]] = Number(!!tmpForm[needConvertBoolFields[i]]);
    }
    let result;
    if (editId.value) {
      result = await edit{{.ModelName}}({ id: editId.value, ...tmpForm });
    } else {
      result = await add{{.ModelName}}(tmpForm);
    }
    if (result && result.code === 200) {
      message(editId.value ? "修改成功" : "新增成功", { type: "success" });
    } else {
      message(result?.message || "操作失败", { type: "error" });
    }
    emit("refresh");
    handleClose();
  } catch (error: any) {
    message(error?.message || "操作失败", { type: "error" });
  } finally {
    submitLoading.value = false;
  }
}

const uploadFileSuccessCallback = async (uploadResult: any, file: any, key: string) => {
  form[key] = form[key].length ? [...form[key], uploadResult.fileId] : [uploadResult.fileId];
  fileListObj.value[key] = await doGetLinkById(key);
};

function removeFile(e: any, key: string) {
  const fileId = e?.raw && e.raw instanceof File ? e.raw.fileId : e.id;
  if (!fileId) return;
  fileListObj.value[key] = fileListObj.value[key].filter((item: any) => item.id != fileId);
  form[key] = form[key].filter((item: any) => item != fileId);
}

const getShowUploadListItem = (fileId: string, fileName: string, filePath: string) => {
  return { id: fileId, name: fileName, url: filePath };
};

const doGetLinkById = async (key: string) => {
  if (!form[key] || !form[key].length) return [];
  return await getLinkById(form[key]).then((res: any) => {
    if (res?.code !== 200) {
      console.error(res?.message || "getLinkById() 获取数据错误");
      return [];
    }
    return (res?.data?.list || []).map((item: any) =>
      getShowUploadListItem(item.id, item.realName, item.path)
    );
  });
};

defineExpose({ open });
</script>
