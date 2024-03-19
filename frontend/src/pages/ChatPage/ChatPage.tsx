import { Message, MessageType } from "@/components/Message/Message/Message";
import { Panel } from "@/components/Panel/Panel/Panel";

import { messagesWrapper, titleStyles, wrapperStyles } from "./ChatPage.css";

const fakeData = {
  title: "This is a chat",
  messages: [
    {
      id: 1,
      message: "This is a test message",
      time: new Date("2024-01-10"),
      type: MessageType.Browser,
    },
    {
      id: 2,
      message:
        "This is a some pretty long message. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum placerat felis ante, non semper mi hendrerit id. Praesent sodales est dui, ut semper sem consectetur nec. Praesent vitae euismod metus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Maecenas vitae ante interdum erat blandit eleifend.",
      time: new Date("2024-03-18"),
      type: MessageType.Terminal,
    },
    {
      id: 3,
      message: "This is some random message",
      time: new Date("2024-03-18"),
      type: MessageType.Code,
    },
    {
      id: 4,
      message: "This is some ask message",
      time: new Date("2024-03-18"),
      type: MessageType.Ask,
    },
    {
      id: 5,
      message: "This task is done",
      time: new Date("2024-03-18"),
      type: MessageType.Done,
    },
  ],
};

export const ChatPage = () => {
  return (
    <div className={wrapperStyles}>
      <Panel>
        <div className={titleStyles}>{fakeData.title}</div>
        <div className={messagesWrapper}>
          {fakeData.messages.map((message) => (
            <Message key={message.id} {...message} />
          ))}
        </div>
      </Panel>
      <Panel>test</Panel>
    </div>
  );
};
