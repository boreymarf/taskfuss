import { defineConfig } from 'unocss'

// Docs: https://unocss.dev/guide/

export default defineConfig({
  theme: {
    colors: {
      /** Pure white in HSL. */
      white: 'hsl(0 0% 100%)',
      /** Pure red in HSL. */
      red: 'hsl(0 100% 50%)',
      /** Pure blue in HSL. */
      blue: 'hsl(240 100% 50%)',
    },
    spacing: {},
  },

  rules: [
    /** Padding scale: 1 unit = 0.25rem. */
    [/^p-(\d+(?:\.\d+)?)$/, ([, d]) => ({ padding: `${Number(d) / 4}rem` })],

    /** Margin scale: 1 unit = 0.25rem, supports negatives. */
    [/^m-(-?\d+(?:\.\d+)?)$/, ([, d]) => ({ margin: `${Number(d) / 4}rem` })],

    /** Background color resolved from theme.colors. */
    [/^bg-([\w-]+)$/, ([, c], { theme }) => {
      const color = theme.colors?.[c]
      if (color) return { 'background-color': color }
    }],
  ],

  shortcuts: {},
})