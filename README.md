# tracking

docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

may be port problem so update this

docker run -d -p 6831:6831/udp -p 14268:14268 -p 16686:16686 jaegertracing/all-in-one:latest


http://localhost:16686    --> jagaur

