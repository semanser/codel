import * as Tabs from "@radix-ui/react-tabs";
import { useNavigate, useParams } from "react-router-dom";

import { Messages } from "@/components/Messages/Messages";
import { Panel } from "@/components/Panel/Panel";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsRootStyles,
  tabsTriggerStyles,
} from "@/components/Tabs/Tabs.css";
import { Terminal } from "@/components/Terminal/Terminal";
import {
  useCreateFlowMutation,
  useCreateTaskMutation,
  useFlowQuery,
  useTaskAddedSubscription,
} from "@/generated/graphql";

import { wrapperStyles } from "./ChatPage.css";

export const ChatPage = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [, createFlowMutation] = useCreateFlowMutation();
  const [, createTaskMutation] = useCreateTaskMutation();
  const shouldPause = !id || id === "new";

  const [{ data }] = useFlowQuery({
    pause: shouldPause,
    variables: { id: id },
  });

  const isNew = id === "new";

  const tasks = id && !isNew ? data?.flow.tasks ?? [] : [];
  const name = id && !isNew ? data?.flow.name ?? "" : "";

  useTaskAddedSubscription({
    variables: { flowId: Number(id) },
    pause: shouldPause,
  });

  const handleSubmit = async (message: string) => {
    if (isNew) {
      const result = await createFlowMutation({});

      const flowId = result?.data?.createFlow.id;
      if (flowId) {
        navigate(`/chat/${flowId}`, { replace: true });

        setTimeout(() => {
          createTaskMutation({
            id: flowId,
            query: message,
          });
        }, 1000);
      }
    } else {
      createTaskMutation({
        id: id,
        query: message,
      });
    }
  };

  return (
    <div className={wrapperStyles}>
      <Panel>
        <Messages tasks={tasks} name={name} onSubmit={handleSubmit} />
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
            <Terminal key={id} id={isNew ? "" : id} />
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
