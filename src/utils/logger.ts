import * as tauri from "@tauri-apps/plugin-log";

type Opts = tauri.LogOptions | undefined;

/**
 * Extracts a caller location (`file` and `line`) for log messages. 
 * Location can be off by a few lines.
 *
 * - `skip` is the number of stack frames to skip before the caller we care about.
 * - Strips protocol, host, leading slash, and any `?query`.
 * - On Windows, converts `/` to `\` so JS log paths match Rust log paths.
 * - Returns `{}` if the stack frame cannot be parsed.
 *
 * Example:
 * ```ts
 * const { file, line } = loc(2);
 * // { file: "src\\main.ts", line: 8 }
 * ```
 */
function loc(skip: number): { file?: string; line?: number } {
  const raw = new Error().stack?.split("\n")[skip + 1];
  if (!raw) return {};
  const m = raw.match(/([^()\s]+):(\d+):\d+\)?\s*$/);
  if (!m) return {};
  let file = m[1]
    .replace(/^[a-z]+:\/\/[^/]*\//i, "")
    .replace(/^\//, "")
    .replace(/\?.*$/, "");
  if (navigator.userAgent.includes("Windows")) {
    file = file.replace(/\//g, "\\");
  }
  return { file, line: +m[2] };
}

export const debug = (msg: string, opts?: Opts) => tauri.debug(msg, { ...loc(2), ...opts });
export const info  = (msg: string, opts?: Opts) => tauri.info (msg, { ...loc(2), ...opts });
export const warn  = (msg: string, opts?: Opts) => tauri.warn (msg, { ...loc(2), ...opts });
export const error = (msg: string, opts?: Opts) => tauri.error(msg, { ...loc(2), ...opts });

export const logger = { debug, info, warn, error };