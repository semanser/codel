import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const headerStyles = style([
  font.textXsSemibold,
  {
    backgroundColor: vars.color.gray6,
    color: vars.color.gray11,
    padding: "8px 12px",
    borderRadius: "8px 8px 0 0",
    display: "flex",
    alignItems: "center",
    gap: 8,
  },
]);

export const wrapperStyles = style({
  backgroundColor: vars.color.gray2,
  borderRadius: 8,
  border: `1px solid ${vars.color.gray3}`,
  overflow: "hidden",
});

export const imgStyles = style({
  width: "100%",
});

export const imgWrapperStyles = style({
  backgroundColor: vars.color.gray12,
  width: "100%",
  minHeight: "auto",
});
