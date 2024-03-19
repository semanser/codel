import * as Tabs from "@radix-ui/react-tabs";

import {
  Message,
  MessageStatus,
  MessageType,
} from "@/components/Message/Message";
import { Panel } from "@/components/Panel/Panel";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsRootStyles,
  tabsTriggerStyles,
} from "@/components/Tabs/Tabs.css";
import { Terminal } from "@/components/Terminal/Terminal";

import { messagesWrapper, titleStyles, wrapperStyles } from "./ChatPage.css";

const fakeData = {
  title: "This is a chat",
  messages: [
    {
      id: 1,
      message: "This is a test message",
      time: new Date("2024-01-10"),
      type: MessageType.Browser,
      output: "This is the output of the message",
      status: MessageStatus.Finished,
    },
    {
      id: 2,
      message:
        "This is a some pretty long message. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum placerat felis ante, non semper mi hendrerit id. Praesent sodales est dui, ut semper sem consectetur nec. Praesent vitae euismod metus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Maecenas vitae ante interdum erat blandit eleifend.",
      time: new Date("2024-03-18"),
      output: "This is the output of the message",
      type: MessageType.Terminal,
      status: MessageStatus.Finished,
    },
    {
      id: 3,
      message: "This is some random message",
      time: new Date("2024-03-18"),
      output: "This is the output of the message",
      type: MessageType.Code,
      status: MessageStatus.Finished,
    },
    {
      id: 4,
      message: "This is some ask message",
      time: new Date("2024-03-18"),
      output: "This is the output of the message",
      type: MessageType.Ask,
      status: MessageStatus.Failed,
    },
    {
      id: 5,
      message: "This task is done",
      time: new Date("2024-03-18"),
      output: "This is the output of the message",
      type: MessageType.Done,
      status: MessageStatus.InProgress,
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
      <Panel>
        <Tabs.Root className={tabsRootStyles} defaultValue="terminal">
          <Tabs.List className={tabsListStyles}>
            <Tabs.Trigger className={tabsTriggerStyles} value="terminal">
              Terminal
            </Tabs.Trigger>
            <Tabs.Trigger
              className={tabsTriggerStyles}
              value="browser"
              disabled
            >
              Browser (Soon)
            </Tabs.Trigger>
            <Tabs.Trigger className={tabsTriggerStyles} value="code" disabled>
              Code (Soon)
            </Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content className={tabsContentStyles} value="terminal">
            <Terminal
              options={{
                allowProposedApi: true,
              }}
            />
          </Tabs.Content>
          <Tabs.Content className={tabsContentStyles} value="browser">
            browser
          </Tabs.Content>
          <Tabs.Content className={tabsContentStyles} value="code">
            code
          </Tabs.Content>
        </Tabs.Root>
      </Panel>
    </div>
  );
};
