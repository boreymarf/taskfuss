
// TODO: Please review this AI bullshit

type RequiredKeys<T> = {
  [K in keyof T]-?: {} extends Pick<T, K> ? never : K;
}[keyof T] & string;

export type SpectaErrorTag<T> = T extends string ? T : RequiredKeys<T>;

type PayloadOf<T, Tag extends string> =
  | (T extends string ? (T extends Tag ? undefined : never) : never)
  | (T extends object
      ? Tag extends keyof T
        ? NonNullable<T[Tag]>
        : never
      : never);

export type SpectaErrorPayloads<T> = {
  [K in SpectaErrorTag<T>]: PayloadOf<T, K>;
};

export type SpectaErrorEntry<T extends string | object> = {
  [K in SpectaErrorTag<T>]: { tag: K; payload: SpectaErrorPayloads<T>[K] };
}[SpectaErrorTag<T>];

export function extractSpectaError<T extends string | object>(
  err: T,
): SpectaErrorEntry<T> {
  if (typeof err === "string") {
    return { tag: err, payload: undefined } as SpectaErrorEntry<T>;
  }

  const entries = Object.entries(err as Record<string, unknown>);
  for (const [key, value] of entries) {
    if (value !== undefined) {
      return { tag: key, payload: value } as SpectaErrorEntry<T>;
    }
  }

  // Инвариант `?: never` гарантирует, что сюда не попадём
  throw new Error(`Invalid specta error: no active tag in ${JSON.stringify(err)}`);
}