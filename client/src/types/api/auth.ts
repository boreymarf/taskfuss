export interface User {
  id: number;
  username: string;
  created_at: Date;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  user: User;
  auth_token: string;
  expires_at: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  auth_token: string;
  expires_at: number;
}

export interface ValidationError {
  code: string;
  message: string;
  details: FieldError[];
}

export interface FieldError {
  field: string;
  code: string;
  message: string;
}
