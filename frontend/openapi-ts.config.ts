// openapi-ts.config.ts
import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../shared/temp/openapi.json',
  output: {
    path: 'src/api/generated',
    postProcess: ['prettier'],
  },
  plugins: [
    '@hey-api/typescript',
    '@hey-api/client-fetch',
    {
      name: '@tanstack/vue-query',
    },
  ],
})
