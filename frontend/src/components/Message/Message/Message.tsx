import { formatDistance } from "date-fns";

import { Icon } from "@/components/Icon/Icon";

import {
  avatarStyles,
  iconStyles,
  messageStyles,
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
};

export const Message = ({ time, message, type }: MessageProps) => {
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
        <div className={messageStyles}>
          <span className={iconStyles}>{getIcon(type)}</span>
          <div>{message}</div>
        </div>
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
