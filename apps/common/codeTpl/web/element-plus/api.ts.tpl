import { http } from "@/utils/http";
import type { Result, ResultTable } from "@/utils/http/types.api";

const BASE_API = "/api{{.ApiUrlPrefix}}";

export const get{{.ModelName}}List = (params?: object) => {
  return http.request<ResultTable>("get", BASE_API + "/get", { params });
};

export const get{{.ModelName}}Detail = (params?: object) => {
  return http.request<Result>("get", BASE_API + "/detail", { params });
};

export const add{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_API + "/add", { data });
};

export const edit{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_API + "/edit", { data });
};

export const del{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_API + "/del", { data });
};
