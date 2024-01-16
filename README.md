# global_docker_compose

`global_docker_compose` is a centralized way to manage your external dependencies across multiple projects. You start up your Docker services once, and all relevant ports are exposed so that you can use them in your apps. You don't need special `docker-compose.yml` files in each app, nor do you have multiple versions of MySQL or Kafka running around. 

The idea behind `global_docker_compose` is to have everything *but* your app running in a Docker container. `global_docker_compose` is the central place to manage making those containers "good", including volumes, correct port exposure, hostnames, etc.

This tool is specifically to be used for *local development*, not for integration testing on CI or production.

## Installing

If you're running MacOS, you can install `global_docker_compose` via `homebrew`:

`brew install wishabi/flipp/global_docker_compose`

Alternatively, you can download executables from the [Releases page](https://github.com/wishabi/global-docker-compose/releases), unzip it and put it somewhere in your `PATH` variable.

Or you can build it from source by running the following from the root directory of this repo:

`go build -o global_docker_compose cmd/gdc/main.go`

## Usage

`global_docker_compose` has multiple sub-commands, most of which should be familiar:

* `global_docker_compose up --service=<service1>,<service2>`: Bring up a list of services as defined by the table below.
* `global_docker_compose down {service}`: Bring down the specificed service, or all services if not provided.
* `global_docker_compose down`: Bring down all services.
* `global_docker_compose ps`: Show all running services that were configured using the tool.
* `global_docker_compose config`: Print out the docker compose config file being used.
* `global_docker_compose logs {service}`: Print out logs for the specified service, or all services if not provided.
* `global_docker_compose exec <service> <command>` Execute a command on an existing service.
* `global_docker_compose mysql --service=<service> {input_file}` Start a MySQL client against whatever MySQL service is provided (e.g. `mysql56`). If an input file is provided, execute the statements in the input file. Additional services can be specified in the `<service>` parameter; they will be ignored.
* `global_docker_compose redis_cli` Start the Redis CLI (assuming `redis` is running)

The recommended usage of this command is via a shell script that lives in your project which automatically passes through the services that the app cares about. For example, in an executable file called `gdc`:

```shell
global_docker_compose "$@" --services=mysql57,redis,kafka
```

When you call e.g. `gdc up` it will automatically pass everything through to the `global_docker_compose` command which will correspond to `global_docker_compose up --services=mysql57,redis,kafka`. All commands will understand this option and use it to tailor the subcommands to the project settings. This allows your dev setup to be both simple and consistent: in all projects you use the same commands, `gdc up`, `gdc down`, `gdc mysql` etc. without having to worry about which versions or dependencies are installed.

Note that it's recommended to have the current directory in your PATH so you don't have to keep typing `./gdc`. In your `~/.bashrc` or `~/.zshrc` add:
```bash
EXPORT PATH=.:$PATH
```

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
global_docker_compose up --services=redis,postgres --compose_file=./docker-compose.yml
```

## Supported Services

Key|Service|Ports
---|-------|-----
`mysql56`|MySQL 5.6|3307
`mysql57`|MySQL 5.7|3306
`mysql8`|MySQL 8.0|3308
`redis`|Redis|<ul><li>6379</li><li>8001 (Insights)</li></ul>
`kafka`|Kafka with Lenses Box|<ul><li>9092 (Kafka broker)</li><li>8081 (Schema Registry)</li><li>3030 (Lenses)</li></ul>
`mailcatcher`|Mailcatcher|<ul><li>1025 (SMTP server)</li><li>1080 (UI)</li></ul>
`dynamodb`|DynamoDB|<ul><li>8000</li><li>8099 (Admin Dashboard)</li></ul>
`opensearch`|OpenSearch|<ul><li>9200</li><li>5601 (Dashboard)</li></ul>

### MySQL

global_docker_compose supports MySQL 5.6, 5.7 and 5.8. To avoid port conflicts, the exported ports are as follows:

* 5.6: 3307
* 5.7: 3306
* 5.8: 3308

The reason 5.7 was given the "default" of 3306 is that it is currently the most common / default version to use.

#### Exporting and Importing databases to GDC
To get a dump ready from your local databases to GDC, you must:

1. First set up permissions on your users for localhost by going into your local mysql server as the root user `mysql -u root`
```sql
-- Grant these perms for all users
> GRANT SHOW VIEW, SELECT, LOCK TABLES ON *.* TO ''@'localhost';
-- Or just grant them for a specific user (like fadmin)
> GRANT SHOW VIEW, SELECT, LOCK TABLES ON *.* TO 'fadmin'@'localhost';
```
2. Then actually dump the databases or selected databases with the `--databases` option
```sh
mysqldump --single-transaction -h 127.0.0.1 -P 3306 --all-databases --no-tablespaces > ./dump.sql
# In a space-separated list, list it after `--databases` flag, e,g. fadmin dbs
mysqldump --single-transaction -h 127.0.0.1 -P 3306 --databases fadmin_development fadmin_test --no-tablespaces > ./dump.sql
```
3. Now import the dump through gdc. If your app is already on gdc, the `mysql` command will default to the version of MySQL that it currently uses. Otherwise, you can speficy the selected databases in Step 2. and MySQL version using `global_docker_compose mysql --service=<service> <dump_file>`
```sh
./gdc mysql ./dump.sql
```
4. Rinse and repeat! Don't forget to remove the dump-file to reduce clutter on your local machine.
```
rm ./dump.sql
```

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

If you need to advertise a different host for the broker and schema registry, you can set the `KAFKA_ADV_HOST` environment variable before invoking `global_docker_compose`. Here's an example:

```
MY_IP=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')
KAFKA_ADV_HOST=$MY_IP global_docker_compose --services=kafka
```

### Redis

Redis comes with a built-in `redisinsight` task which can show you the contents of your Redis installation. You can access Insights at [http://localhost:8001](http://localhost:8001).

To allow it to access your local Redis, you *can't* use `127.0.0.1`. Instead [you need to use your host IP address](https://collabnix.com/running-redisinsight-using-docker-compose/) - use `ifconfig en0` to find it:

<pre>
➜  ~/RubymineProjects/cp-legacy git:(main) ✗ ifconfig en0
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	options=400<CHANNEL_IO>
	ether 14:7d:da:49:ca:df 
	inet6 fe80::10f0:fee7:2a6f:be0e%en0 prefixlen 64 secured scopeid 0xa 
	inet <b>192.168.1.9</b> netmask 0xffffff00 broadcast 192.168.1.255
	nd6 options=201<PERFORMNUD,DAD>
	media: autoselect
	status: active
</pre>

In this case you'd enter `192.168.1.9` as your hostname in Redis Insight.

### Mailcatcher

[Mailcatcher](https://mailcatcher.me/) is a local SMTP server you can use to send and view e-mails. Set up your mail sending code to talk
to port 1025 and you can view the mail that got sent on port 1080.

On Rails, this would look like:

```ruby
  config.action_mailer.delivery_method = :smtp
  config.action_mailer.smtp_settings = { address: '127.0.0.1', port: 1025 }
```

## Releasing

Releases are done via [GoReleaser](https://goreleaser.com/intro/) which is run on GitHub Actions whenever a new tag is pushed. See the file `.goreleaser.yml` for more information.

To bump a new version, add a Git tag and also update the version reported in `cmd/gdc/commands/root.go`.

## Contributing

Feel free to fork and add pull requests - we can add more services as necessary or tweak the ones we have.

You can build the executable locally by running:

    go build -o global_docker_compose cmd/gdc/main.go
### Adding a new service

The steps to add a new service are:

1. Add the service to `gdc/docker-compose.yml`.
2. Add the functionality for your command in `gdc/docker.go`.
3. Add a new command under `cmd/gdc/commands`. You can copy and paste an existing one or make changes. `global_docker_compose` uses [Cobra](https://github.com/spf13/cobra) for command-line flags, validations, help text and arguments, so please read that documentation for more info.
4. Put up your PR!


## Troubleshooting

If you are on an M* macOS machine, you might run into trouble getting GDC to run. 

### Error: `no matching manifest for linux/arm64/v8 in the manifest list entries`

We had to apply these fixes:

```sh
docker pull --platform linux/x86_64 MySQL
export DOCKER_DEFAULT_PLATFORM=linux/x86_64/v8 # you might want to put this in your rc or profile file
```

Then re-run `gdc up`.

#### References
[StackOverflow Ref]([url](https://stackoverflow.com/questions/65456814/docker-apple-silicon-m1-preview-mysql-no-matching-manifest-for-linux-arm64-v8/67361161#67361161)https://stackoverflow.com/questions/65456814/docker-apple-silicon-m1-preview-mysql-no-matching-manifest-for-linux-arm64-v8/67361161#67361161)
[Use Colima to run Docker containers]([url](https://smallsharpsoftwaretools.com/tutorials/use-colima-to-run-docker-containers-on-macos/)https://smallsharpsoftwaretools.com/tutorials/use-colima-to-run-docker-containers-on-macos/)

### Error: `unknown shorthand flag: 'p' in -p`

#### References

To fix this, we had to run:

```sh
ln -sfn $(brew --prefix)/opt/docker-compose/bin/docker-compose ~/.docker/cli-plugins/docker-compose
```

We also made sure `colima`, `docker`, `docker-compose` and `XCode` were all up to date.

[Fadmin Troubleshooting docs]([url](https://flippit.atlassian.net/wiki/spaces/CTLR/pages/10448437272/End-to-end+Fadmin+WES+local+setup#Troubleshooting)https://flippit.atlassian.net/wiki/spaces/CTLR/pages/10448437272/End-to-end+Fadmin+WES+local+setup#Troubleshooting) 

