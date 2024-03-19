import { style } from "@vanilla-extract/css";

export const wrapperStyles = style({
  width: "250px",
  display: "flex",
  flexDirection: "column",
  gap: 8,
  padding: 16,
  overflowY: "scroll",
});
