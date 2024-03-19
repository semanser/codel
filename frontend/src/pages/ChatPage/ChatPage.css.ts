import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const wrapperStyles = style({
  display: "flex",
  flex: 1,
  padding: 16,
  gap: 16,
  maxWidth: 2000,
  margin: "0 auto",
});

export const titleStyles = style([
  font.textSmSemibold,
  {
    color: vars.color.gray11,
    textAlign: "center",
  },
]);
