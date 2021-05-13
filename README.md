# global_docker_compose

`global_docker_compose` is a centralized way to manage your external dependencies across multiple projects. You start up your Docker services once, and all relevant ports are exposed so that you can use them in your apps. You don't need special `docker-compose.yml` files in each app, nor do you have multiple versions of MySQL or Kafka running around.

The idea behind `global_docker_compose` is to have everything *but* your app running in a Docker container. `global_docker_compose` is the central place to manage making those containers "good", including volumes, correct port exposure, hostnames, etc.

This tool is specifically to be used for *local development*, not for integration testing on CI or production.

## Usage

You must have Ruby and Docker installed to use this tool.

First install it as a Ruby gem:

`gem install global_docker_compose`

You now should be able to access the `global_docker_compose` command line from anywhere on your computer that has access to the same gemset. If you are using RVM, you will need to do this for each installed Ruby version, or simply add this to your Gemfile.

`global_docker_compose` has multiple sub-commands, most of which should be familiar:

* `global_docker_compose up --service=<service1> <service2>`: Bring up a list of services as defined by the table below.
* `global_docker_compose down --service=<service1> <service2>`: Bring down the specificed services.
* `global_docker_compose down`: Bring down all services.
* `global_docker_compose ps`: Show all running services that were configured using the tool.
* `global_docker_compose logs`: Print out logs.
* `global_docker_compose exec <service> <command>` Execute a command on an existing service.
* `global_docker_compose mysql --service=<service>` Start a MySQL client against whatever MySQL service is provided (e.g. `mysql56`).
* `global_docker_compose redis_cli` Start the Redis CLI (assuming `redis` is running)

The recommended usage of this command is via a shell script that lives in your project which automatically passes through the services that the app cares about. For example, in an executable file called `gdc`:

```shell
global_docker_compose "$@" --services=mysql57 redis kafka
```

When you call e.g. `gdc up` it will automatically pass everything through to the `global_docker_compose` command which will correspond to `global_docker_compose up --services=mysql57 redis kafka`. All commands will understand this option and use it to tailor the subcommands to the project settings. 

## Important Note

All services are exposed with the host IP of `127.0.0.1`. If you use `localhost`, it may not work. Whenever accessing local services (e.g. in configuration for your app), you should always use the IP address, not `localhost`.

## Additional Compose Files

`global_docker_compose` allows to supply an additional docker-compose file to augment the built-in ones with the `--compose_file` option. This file will be merged with the built-in ones using [docker-compose's merging rules](https://docs.docker.com/compose/extends/#adding-and-overriding-configuration). 

Note that if you define new services with this file, you must pass in the service name with the `--services` option along with the other ones.

As an example, if you have a separate docker-compose file that looks like this:

```yaml
version: '3.6'
services:
  postgres:
    image: postgres:11.1
    expose:
      - 5432
    environment:
      POSTGRES_PASSWORD: root
```

...you can start up Redis and Postgres with the following command:

```bash
global_docker_compose up --services=redis postgres --compose_file=./docker-compose.yml
```

## Supported Services

Key|Service|Ports
---|-------|-----
`mysql56`|MySQL 5.6|3307
`mysql57`|MySQL 5.7|3306
`mysql8`|MySQL 8.0|3308
`redis`|Redis|<ul><li>6379</li><li>8001 (Insights)</li></ul>
`kafka`|Kafka with Lenses Box|<ul><li>9092 (Kafka broker)</li><li>8081 (Schema Registry)</li><li>3030 (Lenses)</li></ul>

### MySQL

global_docker_compose supports MySQL 5.6, 5.7 and 5.8. To avoid port conflicts, the exported ports are as follows:

* 5.6: 3307
* 5.7: 3306
* 5.8: 3308

The reason 5.7 was given the "default" of 3306 is that it is currently the most common / default version to use.

### Kafka with Lenses

The recommended way to have your app talk to Kafka is to use Lenses Box, which is denoted by the `kafka` service. This includes the Kafka brokers, Zookeeper, schema registry, and Kafka Connect.

In order to use Lenses, you need to do the following:

1. Go [here](https://lenses.io/box/) and enter your e-mail address. You will get a free license key which is good for 6 months.
2. Edit your `~/.bashrc` or `~/.zshrc` and add the following lines:

```
export LENSES_KEY="https://licenses.lenses.io/download/lensesdl?id=<<<YOUR KEY HERE>>>"
```

3. Make sure that your Docker process can use at least 5GB of memory. On Mac you can do this via Docker for Mac -> Resources.

That's it! You can access your local Lenses at [http://localhost:3030](http://localhost:3030), with both username and password set to `admin`. It typically takes about 30-45 seconds to start Lenses. If you're not seeing it past that, make sure you've given Docker enough memory (see #3 above).

### Redis

Redis comes with a built-in `redisinsight` task which can show you the contents of your Redis installation. You can access Insights at [http://localhost:8001](http://localhost:8001).

## Contributing

Feel free to fork and add pull requests - we can add more services as necessary or tweak the ones we have.
