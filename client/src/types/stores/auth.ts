import type { FieldError } from "../api";


export class RegisterError extends Error {
  constructor(
    public message: string,
    public code: string,
    public details: FieldError[]
  ) {
    super(message);
    this.name = 'RegisterError';
  }
}
