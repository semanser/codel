import { style } from "@vanilla-extract/css";

import { vars } from "@/styles/theme.css";

export const wrapperStyles = style({
  display: "flex",
  height: "100vh",
  backgroundColor: vars.color.gray1,
});
