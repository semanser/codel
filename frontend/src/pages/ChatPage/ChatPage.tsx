import * as Tabs from "@radix-ui/react-tabs";
import { useParams } from "react-router-dom";

import { Message, MessageType } from "@/components/Message/Message";
import { Panel } from "@/components/Panel/Panel";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsRootStyles,
  tabsTriggerStyles,
} from "@/components/Tabs/Tabs.css";
import { Terminal } from "@/components/Terminal/Terminal";
import { useFlowQuery } from "@/generated/graphql";

import { messagesWrapper, titleStyles, wrapperStyles } from "./ChatPage.css";

export const ChatPage = () => {
  const { id } = useParams<{ id: string }>();
  const [{ data }] = useFlowQuery({
    pause: !id,
    variables: { id: id },
  });

  const messages =
    data?.flow?.tasks.map((task) => ({
      id: task.id,
      message: task.message,
      time: task.createdAt,
      status: task.status,
      // TODO Add the correct type and output
      type: MessageType.Terminal,
      output: "Test output",
    })) ?? [];

  return (
    <div className={wrapperStyles}>
      <Panel>
        <div className={titleStyles}>{data?.flow.name}</div>
        <div className={messagesWrapper}>
          {messages.map((message) => (
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
