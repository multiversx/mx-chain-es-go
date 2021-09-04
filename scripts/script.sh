IMAGE_NAME=elastic-container

start() {
  docker pull docker.elastic.co/elasticsearch/elasticsearch:7.9.0

  docker run -d --name ${IMAGE_NAME} -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.9.0

  sleep 30s
}

stop() {
  docker stop ${IMAGE_NAME}
}

delete() {
  curl -XDELETE http://localhost:9200/_all

  curl -X DELETE "http://localhost:9200/_template/*"

}

"$@"