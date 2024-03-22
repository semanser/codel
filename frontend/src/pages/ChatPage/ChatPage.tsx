import * as Tabs from "@radix-ui/react-tabs";
import { useNavigate, useParams } from "react-router-dom";

import { Messages } from "@/components/Messages/Messages";
import { Panel } from "@/components/Panel/Panel";
import dockerSvg from "@/assets/docker.svg";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsPillStyles,
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
      const result = await createFlowMutation({
        query: message,
      });

      const flowId = result?.data?.createFlow.id;
      if (flowId) {
        navigate(`/chat/${flowId}`, { replace: true });
      }
    } else {
      createTaskMutation({
        id: id,
        query: message,
      });
    }
  };

  const containerName = data?.flow?.containerName;

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
              {containerName && (<div className={tabsPillStyles}>
                <img src={dockerSvg} alt="Docker" width="14" height="14" />
                {containerName}
              </div>)}
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
            <Terminal id={isNewFlow ? "" : id} />
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
