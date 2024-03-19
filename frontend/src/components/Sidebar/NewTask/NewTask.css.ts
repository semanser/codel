import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const wrapperStyles = style([
  font.textSmSemibold,
  {
    background: "none",
    border: "none",
    textAlign: "left",
    color: vars.color.gray12,
    padding: "9px 16px",
    cursor: "pointer",
    marginBottom: 16,

    ":hover": {
      color: vars.color.primary9,
      backgroundColor: vars.color.gray2,
    },
  },
]);
