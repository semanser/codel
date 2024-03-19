import { Panel } from "@/components/Panel/Panel/Panel";

import { titleStyles, wrapperStyles } from "./ChatPage.css";

const fakeData = {
  title: "This is a chat",
};

export const ChatPage = () => {
  return (
    <div className={wrapperStyles}>
      <Panel>
        <div className={titleStyles}>{fakeData.title}</div>
      </Panel>
      <Panel>test</Panel>
    </div>
  );
};
