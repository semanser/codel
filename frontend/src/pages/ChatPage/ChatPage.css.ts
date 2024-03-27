import { globalStyle, style } from "@vanilla-extract/css";

export const wrapperStyles = style({
  display: "flex",
  flex: 1,
  padding: 16,
  gap: 16,
  maxWidth: 2000,
  margin: "0 auto",
});

export const tabsStyles = style({
  display: "flex",
  justifyContent: "space-between",
  width: "100%",
});

export const leftColumnStyles = style({
  display: "flex",
});

export const followButtonStyles = style({});

globalStyle(`${followButtonStyles} > svg`, {
  width: 20,
});
