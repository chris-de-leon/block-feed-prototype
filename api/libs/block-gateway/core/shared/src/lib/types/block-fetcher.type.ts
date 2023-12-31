import { QueueNames } from "../enums/queue-names.enum"
import { JobNames } from "../enums/job-names.enum"
import { Queue, Worker, Processor } from "bullmq"

export namespace TBlockFetcher {
  export type TQueueName = QueueNames.BLOCK_FETCHER
  export type TJobName = JobNames.FETCH_BLOCK
  export type TQueueOutput = void
  export type TQueueInput = number
  export type TWorker = Worker<TQueueInput, TQueueOutput, TJobName>
  export type TQueue = Queue<TQueueInput, TQueueOutput, TJobName>
  export type TProcessor = Processor<TQueueInput, TQueueOutput, TJobName>
}
