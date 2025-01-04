# sandbox

Experimental Go microservice.

## Sanity tests

```sh
# success
echo '{}' | grpc-client-cli -a localhost:5000 -s Health -m Check

# invalid argument
echo '{}' | grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource

# success
echo '{"resource":{"name":"example"}}' | grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource
```
