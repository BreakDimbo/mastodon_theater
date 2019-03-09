This is a bot client for mastodon, it supports:

* post toot according to the arrangemant in script Regularly (default 60m for every line)
* reply LOVE_YOU
* reply anything according to any key word (if you are able to program) 


# Install
* clone the repo
```
git clone git@github.com:BreakDimbo/mastodon_theater.git
```

* [install golang](https://golang.org/doc/install)

* install dependency(not necessary)
```
cd mastodon_theater
make deps
```

* make theater
```
make theater
```

* [install redis](https://redis.io/download)
```
$ wget http://download.redis.io/releases/redis-5.0.3.tar.gz
$ tar xzf redis-5.0.3.tar.gz
$ cd redis-5.0.3
$ make
```
* start the redis-server
```
redis-5.0.3/src/redis-server
```

# Set Up and Run

* set the config file
```
cp mastodon_theater/config/production.demo.toml mastodon_theater/config/production.toml
```
* open the file production.toml and set the client_email and client_password
* edit the mastodon_theater/config/steinGate1.1.txt file as you like according to the format
* if you want to refresh the script in steinGate1.1.txt from the scratch, remember to flush the redis
```
redis-5.0.3/src/redis-cli flushdb
```

* start the service
```
mastodon_theater/bin/theater -env=production
```


