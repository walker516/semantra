import { ApolloServer, gql } from "apollo-server";
import fetch from "node-fetch";

// スキーマ定義
const typeDefs = gql`
  type User {
    id: ID!
    name: String!
  }

  type Query {
    user(id: ID!): User
  }
`;

// リゾルバ定義
const resolvers = {
  Query: {
    user: async (_: any, { id }: { id: string }) => {
      const res = await fetch(`http://localhost:8080/api/user/${id}`);
      return await res.json();
    },
  },
};

// サーバ起動
const server = new ApolloServer({ typeDefs, resolvers });

server.listen({ port: 4000 }).then(({ url }) => {
  console.log(`🚀 BFF GraphQL Server ready at ${url}`);
});
