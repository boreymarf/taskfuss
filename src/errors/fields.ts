import type { FieldValueError } from "@/generated/bindings";
import type {
  SpectaErrorEntry,
  SpectaErrorPayloads,
  SpectaErrorTag,
} from "@/utils/specta-error";
import { extractSpectaError } from "@/utils/specta-error";

export type FieldValueErrorTag = SpectaErrorTag<FieldValueError>;
export type FieldValueErrorPayloads = SpectaErrorPayloads<FieldValueError>;
export type ExtractedFieldValueError = SpectaErrorEntry<FieldValueError>;

export const extractFieldValueError = (err: FieldValueError) =>
  extractSpectaError(err);