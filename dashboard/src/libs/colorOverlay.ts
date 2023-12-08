export const colorOverlay = (baseColor: string, overlayColor: string) =>
  `linear-gradient(0deg, ${overlayColor} 0%, ${overlayColor} 100%), ${baseColor}`
