# sandbox

Experimental Go microservice.

## Sanity tests

```sh
# success
echo '{}' | grpc-client-cli -a localhost:5000 -s Health -m Check

# invalid argument
echo '{}' | grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource

# unimplemented
echo '{"resource":{}}' | grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource
```
