import { Task } from "@/generated/graphql";

import { Message } from "./Message/Message";
import {
  messagesListWrapper,
  messagesWrapper,
  newMessageTextarea,
  titleStyles,
} from "./Messages.css";

type MessagesProps = {
  tasks: Task[];
  name: string;
  onSubmit: (message: string) => void;
};

export const Messages = ({ tasks, name, onSubmit }: MessagesProps) => {
  const messages =
    tasks.map((task) => ({
      id: task.id,
      message: task.message,
      time: task.createdAt,
      status: task.status,
      type: task.type,
      output: "Test output",
    })) ?? [];

  const handleKeyPress = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();

      const message = e.currentTarget.value;

      e.currentTarget.value = "";

      onSubmit(message);
    }
  };

  return (
    <div className={messagesWrapper}>
      {name && <div className={titleStyles}>{name}</div>}
      <div className={messagesListWrapper}>
        {messages.map((message) => (
          <Message key={message.id} {...message} />
        ))}
      </div>
      <textarea
        autoFocus
        className={newMessageTextarea}
        placeholder="Enter your message..."
        onKeyPress={handleKeyPress}
      />
    </div>
  );
};
