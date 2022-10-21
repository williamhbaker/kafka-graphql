# kafka-graphql

This demonstrates an easy way to query streaming data from Kafka using GraphQL, powered by Estuary
Flow and Hasura. Streaming events from Kafka are captured by Estuary Flow and materialized into a
Postgres database. The Hasura GraphQL engine automatically generates a GraphQL API can be used for
queries or real-time subscriptions of the data in the database.

To get started, you will need:

- An Estuary Flow account. Sign up at https://www.estuary.dev/.
- A Kafka cluster. Various managed providers exist for this, or try the `docker-compose.yaml` for an
  easy way to run Kafka for demo purposes in [KRaft
  mode](https://developer.confluent.io/learn/kraft/).
- A Postgres database. Managed Postgres can also be used here. The `docker-compose.yaml` also
  includes a Postgres instance that can be used as a demo.
- Hasura GraphQL engine. This can be run using docker (again, in the `docker-compose.yaml`), or see
  Hasura's [Getting Started](https://hasura.io/docs/latest/getting-started/index/) docs.

### Modeling Your Data

Representing streaming event data in a useful way can be challenging. Estuary Flow provides many
different ways to model this data, like
[reductions](https://docs.estuary.dev/concepts/schemas/#reductions) and
[derivations](https://docs.estuary.dev/concepts/derivations/). For this example, we are modeling a
stream of stock trade information, containing stock ticker symbols, trade price, trade volume, and
the time of the trade. We'll use reductions to find the total volume traded for each symbol and the
last traded price. A sample data generator is included in this repository, and the rest of this
example is based on this data.

### Configuring Estuary Flow

Estuary Flow will be configured with a Capture to ingest the data from Kafka, and a Materialization
to add the data to Postgres.

Both of these can be created in the Flow UI, or using the `flowctl` CLI. The `flow.yaml` file contains a full example of creating these tasks

For more details about setting up Flow, see the docs for [Kafka capture connector](https://docs.estuary.dev/reference/Connectors/capture-connectors/apache-kafka/) and [Postgres materialization connector](https://docs.estuary.dev/reference/Connectors/materialization-connectors/PostgreSQL/).

### Configuring Kafka and Postgres

Once the capture and materialization are published and running, no further configuration is required for Kafka or Postgres. The Postgres materialization will automatically create a table in your database that you specificy in the materialization configuration.

### Querying the Data with GraphQL

To connect the Hasura GraphQL engine to your Postgres database, refer to the [Hasura docs](https://hasura.io/docs/latest/index/). You will need to connect Hasura to the database using its database URL, and then track the table that your materialization created. Once that's done, you can use the Hasura Engine console's API explorer to build queries over your data as it is being materialized in real time.

For example, a query to list all symbols along with their summed trade volume and most recent trade price would be:

```
query {
  trades {
    price
    symbol
    Size
  }
}
```

It would produce data in the form of:

```json
{
  "data": {
    "trades": [
      {
        "price": 61.23,
        "symbol": "META",
        "Size": 237578
      },
      {
        "price": 4.03,
        "symbol": "MSFT",
        "Size": 233248
      },
      {
        "price": 84.85,
        "symbol": "NVDA",
        "Size": 238257
      },
      ...etc
    ]
  }
}
```

Real-time queries are possible using GraphQL subscriptions. A subscription query to list the same data that will update in real-time would be:

```
subscription {
  trades {
    price
    symbol
    Size
  }
}
```

The data returned will have the same shape, and be continuously updated via the GraphQL subscription.
