import { vars } from "@/styles/theme.css";
import { font } from "@/styles/font.css";
import { style } from "@vanilla-extract/css";

export const wrapperStyles = style([font.textSmSemibold, {
  padding: '9px 16px',

  ":hover": {
    color: vars.color.gray12,
    backgroundColor: vars.color.gray2,
  },
  selectors: {
    "&.active": {
      color: vars.color.gray12,
      backgroundColor: vars.color.gray2,
    }
  }
}]);

export const linkStyles = style({
  color: vars.color.gray11,
  textDecoration: "none",
});
