import * as Tabs from "@radix-ui/react-tabs";
import { useParams } from "react-router-dom";

import { Messages } from "@/components/Messages/Messages";
import { Panel } from "@/components/Panel/Panel";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsRootStyles,
  tabsTriggerStyles,
} from "@/components/Tabs/Tabs.css";
import { Terminal } from "@/components/Terminal/Terminal";
import { useFlowQuery } from "@/generated/graphql";

import { wrapperStyles } from "./ChatPage.css";

export const ChatPage = () => {
  const { id } = useParams<{ id: string }>();
  const [{ data }] = useFlowQuery({
    pause: !id && id !== "new",
    variables: { id: id },
  });

  const isNew = id === "new";

  const tasks = id && !isNew ? data?.flow.tasks ?? [] : [];
  const name = id && !isNew ? data?.flow.name ?? "" : "";

  return (
    <div className={wrapperStyles}>
      <Panel>
        <Messages tasks={tasks} name={name} />
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
