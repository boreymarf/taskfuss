import type { ApiResponse as GeneratedApiResponse } from "../../api/generated";

export interface ApiResponse<T = any> extends Omit<GeneratedApiResponse, 'data'> {
  data?: T;
}
