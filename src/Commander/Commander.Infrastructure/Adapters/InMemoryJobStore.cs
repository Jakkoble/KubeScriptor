using System.Collections.Concurrent;
using Commander.Core.Entities;
using Commander.Core.Ports;

namespace Commander.Infrastructure.Adapters;

public class InMemoryJobStore : IJobStore
{
  private readonly ConcurrentDictionary<Guid, Job> _jobs = new();

  public bool StoreJob(Job job)
  {
    return _jobs.TryAdd(job.Id, job);
  }

  public Job GetJob(Guid jobId)
  {
    if (_jobs.TryGetValue(jobId, out var job))
    {
      return job;
    }

    throw new KeyNotFoundException($"Job with ID {jobId} not found.");
  }
}
