# wwenet

> A tool for downloading from the WWE Network before it goes the way of the dodo!

Inspired by [WWE-Network-Downloader](https://github.com/freyta/WWE-Network-Downloader) with extra functionality I was interested in

I already pay for Netflix on top of the WWE Network but I don't want to watch everything while perpetually connected to a French VPN so this tool is nice for having a local archive for when connectivity is questionable.

Some of the content on the WWE Network is also not legally available in many countries, while still not being available on Netflix either.

## Network

Minimum headers:
- Authorization
- Realm
- x-api-key

Recommended:
- x-app-var
- User-Agent

## Development

### Accessing WWE Network

[Mullvad](https://mullvad.net/en) can be used with a French exit node.

A paid WWE Network account is required to access the network.

### Database migrations

Data is stored in `wwenet.db` with migrations being created using `sqlc`.

Queries are added in `query.sql` which generates Go code inside `storage`.

### Taskfile commands

#### Building a binary

```console
task build
```

#### Creating a new migration

```console
task migrate:create -- <name>
```

Creates `migrations/<timestamp>_<name>.sql`