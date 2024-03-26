import gql from 'graphql-tag';
import * as Urql from 'urql';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  JSON: { input: any; output: any; }
  Time: { input: any; output: any; }
  Uint: { input: any; output: any; }
};

export type Flow = {
  __typename?: 'Flow';
  id: Scalars['Uint']['output'];
  name: Scalars['String']['output'];
  status: FlowStatus;
  tasks: Array<Task>;
  terminal: Terminal;
};

export enum FlowStatus {
  Finished = 'finished',
  InProgress = 'inProgress'
}

export type Log = {
  __typename?: 'Log';
  id: Scalars['Uint']['output'];
  text: Scalars['String']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  _exec: Scalars['String']['output'];
  createFlow: Flow;
  createTask: Task;
};


export type Mutation_ExecArgs = {
  command: Scalars['String']['input'];
  containerId: Scalars['String']['input'];
};


export type MutationCreateTaskArgs = {
  flowId: Scalars['Uint']['input'];
  query: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  flow: Flow;
  flows: Array<Flow>;
};


export type QueryFlowArgs = {
  id: Scalars['Uint']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  flowUpdated: Flow;
  taskAdded: Task;
  taskUpdated: Task;
  terminalLogsAdded: Log;
};


export type SubscriptionFlowUpdatedArgs = {
  flowId: Scalars['Uint']['input'];
};


export type SubscriptionTaskAddedArgs = {
  flowId: Scalars['Uint']['input'];
};


export type SubscriptionTerminalLogsAddedArgs = {
  flowId: Scalars['Uint']['input'];
};

export type Task = {
  __typename?: 'Task';
  args: Scalars['JSON']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['Uint']['output'];
  message: Scalars['String']['output'];
  results: Scalars['JSON']['output'];
  status: TaskStatus;
  type: TaskType;
};

export enum TaskStatus {
  Failed = 'failed',
  Finished = 'finished',
  InProgress = 'inProgress',
  Stopped = 'stopped'
}

export enum TaskType {
  Ask = 'ask',
  Browser = 'browser',
  Code = 'code',
  Done = 'done',
  Input = 'input',
  Terminal = 'terminal'
}

export type Terminal = {
  __typename?: 'Terminal';
  connected: Scalars['Boolean']['output'];
  containerName: Scalars['String']['output'];
  logs: Array<Log>;
};

export const FlowOverviewFragmentFragmentDoc = gql`
    fragment flowOverviewFragment on Flow {
  id
  name
  status
}
    `;
export const LogFragmentFragmentDoc = gql`
    fragment logFragment on Log {
  text
  id
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
  status
  terminal {
    containerName
    connected
    logs {
      ...logFragment
    }
  }
  tasks {
    ...taskFragment
  }
}
    ${LogFragmentFragmentDoc}
${TaskFragmentFragmentDoc}`;
export const FlowsDocument = gql`
    query flows {
  flows {
    ...flowOverviewFragment
  }
}
    ${FlowOverviewFragmentFragmentDoc}`;

export function useFlowsQuery(options?: Omit<Urql.UseQueryArgs<FlowsQueryVariables>, 'query'>) {
  return Urql.useQuery<FlowsQuery, FlowsQueryVariables>({ query: FlowsDocument, ...options });
};
export const FlowDocument = gql`
    query flow($id: Uint!) {
  flow(id: $id) {
    ...flowFragment
  }
}
    ${FlowFragmentFragmentDoc}`;

export function useFlowQuery(options: Omit<Urql.UseQueryArgs<FlowQueryVariables>, 'query'>) {
  return Urql.useQuery<FlowQuery, FlowQueryVariables>({ query: FlowDocument, ...options });
};
export const CreateFlowDocument = gql`
    mutation createFlow {
  createFlow {
    id
    name
  }
}
    `;

export function useCreateFlowMutation() {
  return Urql.useMutation<CreateFlowMutation, CreateFlowMutationVariables>(CreateFlowDocument);
};
export const CreateTaskDocument = gql`
    mutation createTask($flowId: Uint!, $query: String!) {
  createTask(flowId: $flowId, query: $query) {
    ...taskFragment
  }
}
    ${TaskFragmentFragmentDoc}`;

export function useCreateTaskMutation() {
  return Urql.useMutation<CreateTaskMutation, CreateTaskMutationVariables>(CreateTaskDocument);
};
export const TaskAddedDocument = gql`
    subscription taskAdded($flowId: Uint!) {
  taskAdded(flowId: $flowId) {
    ...taskFragment
  }
}
    ${TaskFragmentFragmentDoc}`;

export function useTaskAddedSubscription<TData = TaskAddedSubscription>(options: Omit<Urql.UseSubscriptionArgs<TaskAddedSubscriptionVariables>, 'query'>, handler?: Urql.SubscriptionHandler<TaskAddedSubscription, TData>) {
  return Urql.useSubscription<TaskAddedSubscription, TData, TaskAddedSubscriptionVariables>({ query: TaskAddedDocument, ...options }, handler);
};
export const TerminalLogsAddedDocument = gql`
    subscription terminalLogsAdded($flowId: Uint!) {
  terminalLogsAdded(flowId: $flowId) {
    ...logFragment
  }
}
    ${LogFragmentFragmentDoc}`;

export function useTerminalLogsAddedSubscription<TData = TerminalLogsAddedSubscription>(options: Omit<Urql.UseSubscriptionArgs<TerminalLogsAddedSubscriptionVariables>, 'query'>, handler?: Urql.SubscriptionHandler<TerminalLogsAddedSubscription, TData>) {
  return Urql.useSubscription<TerminalLogsAddedSubscription, TData, TerminalLogsAddedSubscriptionVariables>({ query: TerminalLogsAddedDocument, ...options }, handler);
};
export const FlowUpdatedDocument = gql`
    subscription flowUpdated($flowId: Uint!) {
  flowUpdated(flowId: $flowId) {
    id
    name
    terminal {
      containerName
      connected
    }
  }
}
    `;

export function useFlowUpdatedSubscription<TData = FlowUpdatedSubscription>(options: Omit<Urql.UseSubscriptionArgs<FlowUpdatedSubscriptionVariables>, 'query'>, handler?: Urql.SubscriptionHandler<FlowUpdatedSubscription, TData>) {
  return Urql.useSubscription<FlowUpdatedSubscription, TData, FlowUpdatedSubscriptionVariables>({ query: FlowUpdatedDocument, ...options }, handler);
};
export type FlowOverviewFragmentFragment = { __typename?: 'Flow', id: any, name: string, status: FlowStatus };

export type FlowsQueryVariables = Exact<{ [key: string]: never; }>;


export type FlowsQuery = { __typename?: 'Query', flows: Array<{ __typename?: 'Flow', id: any, name: string, status: FlowStatus }> };

export type TaskFragmentFragment = { __typename?: 'Task', id: any, type: TaskType, message: string, status: TaskStatus, args: any, results: any, createdAt: any };

export type LogFragmentFragment = { __typename?: 'Log', text: string, id: any };

export type FlowFragmentFragment = { __typename?: 'Flow', id: any, name: string, status: FlowStatus, terminal: { __typename?: 'Terminal', containerName: string, connected: boolean, logs: Array<{ __typename?: 'Log', text: string, id: any }> }, tasks: Array<{ __typename?: 'Task', id: any, type: TaskType, message: string, status: TaskStatus, args: any, results: any, createdAt: any }> };

export type FlowQueryVariables = Exact<{
  id: Scalars['Uint']['input'];
}>;


export type FlowQuery = { __typename?: 'Query', flow: { __typename?: 'Flow', id: any, name: string, status: FlowStatus, terminal: { __typename?: 'Terminal', containerName: string, connected: boolean, logs: Array<{ __typename?: 'Log', text: string, id: any }> }, tasks: Array<{ __typename?: 'Task', id: any, type: TaskType, message: string, status: TaskStatus, args: any, results: any, createdAt: any }> } };

export type CreateFlowMutationVariables = Exact<{ [key: string]: never; }>;


export type CreateFlowMutation = { __typename?: 'Mutation', createFlow: { __typename?: 'Flow', id: any, name: string } };

export type CreateTaskMutationVariables = Exact<{
  flowId: Scalars['Uint']['input'];
  query: Scalars['String']['input'];
}>;


export type CreateTaskMutation = { __typename?: 'Mutation', createTask: { __typename?: 'Task', id: any, type: TaskType, message: string, status: TaskStatus, args: any, results: any, createdAt: any } };

export type TaskAddedSubscriptionVariables = Exact<{
  flowId: Scalars['Uint']['input'];
}>;


export type TaskAddedSubscription = { __typename?: 'Subscription', taskAdded: { __typename?: 'Task', id: any, type: TaskType, message: string, status: TaskStatus, args: any, results: any, createdAt: any } };

export type TerminalLogsAddedSubscriptionVariables = Exact<{
  flowId: Scalars['Uint']['input'];
}>;


export type TerminalLogsAddedSubscription = { __typename?: 'Subscription', terminalLogsAdded: { __typename?: 'Log', text: string, id: any } };

export type FlowUpdatedSubscriptionVariables = Exact<{
  flowId: Scalars['Uint']['input'];
}>;


export type FlowUpdatedSubscription = { __typename?: 'Subscription', flowUpdated: { __typename?: 'Flow', id: any, name: string, terminal: { __typename?: 'Terminal', containerName: string, connected: boolean } } };

import { IntrospectionQuery } from 'graphql';
export default {
  "__schema": {
    "queryType": {
      "name": "Query"
    },
    "mutationType": {
      "name": "Mutation"
    },
    "subscriptionType": {
      "name": "Subscription"
    },
    "types": [
      {
        "kind": "OBJECT",
        "name": "Flow",
        "fields": [
          {
            "name": "id",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "name",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "status",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "tasks",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "LIST",
                "ofType": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "OBJECT",
                    "name": "Task",
                    "ofType": null
                  }
                }
              }
            },
            "args": []
          },
          {
            "name": "terminal",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Terminal",
                "ofType": null
              }
            },
            "args": []
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Log",
        "fields": [
          {
            "name": "id",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "text",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Mutation",
        "fields": [
          {
            "name": "_exec",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": [
              {
                "name": "command",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              },
              {
                "name": "containerId",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          },
          {
            "name": "createFlow",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Flow",
                "ofType": null
              }
            },
            "args": []
          },
          {
            "name": "createTask",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Task",
                "ofType": null
              }
            },
            "args": [
              {
                "name": "flowId",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              },
              {
                "name": "query",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Query",
        "fields": [
          {
            "name": "flow",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Flow",
                "ofType": null
              }
            },
            "args": [
              {
                "name": "id",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          },
          {
            "name": "flows",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "LIST",
                "ofType": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "OBJECT",
                    "name": "Flow",
                    "ofType": null
                  }
                }
              }
            },
            "args": []
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Subscription",
        "fields": [
          {
            "name": "flowUpdated",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Flow",
                "ofType": null
              }
            },
            "args": [
              {
                "name": "flowId",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          },
          {
            "name": "taskAdded",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Task",
                "ofType": null
              }
            },
            "args": [
              {
                "name": "flowId",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          },
          {
            "name": "taskUpdated",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Task",
                "ofType": null
              }
            },
            "args": []
          },
          {
            "name": "terminalLogsAdded",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "OBJECT",
                "name": "Log",
                "ofType": null
              }
            },
            "args": [
              {
                "name": "flowId",
                "type": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "SCALAR",
                    "name": "Any"
                  }
                }
              }
            ]
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Task",
        "fields": [
          {
            "name": "args",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "createdAt",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "id",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "message",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "results",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "status",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "type",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          }
        ],
        "interfaces": []
      },
      {
        "kind": "OBJECT",
        "name": "Terminal",
        "fields": [
          {
            "name": "connected",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "containerName",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "SCALAR",
                "name": "Any"
              }
            },
            "args": []
          },
          {
            "name": "logs",
            "type": {
              "kind": "NON_NULL",
              "ofType": {
                "kind": "LIST",
                "ofType": {
                  "kind": "NON_NULL",
                  "ofType": {
                    "kind": "OBJECT",
                    "name": "Log",
                    "ofType": null
                  }
                }
              }
            },
            "args": []
          }
        ],
        "interfaces": []
      },
      {
        "kind": "SCALAR",
        "name": "Any"
      }
    ],
    "directives": []
  }
} as unknown as IntrospectionQuery;