include Makefile.vars
include Makefile.common



# to override a target from Makefile.common just redefine the target.
# you can also chain the original atlas target by adding
# -atlas to the dependency of the redefined target
dapr run --app-id pub --components-path ./conf  --app-port 1250 --dapr-grpc-port 3501   --app-protocol grpc  go run cmd/server/*