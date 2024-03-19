import { IntrospectionQuery } from "graphql";
import gql from "graphql-tag";
import * as Urql from "urql";

export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  JSON: { input: any; output: any };
  Time: { input: any; output: any };
  Uint: { input: any; output: any };
};

export type Flow = {
  __typename?: "Flow";
  id: Scalars["Uint"]["output"];
  name: Scalars["String"]["output"];
  tasks: Array<Task>;
};

export type Mutation = {
  __typename?: "Mutation";
  createFlow: Flow;
  createTask: Task;
  stopTask: Task;
};

export type MutationCreateTaskArgs = {
  id: Scalars["Uint"]["input"];
  query: Scalars["String"]["input"];
};

export type MutationStopTaskArgs = {
  id: Scalars["Uint"]["input"];
};

export type Query = {
  __typename?: "Query";
  flow: Flow;
  flows: Array<Flow>;
};

export type QueryFlowArgs = {
  id: Scalars["Uint"]["input"];
};

export type Subscription = {
  __typename?: "Subscription";
  taskAdded: Task;
  taskUpdated: Task;
};

export type Task = {
  __typename?: "Task";
  args: Scalars["JSON"]["output"];
  createdAt: Scalars["Time"]["output"];
  id: Scalars["Uint"]["output"];
  message: Scalars["String"]["output"];
  results: Scalars["JSON"]["output"];
  status: TaskStatus;
  type: TaskType;
};

export enum TaskStatus {
  Failed = "failed",
  Finished = "finished",
  InProgress = "inProgress",
  Stopped = "stopped",
}

export enum TaskType {
  Action = "action",
  Input = "input",
}

export const FlowOverviewFragmentFragmentDoc = gql`
  fragment flowOverviewFragment on Flow {
    id
    name
  }
`;
export const TaskFragmentFragmentDoc = gql`
  fragment taskFragment on Task {
    id
    type
    message
    status
    args
    results
    createdAt
  }
`;
export const FlowFragmentFragmentDoc = gql`
  fragment flowFragment on Flow {
    id
    name
    tasks {
      ...taskFragment
    }
  }
  ${TaskFragmentFragmentDoc}
`;
export const FlowsDocument = gql`
  query flows {
    flows {
      ...flowOverviewFragment
    }
  }
  ${FlowOverviewFragmentFragmentDoc}
`;

export function useFlowsQuery(
  options?: Omit<Urql.UseQueryArgs<FlowsQueryVariables>, "query">,
) {
  return Urql.useQuery<FlowsQuery, FlowsQueryVariables>({
    query: FlowsDocument,
    ...options,
  });
}
export const FlowDocument = gql`
  query flow($id: Uint!) {
    flow(id: $id) {
      ...flowFragment
    }
  }
  ${FlowFragmentFragmentDoc}
`;

export function useFlowQuery(
  options: Omit<Urql.UseQueryArgs<FlowQueryVariables>, "query">,
) {
  return Urql.useQuery<FlowQuery, FlowQueryVariables>({
    query: FlowDocument,
    ...options,
  });
}
export const CreateFlowDocument = gql`
  mutation createFlow {
    createFlow {
      id
    }
  }
`;

export function useCreateFlowMutation() {
  return Urql.useMutation<CreateFlowMutation, CreateFlowMutationVariables>(
    CreateFlowDocument,
  );
}
export const CreateTaskDocument = gql`
  mutation createTask($id: Uint!, $query: String!) {
    createTask(id: $id, query: $query) {
      ...taskFragment
    }
  }
  ${TaskFragmentFragmentDoc}
`;

export function useCreateTaskMutation() {
  return Urql.useMutation<CreateTaskMutation, CreateTaskMutationVariables>(
    CreateTaskDocument,
  );
}
export type FlowOverviewFragmentFragment = {
  __typename?: "Flow";
  id: any;
  name: string;
};

export type FlowsQueryVariables = Exact<{ [key: string]: never }>;

export type FlowsQuery = {
  __typename?: "Query";
  flows: Array<{ __typename?: "Flow"; id: any; name: string }>;
};

export type TaskFragmentFragment = {
  __typename?: "Task";
  id: any;
  type: TaskType;
  message: string;
  status: TaskStatus;
  args: any;
  results: any;
  createdAt: any;
};

export type FlowFragmentFragment = {
  __typename?: "Flow";
  id: any;
  name: string;
  tasks: Array<{
    __typename?: "Task";
    id: any;
    type: TaskType;
    message: string;
    status: TaskStatus;
    args: any;
    results: any;
    createdAt: any;
  }>;
};

export type FlowQueryVariables = Exact<{
  id: Scalars["Uint"]["input"];
}>;

export type FlowQuery = {
  __typename?: "Query";
  flow: {
    __typename?: "Flow";
    id: any;
    name: string;
    tasks: Array<{
      __typename?: "Task";
      id: any;
      type: TaskType;
      message: string;
      status: TaskStatus;
      args: any;
      results: any;
      createdAt: any;
    }>;
  };
};

export type CreateFlowMutationVariables = Exact<{ [key: string]: never }>;

export type CreateFlowMutation = {
  __typename?: "Mutation";
  createFlow: { __typename?: "Flow"; id: any };
};

export type CreateTaskMutationVariables = Exact<{
  id: Scalars["Uint"]["input"];
  query: Scalars["String"]["input"];
}>;

export type CreateTaskMutation = {
  __typename?: "Mutation";
  createTask: {
    __typename?: "Task";
    id: any;
    type: TaskType;
    message: string;
    status: TaskStatus;
    args: any;
    results: any;
    createdAt: any;
  };
};

export default {
  __schema: {
    queryType: {
      name: "Query",
    },
    mutationType: {
      name: "Mutation",
    },
    subscriptionType: {
      name: "Subscription",
    },
    types: [
      {
        kind: "OBJECT",
        name: "Flow",
        fields: [
          {
            name: "id",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "name",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "tasks",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "LIST",
                ofType: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "OBJECT",
                    name: "Task",
                    ofType: null,
                  },
                },
              },
            },
            args: [],
          },
        ],
        interfaces: [],
      },
      {
        kind: "OBJECT",
        name: "Mutation",
        fields: [
          {
            name: "createFlow",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Flow",
                ofType: null,
              },
            },
            args: [],
          },
          {
            name: "createTask",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Task",
                ofType: null,
              },
            },
            args: [
              {
                name: "id",
                type: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "SCALAR",
                    name: "Any",
                  },
                },
              },
              {
                name: "query",
                type: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "SCALAR",
                    name: "Any",
                  },
                },
              },
            ],
          },
          {
            name: "stopTask",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Task",
                ofType: null,
              },
            },
            args: [
              {
                name: "id",
                type: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "SCALAR",
                    name: "Any",
                  },
                },
              },
            ],
          },
        ],
        interfaces: [],
      },
      {
        kind: "OBJECT",
        name: "Query",
        fields: [
          {
            name: "flow",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Flow",
                ofType: null,
              },
            },
            args: [
              {
                name: "id",
                type: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "SCALAR",
                    name: "Any",
                  },
                },
              },
            ],
          },
          {
            name: "flows",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "LIST",
                ofType: {
                  kind: "NON_NULL",
                  ofType: {
                    kind: "OBJECT",
                    name: "Flow",
                    ofType: null,
                  },
                },
              },
            },
            args: [],
          },
        ],
        interfaces: [],
      },
      {
        kind: "OBJECT",
        name: "Subscription",
        fields: [
          {
            name: "taskAdded",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Task",
                ofType: null,
              },
            },
            args: [],
          },
          {
            name: "taskUpdated",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "OBJECT",
                name: "Task",
                ofType: null,
              },
            },
            args: [],
          },
        ],
        interfaces: [],
      },
      {
        kind: "OBJECT",
        name: "Task",
        fields: [
          {
            name: "args",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "createdAt",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "id",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "message",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "results",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "status",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
          {
            name: "type",
            type: {
              kind: "NON_NULL",
              ofType: {
                kind: "SCALAR",
                name: "Any",
              },
            },
            args: [],
          },
        ],
        interfaces: [],
      },
      {
        kind: "SCALAR",
        name: "Any",
      },
    ],
    directives: [],
  },
} as unknown as IntrospectionQuery;
