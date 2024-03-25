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
  useFlowUpdatedSubscription,
  useTaskAddedSubscription,
  useTerminalLogsAddedSubscription,
} from "@/generated/graphql";

import { wrapperStyles } from "./ChatPage.css";

export const ChatPage = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [, createFlowMutation] = useCreateFlowMutation();
  const [, createTaskMutation] = useCreateTaskMutation();
  const isNewFlow = !id || id === "new";

  const [{ operation, data }] = useFlowQuery({
    pause: isNewFlow,
    variables: { id },
  });

  // https://github.com/urql-graphql/urql/issues/2507#issuecomment-1159281108
  const isStaleData = operation?.variables.id !== id;

  const tasks = !isStaleData ? data?.flow.tasks ?? [] : [];
  const name = !isStaleData ? data?.flow.name ?? "" : "";
  const status = !isStaleData ? data?.flow.status : undefined;
  const terminal = !isStaleData ? data?.flow.terminal : undefined;

  useTerminalLogsAddedSubscription({
    variables: { flowId: Number(id) },
    pause: isNewFlow,
  });

  useTaskAddedSubscription({
    variables: { flowId: Number(id) },
    pause: isNewFlow,
  });

  useFlowUpdatedSubscription({
    variables: { flowId: Number(id) },
    pause: isNewFlow,
  });

  const handleSubmit = async (message: string) => {
    if (isNewFlow) {
      const result = await createFlowMutation({});

      const flowId = result?.data?.createFlow.id;
      if (flowId) {
        navigate(`/chat/${flowId}`, { replace: true });

        createTaskMutation({
          flowId: flowId,
          query: message,
        });
      }
    } else {
      createTaskMutation({
        flowId: id,
        query: message,
      });
    }
  };

  return (
    <div className={wrapperStyles}>
      <Panel>
        <Messages
          tasks={tasks}
          name={name}
          onSubmit={handleSubmit}
          flowStatus={status}
          isNew={isNewFlow}
        />
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
              id={isNewFlow ? "" : id}
              status={status}
              title={terminal?.containerName}
              logs={terminal?.logs ?? []}
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
