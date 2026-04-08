using Commander.Core.Entities;

namespace Commander.Core.Ports;

public interface IJobStore
{
  bool StoreJob(Job job);
  Job GetJob(Guid jobId);
}
