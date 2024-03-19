import {
  avatarStyles,
  messageStyles,
  rightColumnStyles,
  timeStyles,
  wrapperStyles,
} from "./Message.css";

type MessageProps = {
  message: string;
  time: Date;
};

export const Message = ({ time, message }: MessageProps) => {
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
        <div className={timeStyles}>{time.toDateString()}</div>
        <div className={messageStyles}>{message}</div>
      </div>
    </div>
  );
};
