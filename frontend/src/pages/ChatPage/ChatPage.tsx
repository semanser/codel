import { Panel } from "@/components/Panel/Panel/Panel";

import { wrapperStyles } from "./ChatPage.css";

export const ChatPage = () => {
  return (
    <div className={wrapperStyles}>
      <Panel />
      <Panel />
    </div>
  );
};
