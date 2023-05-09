# xk6-sql

This is a [k6](https://github.com/grafana/k6) extension using the
[xk6](https://github.com/grafana/xk6) system.

Supported RDBMSs: `postgres`. See the [examples](examples)
directory for usage. Other RDBMSs are not supported, see
[details below](#support-for-other-rdbmss).

## Build

To build a `k6` binary with this plugin, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  xk6 build --with github.com/gcfabri/xk6-sql
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```bash
make
```
Once built, you can run your newly extended `k6` using:
```shell
 ./k6 run examples/postgres_test.js
 ```

## Example

```javascript
// script.js
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
  let result = pgxpool.query(pool, 'SELECT 1;');
  console.log(`result: ${JSON.stringify(result)}`);
}
```

Result output:

```shell
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: /Users/gabriel.fabri/Workspace/go-playground/cmd/k6/run.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0001] result: [{"?column?":1}]                      source=console

     █ setup

     █ teardown

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=436.43ms min=3.84µs med=7.59µs max=1.3s p(90)=1.04s p(95)=1.17s
     iterations...........: 1   0.760586/s
     vus..................: 1   min=1      max=1
     vus_max..............: 1   min=1      max=1


running (00m01.3s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m01.3s/10m0s  1/1 iters, 1 per VU
```

## See also

- [Load Testing SQL Databases with k6](https://k6.io/blog/load-testing-sql-databases-with-k6/)

### Support for other RDBMSs

Note that this project is not accepting support for additional SQL implementations
and RDBMSs. See the discussion in [issue #17](https://github.com/grafana/xk6-sql/issues/17)
for the reasoning.

You can build k6 binaries by simply specifying these project URLs in `xk6 build`.
E.g. `xk6 build --with github.com/gcfabri/xk6-sql`.
Please report any issues with these extensions in their respective GitHub issue trackers,
and not in gcfabri/xk6-sql.

## Docker

For those who do not have a Go development environment available, or simply want
to run an extended version of `k6` as a container, Docker is an option to build
and run the application.

The following command will build a custom `k6` image incorporating the `xk6-sql` extension
built from the local source files.
```shell
docker build -t gcfabri/k6-for-sql:latest .
```
