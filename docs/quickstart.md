# Quickstart

## Dependencies

### Just

```sh
# macos
brew install just

# linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin
```

### Docker

```sh
# macos
brew install docker docker-credential-helper docker-compose docker-buildx colima

mkdir -p ~/.docker
echo '{\n\t"auths": {},\n\t"credsStore": "osxkeychain",\n\t"currentContext": "colima",\n\t"cliPluginsExtraDirs": [\n\t\t"/opt/homebrew/lib/docker/cli-plugins"\n\t]\n}' > ~/.docker/config.json

colima start

docker context ls
```

