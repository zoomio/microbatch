# microbatch

Micro-batching is a subset of the [batch processing](https://en.wikipedia.org/wiki/Batch_processing), which is aimed to improve latencies of the long batch processing tasks. It is sort of a middle ground in-between the conventional batch processing (which handles batches of a big size and takes long time for results to be available) and the streaming (or stream processing) where inputs aren't aggregated in the batches and everything is handled as it appears. The goal is to reduce latency and still to allow for reasonable throughput.

So the **microbatch** is a aimed to provide simple primitive for adopting micro-batch technic inside your projects.

## Usage

- You'd need to implement `BatchProcessor` interface to handle your specific jobs.
- New instance of the `MicroBatch` can be created with the `New(yourBatchProcessor)`.
- Created instance of the `MicroBatch` can be started with the `yourMicroBatch.Start(c)`, where c - is a `context.Context`, `context.Context` is used for signalling when `MicroBatch` needs to stop (so use `context.WithCancel`).
- New running instance of the `MicroBatch` can be created with the `NewRunning(c, yourBatchProcessor)`.
- New job can be submitted to the `MicroBatch` via `yourMicroBatch.Submit(yourJob)`.
- A job is an instance of the `Job` type.