import { Task } from "@/generated/graphql";

import { Message, MessageType } from "./Message/Message";
import {
  messagesListWrapper,
  messagesWrapper,
  newMessageTextarea,
  titleStyles,
} from "./Messages.css";

type MessagesProps = {
  tasks: Task[];
  name: string;
};

export const Messages = ({ tasks, name }: MessagesProps) => {
  const messages =
    tasks.map((task) => ({
      id: task.id,
      message: task.message,
      time: task.createdAt,
      status: task.status,
      // TODO Add the correct type and output
      type: MessageType.Terminal,
      output: "Test output",
    })) ?? [];

  return (
    <div className={messagesWrapper}>
      <div className={titleStyles}>{name}</div>
      <div className={messagesListWrapper}>
        {messages.map((message) => (
          <Message key={message.id} {...message} />
        ))}
      </div>
      <textarea
        className={newMessageTextarea}
        placeholder="Enter your message..."
      />
    </div>
  );
};
