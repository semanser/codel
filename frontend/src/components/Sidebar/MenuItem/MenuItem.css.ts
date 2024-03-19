import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const wrapperStyles = style([
  font.textSmSemibold,
  {
    borderRadius: 6,

    ":hover": {
      color: vars.color.gray12,
      backgroundColor: vars.color.gray2,
    },
    selectors: {
      "&.active": {
        color: vars.color.gray12,
        backgroundColor: vars.color.gray2,
      },
    },
  },
]);

export const linkStyles = style({
  color: vars.color.gray11,
  textDecoration: "none",
  display: "block",
  padding: "9px 16px",
});
