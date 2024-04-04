import { style } from "@vanilla-extract/css";

import { font } from "@/styles/font.css";
import { vars } from "@/styles/theme.css";

export const messagesWrapper = style({
  position: "relative",
  height: "100%",
});

export const messagesListWrapper = style({
  display: "flex",
  flexDirection: "column",
  // 100% - height of the new message textarea - height of the title bar
  maxHeight: "calc(100% - 100px - 90px)",
  overflowY: "scroll",
  gap: 22,
});

export const titleStyles = style([
  font.textSmSemibold,
  {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    gap: 12,
    color: vars.color.gray11,
    textAlign: "center",
    marginBottom: 16,
  },
]);

export const modelStyles = style({
  color: vars.color.gray10,
});

export const newMessageTextarea = style([
  font.textSmMedium,
  {
    position: "absolute",
    bottom: 0,
    height: 120,
    left: 0,
    backgroundColor: vars.color.gray4,
    border: `1px solid ${vars.color.gray5}`,
    borderRadius: "0 0 6px 6px",
    width: "calc(100% + 32px)",
    color: vars.color.gray12,
    padding: 16,
    margin: -16,
    boxShadow: `0 -20px 30px 10px ${vars.color.gray2}`,
    resize: "none",

    ":focus": {
      outline: "none",
      borderColor: vars.color.primary5,
    },

    ":disabled": {
      backgroundColor: vars.color.gray3,
      borderColor: vars.color.gray4,
    },
  },
]);
