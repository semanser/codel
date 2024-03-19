import { formatDistance } from "date-fns";
import { useState } from "react";

import { Icon } from "@/components/Icon/Icon";
import { TaskStatus } from "@/generated/graphql";

import { Button } from "../Button/Button";
import {
  avatarStyles,
  contentStyles,
  iconStyles,
  messageStyles,
  outputStyles,
  rightColumnStyles,
  timeStyles,
  wrapperStyles,
} from "./Message.css";

export enum MessageType {
  Browser,
  Terminal,
  Code,
  Ask,
  Done,
}

type MessageProps = {
  message: string;
  time: Date;
  type: MessageType;
  status: TaskStatus;
  output: string;
};

export const Message = ({
  time,
  message,
  type,
  status,
  output,
}: MessageProps) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleExpand = () => {
    setIsExpanded((prev) => !prev);
  };

  return (
    <div className={wrapperStyles}>
      <img
        src="https://via.placeholder.com/40"
        alt="avatar"
        className={avatarStyles}
        width="40"
        height="40"
      />
      <div className={rightColumnStyles}>
        <div className={timeStyles}>
          {formatDistance(new Date(time), new Date(), { addSuffix: true })}
        </div>
        <div
          className={
            status !== TaskStatus.Failed
              ? messageStyles.Regular
              : messageStyles.Failed
          }
          onClick={toggleExpand}
        >
          <div className={contentStyles}>
            <span
              className={
                status !== TaskStatus.Failed
                  ? iconStyles.Regular
                  : iconStyles.Failed
              }
            >
              {getIcon(type)}
            </span>
            <div>{message}</div>
          </div>
          {status === TaskStatus.InProgress && (
            <Button size="small" hierarchy="danger">
              Stop
            </Button>
          )}
        </div>
        {isExpanded && <div className={outputStyles}>{output}</div>}
      </div>
    </div>
  );
};

const getIcon = (type: MessageType) => {
  let icon = null;

  switch (type) {
    case MessageType.Browser:
      icon = <Icon.Browser />;
      break;
    case MessageType.Terminal:
      icon = <Icon.Terminal />;
      break;
    case MessageType.Code:
      icon = <Icon.Code />;
      break;
    case MessageType.Ask:
      icon = <Icon.MessageQuestion />;
      break;
    case MessageType.Done:
      icon = <Icon.CheckCircle />;
      break;
  }

  return icon;
};
