export interface ApiResponse<T> {
  data: T;
  timestamp: number;
  latency: string;
}

export interface ApiError {
  code: string,
  message: string,
  details: any
  timestamp: number
}
