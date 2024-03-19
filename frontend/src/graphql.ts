import { devtoolsExchange } from "@urql/devtools";
import { cacheExchange } from "@urql/exchange-graphcache";
import { createClient as createWSClient } from "graphql-ws";
import { createClient, fetchExchange, subscriptionExchange } from "urql";

import schema from "../generated/graphql.schema.json";

export const cache = cacheExchange({
  schema: schema,
  updates: {
    Mutation: {},
  },
  keys: {},
});

const wsClient = createWSClient({
  url: "ws://" + import.meta.env.VITE_API_URL + "/graphql",
});

export const graphqlClient = createClient({
  url: "http://" + import.meta.env.VITE_API_URL + "/graphql",
  fetchOptions: {},
  exchanges: [
    devtoolsExchange,
    cache,
    fetchExchange,
    subscriptionExchange({
      forwardSubscription(request) {
        const input = { ...request, query: request.query || "" };
        return {
          subscribe(sink) {
            const unsubscribe = wsClient.subscribe(input, sink);

            wsClient.on("error", (error) => {
              console.error("The subscription errored:", error);
            });

            return { unsubscribe };
          },
        };
      },
    }),
  ],
});
