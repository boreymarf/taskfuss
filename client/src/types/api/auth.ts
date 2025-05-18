export interface User {
  id: number;
  username: string;
  createdAt: Date;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  user: User;
  authToken: string;
  expiresAt: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  authToken: string;
  expiresAt: number;
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
