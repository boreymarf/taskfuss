export const API_CONFIG = {
  VITE_SERVER_URL: import.meta.env.VITE_SERVER_URL || null
}

if (API_CONFIG.VITE_SERVER_URL === null) {
  throw new Error(`VITE_SERVER_URL is not set in the .env file!`);
}
