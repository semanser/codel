import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const messagesWrapper = style({
  display: "flex",
  flexDirection: "column",
  gap: 22,
});

export const titleStyles = style([
  font.textSmSemibold,
  {
    color: vars.color.gray11,
    textAlign: "center",
    marginBottom: 32,
  },
]);
