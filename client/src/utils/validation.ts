export function isValidEmail(email: string): Boolean {
  if (!email.includes('@')) return false;
  return true;
}

export function isValidPassword(password: string): Boolean {
  if (password.length < 6) return false
  return true
}
