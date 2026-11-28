import { createI18n } from "vue-i18n";
import type { ComposerTranslation } from "vue-i18n";
import type { FieldValueError } from "@/generated/bindings";
import { extractSpectaError } from "@/utils/errors";
import ru from "./locales/ru.json";
import en from "./locales/en.json";

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: "en",
  messages: { ru, en },
});

export default i18n;

export function formatFieldError(
  err: FieldValueError,
  t: ComposerTranslation,
): string {
  const entry = extractSpectaError(err);

  // Only Nested holds another error; format it recursively.
  if (entry.tag === "Nested") {
    return t("errors.fields.Nested", {
      message: formatFieldError(entry.payload, t),
    });
  }

  return t(
    `errors.fields.${entry.tag}`,
    (entry.payload ?? {}) as Record<string, unknown>,
  );
}