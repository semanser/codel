import { useEffect, useRef } from "react";

import { FlowStatus, Model, Task } from "@/generated/graphql";

import { Button } from "../Button/Button";
import { Message } from "./Message/Message";
import {
  messagesListWrapper,
  messagesWrapper,
  modelStyles,
  newMessageTextarea,
  titleStyles,
} from "./Messages.css";

type MessagesProps = {
  tasks: Task[];
  name: string;
  onSubmit: (message: string) => void;
  onFlowStop: () => void;
  flowStatus?: FlowStatus;
  isNew?: boolean;
  model?: Model;
};

export const Messages = ({
  tasks,
  name,
  flowStatus,
  onSubmit,
  isNew,
  onFlowStop,
  model,
}: MessagesProps) => {
  const messages =
    tasks.map((task) => ({
      id: task.id,
      message: task.message,
      time: task.createdAt,
      status: task.status,
      type: task.type,
      output: task.results,
    })) ?? [];

  const messagesRef = useRef<HTMLDivElement>(null);
  const autoScrollEnabledRef = useRef(true);

  const handleKeyPress = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();

      const message = e.currentTarget.value;

      e.currentTarget.value = "";

      onSubmit(message);
    }
  };

  useEffect(() => {
    const messagesDiv = messagesRef.current;
    if (!messagesDiv) return;

    const scrollHandler = () => {
      if (
        messagesDiv.scrollTop + messagesDiv.clientHeight + 50 >=
        messagesDiv.scrollHeight
      ) {
        autoScrollEnabledRef.current = true;
      } else {
        autoScrollEnabledRef.current = false;
      }
    };

    messagesDiv.addEventListener("scroll", scrollHandler);

    return () => {
      messagesDiv.removeEventListener("scroll", scrollHandler);
    };
  }, []);

  useEffect(() => {
    const messagesDiv = messagesRef.current;
    if (!messagesDiv) return;

    if (autoScrollEnabledRef.current) {
      messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }
  }, [tasks]);

  const isFlowFinished = flowStatus === FlowStatus.Finished;

  return (
    <div className={messagesWrapper}>
      {name && (
        <div className={titleStyles}>
          {name}
          <span className={modelStyles}>{` - ${model?.id}`}</span>{" "}
          {isFlowFinished ? (
            " (Finished)"
          ) : (
            <Button hierarchy="danger" size="small" onClick={onFlowStop}>
              Finish
            </Button>
          )}{" "}
        </div>
      )}
      <div className={messagesListWrapper} ref={messagesRef}>
        {messages.map((message) => (
          <Message key={message.id} {...message} />
        ))}
      </div>
      <textarea
        autoFocus
        className={newMessageTextarea}
        placeholder={
          isFlowFinished
            ? "The task is finished."
            : isNew
              ? "Type a new message to start the flow..."
              : "Type a message..."
        }
        onKeyPress={handleKeyPress}
        disabled={isFlowFinished}
      />
    </div>
  );
};
