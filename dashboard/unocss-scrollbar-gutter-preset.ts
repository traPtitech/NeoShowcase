import { definePreset } from "unocss";

export const presetScrollbarGutter = definePreset({
  name: "unocss-preset-scrollbar-gutter",
  rules: [
    [
      /^scrollbar-gutter-both$/,
      () => ({
        "scrollbar-gutter": "stable both-edges",
      }),
    ],
    [
      /^scrollbar-gutter-(auto|stable|inherit|initial|revert|revert-layer|unset)$/,
      ([, p]) => ({
        "scrollbar-gutter": p,
      }),
    ],
  ],
});
