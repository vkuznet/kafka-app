### Setup local kafka pipeline
To setup local kafka pipeline we need the following:
- a node with running docker daemon
- a docker compose installation
- a kafka docker github repo

```
# install docker compose
curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /wma/vk/docker-compose
export PATH=$PATH:/wma/vk
chmod +x docker-compose

# clone kafka docker github repo
git clone git@github.com:wurstmeister/kafka-docker.git
cd kafka-docker
# vim docker-compose.yml
version: '2'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"
  kafka:
    build: .
    ports:
      - "9092"
    environment:
      KAFKA_ADVERTISED_HOST_NAME: 188.185.115.163
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

# start kafka
docker-compose ps
docker-compose down
docker-compose up -d
docker-compose ps
          Name                        Command               State                           Ports
--------------------------------------------------------------------------------------------------------------------------
kafka-docker_kafka_1       start-kafka.sh                   Up      0.0.0.0:9094->9094/tcp,:::9094->9094/tcp
kafka-docker_zookeeper_1   /bin/sh -c /usr/sbin/sshd  ...   Up      0.0.0.0:2181->2181/tcp,:::2181->2181/tcp, 22/tcp,
                                                                    2888/tcp, 3888/tcp

# get logs
docker-compose logs

# view existing topics with kafkacat
# see https://github.com/edenhill/kafkacat
docker run -it --network=host edenhill/kafkacat:1.6.0 -b localhost:9094 -L

# run client
cd ..
git clone git@github.com:vkuznet/kafka-app.git
cd kafka-app
go build
./kafka-app
```

### References

- [Docker compose](https://docs.docker.com/compose/install/)
- [Kafka docker](https://github.com/wurstmeister/kafka-docker)
- [Kafka local setup](https://www.kimsereylam.com/kafka/docker/2020/10/16/setup-local-kafka-with-docker.html)
- [Go client](https://www.sohamkamani.com/golang/working-with-kafka/)
