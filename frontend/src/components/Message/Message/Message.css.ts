import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const wrapperStyles = style({
  display: "flex",
  gap: 12,
});

export const avatarStyles = style({
  borderRadius: "50%",
});

export const rightColumnStyles = style({
  display: "flex",
  flexDirection: "column",
  gap: 6,
});

export const timeStyles = style([
  font.textXsRegular,
  {
    color: vars.color.primary8,
  },
]);

export const messageStyles = style([
  font.textSmRegular,
  {
    color: vars.color.primary12,
    padding: "10px 14px",
    background: vars.color.gray1,
    borderRadius: "0 8px 8px 8px",
    border: `1px solid ${vars.color.gray3}`,
  },
]);
