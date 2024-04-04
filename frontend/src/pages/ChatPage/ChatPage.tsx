import * as Tabs from "@radix-ui/react-tabs";
import { useLocalStorage } from "@uidotdev/usehooks";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import Browser from "@/components/Browser/Browser";
import { Button } from "@/components/Button/Button";
import { Icon } from "@/components/Icon/Icon";
import { Messages } from "@/components/Messages/Messages";
import { Panel } from "@/components/Panel/Panel";
import {
  tabsContentStyles,
  tabsListStyles,
  tabsRootStyles,
  tabsTriggerStyles,
} from "@/components/Tabs/Tabs.css";
import { Terminal } from "@/components/Terminal/Terminal";
import { Tooltip } from "@/components/Tooltip/Tooltip";
import {
  Model,
  useBrowserUpdatedSubscription,
  useCreateFlowMutation,
  useCreateTaskMutation,
  useFinishFlowMutation,
  useFlowQuery,
  useFlowUpdatedSubscription,
  useTaskAddedSubscription,
  useTerminalLogsAddedSubscription,
} from "@/generated/graphql";

import {
  followButtonStyles,
  leftColumnStyles,
  tabsStyles,
  wrapperStyles,
} from "./ChatPage.css";

export const ChatPage = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [, createFlowMutation] = useCreateFlowMutation();
  const [, createTaskMutation] = useCreateTaskMutation();
  const [, finishFlowMutation] = useFinishFlowMutation();
  const isNewFlow = !id || id === "new";
  const [isFollowingTabs, setIsFollowingTabs] = useLocalStorage(
    "isFollowingTabs",
    true,
  );
  const [selectedModel] = useLocalStorage<Model>("model");
  const [activeTab, setActiveTab] = useState("terminal");

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
  const browser = !isStaleData ? data?.flow.browser : undefined;
  const model = !isStaleData ? data?.flow.model : undefined;

  useBrowserUpdatedSubscription(
    {
      variables: { flowId: Number(id) },
      pause: isNewFlow,
    },
    () => {
      if (isFollowingTabs) {
        setActiveTab("browser");
      }
    },
  );

  useTerminalLogsAddedSubscription(
    {
      variables: { flowId: Number(id) },
      pause: isNewFlow,
    },
    () => {
      if (isFollowingTabs) {
        setActiveTab("terminal");
      }
    },
  );

  useTaskAddedSubscription({
    variables: { flowId: Number(id) },
    pause: isNewFlow,
  });

  useFlowUpdatedSubscription({
    variables: { flowId: Number(id) },
    pause: isNewFlow,
  });

  const handleSubmit = async (message: string) => {
    if (isNewFlow && selectedModel.id) {
      const result = await createFlowMutation({
        modelProvider: selectedModel.provider,
        modelId: selectedModel.id,
      });

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

  const handleFlowStop = () => {
    finishFlowMutation({ flowId: id });
  };

  const handleChangeIsFollowingTabs = () => {
    setIsFollowingTabs(!isFollowingTabs);
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
          onFlowStop={handleFlowStop}
          model={model}
        />
      </Panel>
      <Panel>
        <Tabs.Root
          className={tabsRootStyles}
          value={activeTab}
          onValueChange={setActiveTab}
        >
          <Tabs.List className={tabsListStyles}>
            <div className={tabsStyles}>
              <div className={leftColumnStyles}>
                <Tabs.Trigger className={tabsTriggerStyles} value="terminal">
                  Terminal
                </Tabs.Trigger>
                <Tabs.Trigger className={tabsTriggerStyles} value="browser">
                  Browser
                </Tabs.Trigger>
                <Tabs.Trigger
                  className={tabsTriggerStyles}
                  value="code"
                  disabled
                >
                  Code (Soon)
                </Tabs.Trigger>
              </div>

              <Tooltip
                content={
                  <>
                    Following the active tab is{" "}
                    <b>{isFollowingTabs ? "enabled" : "disabled"}</b>
                  </>
                }
              >
                <Button
                  size="small"
                  hierarchy={isFollowingTabs ? "primary" : "secondary"}
                  className={followButtonStyles}
                  onClick={handleChangeIsFollowingTabs}
                >
                  {isFollowingTabs ? <Icon.Eye /> : <Icon.EyeOff />}
                </Button>
              </Tooltip>
            </div>
          </Tabs.List>
          <Tabs.Content className={tabsContentStyles} value="terminal">
            <Terminal
              id={isNewFlow ? "" : id}
              status={status}
              title={terminal?.containerName}
              logs={terminal?.logs ?? []}
              isRunning={terminal?.connected}
            />
          </Tabs.Content>
          <Tabs.Content className={tabsContentStyles} value="browser">
            <Browser
              url={browser?.url || undefined}
              screenshotUrl={browser?.screenshotUrl ?? ""}
            />
          </Tabs.Content>
          <Tabs.Content className={tabsContentStyles} value="code">
            code
          </Tabs.Content>
        </Tabs.Root>
      </Panel>
    </div>
  );
};
