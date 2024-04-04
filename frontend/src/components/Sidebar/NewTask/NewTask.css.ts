import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const wrapperStyles = style({
  display: "flex",
  alignItems: "center",
  marginBottom: 16,
  justifyContent: "space-between",
});

export const linkWrapperStyles = style([
  font.textSmSemibold,
  {
    display: "block",
    textDecoration: "none",
    background: vars.color.gray3,
    border: "none",
    textAlign: "left",
    color: vars.color.gray12,
    padding: "9px 16px",
    cursor: "pointer",
    borderRadius: "6px 0 0 6px",
    flex: 1,

    selectors: {
      "&.active": {
        color: vars.color.primary9,
        backgroundColor: vars.color.gray5,
      },
    },

    ":hover": {
      color: vars.color.primary9,
      backgroundColor: vars.color.gray4,
    },
  },
]);
