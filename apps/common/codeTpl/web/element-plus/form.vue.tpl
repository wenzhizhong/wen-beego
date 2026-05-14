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
        <el-input type="textarea" v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" {{if .FormParam}}{{.FormParam}}{{end}} />
      </el-form-item>
      {{else if eq .FormType "datetime"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-date-picker v-model="form.{{.Name}}" type="datetime" placeholder="请选择{{.Comment}}" {{if .FormParam}}{{.FormParam}}{{end}} style="width: 100%" />
      </el-form-item>
      {{else if eq .FormType "switch"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-switch v-model="form.{{.Name}}" {{if .FormParam}}{{.FormParam}}{{end}} />
      </el-form-item>
      {{else if eq .FormType "select"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-select v-model="form.{{.Name}}" placeholder="请选择{{.Comment}}" style="width: 100%" {{if .FormParam}}{{.FormParam}}{{end}}>
          <el-option label="选项1" value="1" />
          <el-option label="选项2" value="2" />
        </el-select>
      </el-form-item>
      {{else if eq .FormType "radio"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-radio-group v-model="form.{{.Name}}" {{if .FormParam}}{{.FormParam}}{{end}}>
          <el-radio :value="1">选项1</el-radio>
          <el-radio :value="2">选项2</el-radio>
        </el-radio-group>
      </el-form-item>
      {{else if eq .FormType "checkbox"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-checkbox-group v-model="form.{{.Name}}" {{if .FormParam}}{{.FormParam}}{{end}}>
          <el-checkbox label="选项1" value="1" />
          <el-checkbox label="选项2" value="2" />
        </el-checkbox-group>
      </el-form-item>
      {{else if eq .FormType "imageUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-upload action="#" list-type="picture-card" :auto-upload="false" {{if .FormParam}}{{.FormParam}}{{end}}>
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
      {{else if eq .FormType "fileUpload"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-upload action="#" :auto-upload="false" {{if .FormParam}}{{.FormParam}}{{end}}>
          <el-button type="primary">上传文件</el-button>
        </el-upload>
      </el-form-item>
      {{else if eq .FormType "editor"}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <div class="editor-placeholder" {{if .FormParam}}{{.FormParam}}{{end}}>富文本编辑器（请自行引入组件）</div>
      </el-form-item>
      {{else}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-input v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" {{if .FormParam}}{{.FormParam}}{{end}} />
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
import { formRules } from "./utils/rule";
import type { FormInstance } from "element-plus";
import {upload} from "@/api/upload";

const emit = defineEmits(["refresh"]);
const visible = ref(false);
const dialogTitle = ref("新增");
const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const editId = ref("");

const form = reactive({
  {{range .Columns}}{{if not .SkipForm}}  {{.Name}}: "" as string | number,
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
    if (editId.value) {
      await edit{{.ModelName}}({ id: editId.value, ...form });
    } else {
      await add{{.ModelName}}(form);
    }
    message(editId.value ? "修改成功" : "新增成功", { type: "success" });
    emit("refresh");
    handleClose();
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
      form[key] = multiple ? [...form[key], res.data.filePath] : res.data.filePath;
    } else {
      message(res.message ||"上传失败", { type: "error" });
    }
  });
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
