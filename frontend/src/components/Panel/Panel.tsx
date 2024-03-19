import { wrapperStyles } from "./Panel.css";

type PanelProps = {
  children: React.ReactNode;
};

export const Panel = ({ children }: PanelProps) => (
  <div className={wrapperStyles}>{children}</div>
);
