# Quickstart

## Dependencies

A handful of dependencies need to be installed locally. If using MacOS, [Homebrew] is also required.

[Homebrew]: https://brew.sh/

### Just

[Just] is a command runner. It's similar to `make`, but without the added complexity of being a
build system.

```sh
# macos
brew install just

# linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin
```

[Just]: https://github.com/casey/just

### Docker

Docker can be installed using a variety of options:

- [Docker Engine](https://docs.docker.com/engine/install/) on Linux.
- [Docker Desktop](https://docs.docker.com/desktop/) on MacOS.

Alternatively, the below installs Docker on MacOS without needing a license for Docker Desktop.

```sh
# macos
brew install docker docker-credential-helper docker-compose docker-buildx colima

mkdir -p ~/.docker
echo '{\n\t"auths": {},\n\t"credsStore": "osxkeychain",\n\t"currentContext": "colima",\n\t"cliPluginsExtraDirs": [\n\t\t"/opt/homebrew/lib/docker/cli-plugins"\n\t]\n}' > ~/.docker/config.json

colima start

docker context ls
```

## Environment

The development environment can be started, stopped and restarted with:

```sh
# start
just dev-up

# stop
just dev-stop

# restart
dev dev-restart
```

To enter the development environment to run other commands, run:

```sh
just dev-exec
```

## Development

Once in the development environment, some extra tools need to be installed. These tools are updated
alongside other dependencies, so try to run this command fairly often to stay up-to-date. They will
also need to be reinstalled each time the development environment gets restarted.

```sh
just install-tools
```

There are many recipes that can be executed via `just`. The main ones are:

```sh
# sync dependencies, generate and format code
just checks

# scan the repository for issues
just scan

# lint the code
just lint

# run the unit tests
just unit
```

Some recipes call other recipes. For example, `checks` runs `fmt` which runs `fmt-buf` and `fmt-go`.
These more specific recipes can be run by themselves with:

```sh
just fmt-go lint-go
```

For a full list of recipes and what each does, run:

```sh
just

# or
just --list
```
