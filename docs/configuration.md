# Loader configuration file format

| Parameter name               | Data type | Possible values                   | Default value | Description                                                                     |
|------------------------------|-----------|-----------------------------------|---------------|---------------------------------------------------------------------------------|
| Seed                         | int64     | any                               | 42            | Seed for specification generator (for reproducibility)                          |
| YAMLSelector                 | string    | wimpy, container, firecracker     | container     | Service YAML depending on sandbox type                                          |
| EndpointPort                 | int       | >0                                | 80            | Port to be appended to the service URL                                          |
| TracePath                    | string    | string                            | data/traces   | Folder with Azure trace dimensions (invocations.csv, durations.csv, memory.csv) |
| IATDistribution              | string    | exponential, uniform, equidistant | exponential   | IAT distribution                                                                |
| ExperimentDuration           | int       | >0                                | 1             | Experiment duration in minutes of trace to execute excluding warmup             |
| WarmupDuration               | int       | >0                                | 0             | Warmup duration in minutes(disabled if zero)                                    |
| IsPartiallyPanic             | bool      | true/false                        | false         | Pseudo-panic-mode only in Knative                                               |
| EnableZipkinTracing          | bool      | true/false                        | false         | Show loader span in Zipkin traces                                               |
| EnableMetricsScrapping       | bool      | true/false                        | false         | Scrap cluster-wide metrics                                                      |
| GRPCConnectionTimeoutSeconds | int       | > 0                               | 60            | Timeout for establishing a gRPC connection                                      |
| GRPCFunctionTimeoutSeconds   | int       | > 0                               | 90            | Maximum time given to function to execute[^1]                                   |

[^1] - Function can execute for at most 15 minutes as in AWS
Lambda; https://aws.amazon.com/about-aws/whats-new/2018/10/aws-lambda-supports-functions-that-can-run-up-to-15-minutes/