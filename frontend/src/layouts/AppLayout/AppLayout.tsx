import { Outlet } from "react-router-dom";

import { Sidebar } from "@/components/Sidebar/Sidebar";
import { FlowStatus, useFlowsQuery } from "@/generated/graphql";

import { wrapperStyles } from "./AppLayout.css";

export const AppLayout = () => {
  const [{ data }] = useFlowsQuery();

  const sidebarItems =
    data?.flows.map((flow) => ({
      id: flow.id,
      title: flow.name,
      done: flow.status === FlowStatus.Finished,
    })) ?? [];

  return (
    <div className={wrapperStyles}>
      <Sidebar items={sidebarItems} />
      <Outlet />
    </div>
  );
};
