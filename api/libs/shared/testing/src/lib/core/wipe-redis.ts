import { exec } from "child_process"

export const wipeRedis = async () => {
  return await new Promise<
    Readonly<{
      stdout: string
      stderr: string
    }>
  >((res, rej) => {
    exec(`docker exec redis redis-cli "FLUSHALL"`, (err, stdout, stderr) => {
      if (err != null) {
        rej(err)
        return
      }
      res({
        stdout,
        stderr,
      })
    })
  })
}
