import {
  defineConfig,
  presetIcons,
  presetUno,
  transformerVariantGroup,
} from "unocss";

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons(),
  ],
  shortcuts: {
    "h1-regular": "text-7 leading-6",
    "h1-medium": "text-7 font-medium leading-6",
    "h1-bold": "text-7 font-bold leading-6",
    "h2-regular": "text-6 leading-6",
    "h2-medium": "text-6 font-medium leading-6",
    "h2-bold": "text-6 font-bold leading-6",
    "h3-regular": "text-5 leading-6",
    "h3-medium": "text-5 font-medium leading-6",
    "h3-bold": "text-5 font-bold leading-6",
    "h4-regular": "text-4.5 leading-6",
    "h4-medium": "text-4.5 font-medium leading-6",
    "h4-bold": "text-4.5 font-bold leading-6",
    "text-regular": "text-4 leading-6",
    "text-medium": "text-4 font-medium leading-6",
    "text-bold": "text-4 font-bold leading-6",
    "caption-regular": "text-3.5 leading-6",
    "caption-medium": "text-3.5 font-medium leading-6",
    "caption-bold": "text-3.5 font-bold leading-6",
  },
  theme: {
    breakpoints: {
      "md": "768px",
    },
    colors: {
      primary: {
        white: "#FFFFFF",
        main: "#005BAC",
      },
      accent: {
        error: "#F25151",
        warn: "#F1B61E",
        success: "#20BD77",
      },
      transparent: {
        primaryHover: "rgba(0, 91, 172, 0.06)",
        primarySelected: "rgba(0, 91, 172, 0.10)",
        successHover: "rgba(32, 189, 119, 0.06)",
        successSelected: "rgba(32, 189, 119, 0.10)",
        warnHover: "rgba(241, 182, 30, 0.06)",
        warnSelected: "rgba(241, 182, 30, 0.10)",
        errorHover: "rgba(242, 81, 81, 0.06)",
        errorSelected: "rgba(242, 81, 81, 0.10)",
      },
      text: {
        black: "#2F3438",
        white: "#FFFFFF",
        grey: "#606A71",
        link: "#005BAC",
        disabled: "#B9BEC1",
      },
      ui: {
        border: "#CED6DB",
        background: "#F9F9F9",
        primary: "#FFFFFF",
        secondary: "#F0F2F5",
        tertiary: "#E2E5E9",
      },
      blackAlpha: {
        50: "rgba(0, 0, 0, 0.04)",
        100: "rgba(0, 0, 0, 0.06)",
        200: "rgba(0, 0, 0, 0.08)",
        300: "rgba(0, 0, 0, 0.16)",
        600: "rgba(0, 0, 0, 0.48)",
      },
      gray: {
        700: "#2D3748",
        800: "#1A202C",
        900: "#171923",
      },
      blue: {
        500: "#3182CE",
      },
    },
  },
  transformers: [transformerVariantGroup()],
});
