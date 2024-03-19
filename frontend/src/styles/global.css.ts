import { globalStyle } from "@vanilla-extract/css";

import { vars } from "./theme.css";

globalStyle("html, body", {
  margin: 0,
  padding: 0,
  color: vars.color.gray12,
});

globalStyle("*", {
  fontFamily: "Inter var, sans-serif",
  WebkitFontSmoothing: "antialiased",
  boxSizing: "border-box",
});
