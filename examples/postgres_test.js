import pgxpool from 'k6/x/pgxpool';

// The second argument is a PostgreSQL connection string, e.g.
// postgres://myuser:mypass@127.0.0.1:5432/postgres?sslmode=disable
const connString = "postgres://myuser:mypass@127.0.0.1:5432/postgres?sslmode=disable"
const minConns = 20;
const maxConns = 30;

const pool = pgxpool.open(connString, minConns, maxConns);

export function setup() {

}

export function teardown() {

}

export default function () {
  pgxpool.exec(pool, 'SELECT 1;');
  let result = pgxpool.query(pool, 'SELECT 1;');
  console.log(`result: ${JSON.stringify(result)}`);
}